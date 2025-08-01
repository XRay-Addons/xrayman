package nodesyncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeSyncer struct {
	storage NodeStorage
	client  NodeClient
}

func NewNodeSyncer(storage NodeStorage, client NodeClient) (*NodeSyncer, error) {
	if storage == nil {
		return nil, fmt.Errorf("node syncer: init: storage: %w", errdefs.ErrNilArgPassed)
	}
	if client == nil {
		return nil, fmt.Errorf("node syncer: init: client: %w", errdefs.ErrNilArgPassed)
	}
	return &NodeSyncer{
		storage: storage,
		client:  client,
	}, nil
}

// sync node state between node (available via client) and storage.
// required node and user states are described in storage.
// we want to try put node to this state via client,
// and after successful or unsuccessful attempt
// update actual node state according to changes we made or not.
// the situation I hate and try to avoid is
// - i update node via client (for example, remove user)
// - then trying to write it to storage, and all attempts failed
//   due to database connection lost or db host limitations or whatever
// - after i fix it, user marked in database as active,
//   but it's actually not. and i have no clue what is going wrong
//   and what items in database are now incorrect. moreover, and the worst,
//   user tries to made something, some parts of service use database as
//   source of data, other - communicates with node client, and inconsistency
//   between them leads to not-so-interesting errors. hate it.
//
//   to avoid it, let's mark items we are going to modify as 'Unknown value'
//   in storage, and after attempt, try to write to storage actual values.
//   the worst case is node modified but next storage update fails,
//   but now invalid values are explicitly marked as 'Unknown' in storage,
//   so it is possible to detect and handle it.

func (s *NodeSyncer) SyncState(ctx context.Context) error {
	if s == nil || s.storage == nil || s.client == nil {
		return fmt.Errorf("node: sync state: %w", errdefs.ErrNilObjectCall)
	}

	target, previous, err := fetchStatus(ctx, s.storage)
	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	// if node should be running and currently it's not surely stopped, let's
	// check and update its state (node can sometimes switch
	// from any state to stop or disconnected or off due to its internal faults
	// or connection errors)
	current := previous
	if target == models.NodeStatusRunning && previous != models.NodeStatusStopped {
		if current, err = s.client.CheckStatus(ctx); err != nil {
			return fmt.Errorf("node: sync state: %w", err)
		}
	}

	// required node and user states
	// we have 3 options: start/stop node, sync out of sync users.
	// when sync node users, change node state if it differs
	// from current stored state
	switch {
	case target == models.NodeStatusRunning && current == models.NodeStatusStopped:
		err = s.startNode(ctx)
	case target == models.NodeStatusStopped && current == models.NodeStatusRunning:
		err = s.stopNode(ctx)
	case target == models.NodeStatusRunning && current == models.NodeStatusRunning:
		err = s.syncNodeUsers(ctx, current != previous)
	}

	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	return nil
}

func (s *NodeSyncer) startNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = updateCurrentStatus(ctx, s.storage, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		if syncErr := updateCurrentStatus(ctx, s.storage, models.NodeStatusStopped); syncErr != nil {
			err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
		}
	}()

	users, err := listUsers(ctx, s.storage)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// select only enabled users to start node
	activeUsers := selectActiveUsers(users)

	// start node
	nodeConfig, err := s.client.Start(ctx, activeUsers)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// update stored node state
	err = updateNode(ctx, s.storage,
		getUsersPatch(users),
		models.NodeStatusRunning,
		nodeConfig)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	return nil
}

func (s *NodeSyncer) stopNode(ctx context.Context) (err error) {
	// safe state-changing stuff. change state to unknown before
	// stopping, to stopped on success or back to running on fail
	if err = updateCurrentStatus(ctx, s.storage, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		if syncErr := updateCurrentStatus(ctx, s.storage, models.NodeStatusRunning); syncErr != nil {
			err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
		}
	}()

	if err = s.client.Stop(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	// dont update actual status of users on disabled nodes because it has no matter.
	// it updates when node started again.
	if err = updateCurrentStatus(ctx, s.storage, models.NodeStatusStopped); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	return nil
}

func selectActiveUsers(users []models.User) []models.UserProfile {
	active := make([]models.UserProfile, 0, len(users))
	for _, u := range users {
		if u.TargetStatus == models.UserStatusActive {
			active = append(active, u.Profile)
		}
	}
	return active
}

func getUsersPatch(users []models.User) []models.UserStatusPatch {
	patch := make([]models.UserStatusPatch, 0, len(users))
	for _, u := range users {
		patch = append(patch, models.UserStatusPatch{
			UserID: u.Profile.ID,
			Status: u.TargetStatus,
		})
	}
	return patch
}

func (s *NodeSyncer) syncNodeUsers(ctx context.Context, updateNodeState bool) (err error) {
	// get users to update
	pendingSyncs, err := findPendingSyncs(ctx, s.storage)
	if err != nil {
		return fmt.Errorf("sync node users: %w", err)
	}

	// if nothing to update, return
	if len(pendingSyncs) == 0 && !updateNodeState {
		return nil
	}

	// create user state updates to `lock` and `unlock` them safe
	usersUpdate := models.NodeUsersUpdate{}
	preUpdatePatch := make([]models.UserStatusPatch, 0, len(pendingSyncs))
	postUpdatePatch := make([]models.UserStatusPatch, 0, len(pendingSyncs))
	for _, u := range pendingSyncs {
		switch u.User.TargetStatus {
		case models.UserStatusActive:
			usersUpdate.Add = append(usersUpdate.Add, u.User.Profile)
		case models.UserStatusInactive:
			usersUpdate.Remove = append(usersUpdate.Remove, u.User.Profile)
		}
		preUpdatePatch = append(preUpdatePatch, models.UserStatusPatch{
			UserID: u.User.Profile.ID,
			Status: models.UserStatusUnknown,
		})
		postUpdatePatch = append(postUpdatePatch, models.UserStatusPatch{
			UserID: u.User.Profile.ID,
			Status: u.User.TargetStatus,
		})
	}

	// prepare to update. if 'update state' flag passed,
	// change state to actual 'Running' on pre-update.
	// (and don't touch after update)
	if updateNodeState {
		err = updateNode(ctx, s.storage, preUpdatePatch, models.NodeStatusRunning, nil)
	} else {
		err = patchPendingSyncs(ctx, s.storage, preUpdatePatch)
	}
	if err != nil {
		return fmt.Errorf("sync node users: pre update patch: %w", err)
	}

	if err := s.client.UpdateUsers(ctx, usersUpdate); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}

	if err := patchPendingSyncs(ctx, s.storage, postUpdatePatch); err != nil {
		return fmt.Errorf("sync node users: pre update patch: %w", err)
	}

	return nil
}
