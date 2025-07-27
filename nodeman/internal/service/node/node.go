package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Node struct {
	storage Storage
	api     NodeAPI
}

func New(storage Storage, api NodeAPI) (*Node, error) {
	if storage == nil {
		return nil, fmt.Errorf("node init: storage: %w", errdefs.ErrNilArgPassed)
	}
	if api == nil {
		return nil, fmt.Errorf("node init: api: %w", errdefs.ErrNilArgPassed)
	}
	return &Node{
		storage: storage,
		api:     api,
	}, nil
}

// sync node state between node (available via api) and storage.
// required node and user states are described in storage.
// we want to try put node to this state via api,
// and after successful or unsuccessful attempt
// update actual node state according to changes we made or not.
// the situation I hate and try to avoid is
// - i update node via api (for example, remove user)
// - then trying to write it to storage, and all attempts failed
//   due to database connection lost or db host limitations or whatever
// - after i fix it, user marked in database as active,
//   but it's actually not. and i have no clue what is going wrong
//   and what items in database are now incorrect. moreover, and the worst,
//   user tries to made something, some parts of service use database as
//   source of data, other - communicates with node api, and inconsistency
//   between them leads to not-so-interesting errors. hate it.
//
//   to avoid it, let's mark items we are going to modify as 'Unknown value'
//   in storage, and after attempt, try to write to storage actual values.
//   the worst case is node modified but next storage update fails,
//   but now invalid values are explicitly marked as 'Unknown' in storage,
//   so it is possible to detect and handle it.

func (n *Node) SyncState(ctx context.Context) error {
	if n == nil || n.storage == nil || n.api == nil {
		return fmt.Errorf("node: sync state: %w", errdefs.ErrNilObjectCall)
	}

	target, previous, err := fetchStatus(ctx, n.storage)
	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	// if node should be running and currently it's not surely stopped, let's
	// check and update its state (node can sometimes switch
	// from any state to stop or disconnected or off due to its internal faults
	// or connection errors)
	current := previous
	if target == models.NodeStatusRunning && previous != models.NodeStatusStopped {
		if current, err = n.api.CheckStatus(ctx); err != nil {
			return fmt.Errorf("node: sync state: %w", err)
		}
	}

	// required node and user states
	// we have 3 options: start/stop node, sync out of sync users.
	// when sync node users, change node state if it differs
	// from current stored state
	switch {
	case target == models.NodeStatusRunning && current == models.NodeStatusStopped:
		err = n.startNode(ctx)
	case target == models.NodeStatusStopped && current == models.NodeStatusRunning:
		err = n.stopNode(ctx)
	case target == models.NodeStatusRunning && current == models.NodeStatusRunning:
		err = n.syncNodeUsers(ctx, current != previous)
	}

	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	return nil
}

func (n *Node) startNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = updateStatus(ctx, n.storage, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		if syncErr := updateStatus(ctx, n.storage, models.NodeStatusStopped); syncErr != nil {
			err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
		}
	}()

	users, err := listUsers(ctx, n.storage)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// select only enabled users to start node
	activeUsers := selectActiveUsers(users)

	// start node
	nodeConfig, err := n.api.Start(ctx, activeUsers)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// update stored node state
	err = nodeFullUpdate(ctx, n.storage,
		getUsersPatch(users),
		models.NodeStatusRunning,
		nodeConfig)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	return nil
}

func (n *Node) stopNode(ctx context.Context) (err error) {
	// safe state-changing stuff. change state to unknown before
	// stopping, to stopped on success or back to running on fail
	if err = updateStatus(ctx, n.storage, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		if syncErr := updateStatus(ctx, n.storage, models.NodeStatusRunning); syncErr != nil {
			err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
		}
	}()

	if err = n.api.Stop(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	// dont update actual status of users on disabled nodes because it has no matter.
	// it updates when node started again.
	if err = updateStatus(ctx, n.storage, models.NodeStatusStopped); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	return nil
}

func selectActiveUsers(users []models.UserTargetState) []models.UserProfile {
	active := make([]models.UserProfile, 0, len(users))
	for _, u := range users {
		if u.Target == models.UserStatusActive {
			active = append(active, u.User)
		}
	}
	return active
}

func getUsersPatch(users []models.UserTargetState) []models.UserStatusPatch {
	patch := make([]models.UserStatusPatch, 0, len(users))
	for _, u := range users {
		patch = append(patch, models.UserStatusPatch{
			UserID: u.User.ID,
			Status: u.Target,
		})
	}
	return patch
}

func (n *Node) syncNodeUsers(ctx context.Context, updateNodeState bool) (err error) {
	// get users to update
	pendingSyncs, err := findPendingSyncs(ctx, n.storage)
	if err != nil {
		return fmt.Errorf("sync node users: %w", err)
	}

	// if nothing to update, return
	if len(pendingSyncs) == 0 && !updateNodeState {
		return nil
	}

	// create user state updates to `lock` and `unlock` them safe
	usersUpdate := make([]models.UserTargetState, 0, len(pendingSyncs))
	preUpdatePatch := make([]models.UserStatusPatch, 0, len(pendingSyncs))
	postUpdatePatch := make([]models.UserStatusPatch, 0, len(pendingSyncs))
	for _, u := range pendingSyncs {
		usersUpdate = append(usersUpdate, models.UserTargetState{
			User:   u.User,
			Target: u.TargetStatus,
		})
		preUpdatePatch = append(preUpdatePatch, models.UserStatusPatch{
			UserID: u.User.ID,
			Status: models.UserStatusUnknown,
		})
		postUpdatePatch = append(postUpdatePatch, models.UserStatusPatch{
			UserID: u.User.ID,
			Status: u.TargetStatus,
		})
	}

	// prepare to update. if 'update state' flag passed,
	// change state to actual 'Running' on pre-update.
	// (and don't touch after update)
	if updateNodeState {
		err = patchAndUpdateStatus(ctx, n.storage, preUpdatePatch, models.NodeStatusRunning)
	} else {
		err = patchPendingSyncs(ctx, n.storage, preUpdatePatch)
	}
	if err != nil {
		return fmt.Errorf("sync node users: pre update patch: %w", err)
	}

	if err := n.api.UpdateUserStates(ctx, usersUpdate); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}

	if err := patchPendingSyncs(ctx, n.storage, postUpdatePatch); err != nil {
		return fmt.Errorf("sync node users: pre update patch: %w", err)
	}

	return nil
}
