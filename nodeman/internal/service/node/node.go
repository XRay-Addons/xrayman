package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/models"
)

type Node struct {
	storage Storage
	api     API
}

func New(storage Storage, api API) (*Node, error) {
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

	required, previous, err := n.storage.GetNodeState(ctx)
	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}
	if required != models.NodeOn && required != models.NodeOff {
		return fmt.Errorf("invalid required state %v", required)
	}

	// if node should be running and currently it's not surely stopped, let's
	// check and update its state (node can sometimes switch
	// from any state to stop or disconnected or off due to its internal faults
	// or connection errors)
	current := previous
	if required == models.NodeOn && previous != models.NodeOff {
		if current, err = n.api.GetStatus(ctx); err != nil {
			return fmt.Errorf("node: sync state: %w", err)
		}
	}

	// required node and user states
	// we have 3 options: start/stop node, sync out of sync users.
	// when sync node users, change node state if it differs
	// from current stored state
	switch {
	case required == models.NodeOn && current == models.NodeOff:
		err = n.startNode(ctx)
	case required == models.NodeOff && current == models.NodeOn:
		err = n.stopNode(ctx)
	case required == models.NodeOn && current == models.NodeOn:
		err = n.syncNodeUsers(ctx, current != previous)
	}

	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	return nil
}

func (n *Node) startNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = n.updateStoredStatus(ctx, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err != nil {
			if syncErr := n.updateStoredStatus(ctx, models.NodeOff); syncErr != nil {
				err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
			}
		}
	}()

	allUsers, err := n.storage.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// select only enabled users to start node
	enabledUsers := getEnabledUsers(allUsers)

	// start node
	nodeProperties, err := n.api.Start(ctx, enabledUsers)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// write updated state
	writeUoW := n.storage.GetWriteUOW()
	writeUoW.SetActualStatus(models.NodeOn)
	writeUoW.SetActualUserStates(getUserStatusUpdates(allUsers))
	writeUoW.SetNodeProperties(*nodeProperties)
	if err = writeUoW.Do(ctx); err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	return nil
}

func (n *Node) stopNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = n.updateStoredStatus(ctx, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}
	defer func() {
		if err != nil {
			if syncErr := n.updateStoredStatus(ctx, models.NodeOn); syncErr != nil {
				err = errors.Join(err, fmt.Errorf("sync after start node: %w", err))
			}
		}
	}()

	if err = n.api.Stop(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	// dont update actual status of users on disabled nodes because it has no matter.
	// it updates when node started again.
	if err = n.updateStoredStatus(ctx, models.NodeOff); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	return nil
}

func (n *Node) updateStoredStatus(ctx context.Context, s models.NodeStatus) error {
	writeUOW := n.storage.GetWriteUOW()
	writeUOW.SetActualStatus(s)
	if err := writeUOW.Do(ctx); err != nil {
		return fmt.Errorf("update stored status: %w", err)
	}
	return nil
}

func getEnabledUsers(users []UserState) []models.User {
	enabled := make([]models.User, 0, len(users))
	for _, u := range users {
		if u.Status == models.UserEnabled {
			enabled = append(enabled, u.User)
		}
	}
	return enabled
}

func getUserStatusUpdates(users []UserState) []UserStatusUpdate {
	userStates := make([]UserStatusUpdate, 0, len(users))
	for _, u := range users {
		userStates = append(userStates, UserStatusUpdate{
			u.User.ID, u.Status,
		})
	}
	return userStates
}

func (n *Node) syncNodeUsers(ctx context.Context, updateState bool) (err error) {
	// get users to update
	oosUsers, err := n.storage.ListOutOfSyncUsers(ctx)
	if err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}

	// create user state updates to `lock` and `unlock` them safe
	usersUpdate := make([]UserState, 0, len(oosUsers))
	intentUpdate := make([]UserStatusUpdate, 0, len(oosUsers))
	applyUpdate := make([]UserStatusUpdate, 0, len(oosUsers))
	for _, u := range oosUsers {
		usersUpdate = append(usersUpdate, UserState{
			User:   u.User,
			Status: u.Required,
		})
		intentUpdate = append(intentUpdate, UserStatusUpdate{
			ID:     u.User.ID,
			Actual: models.UserStatusUnknown,
		})
		applyUpdate = append(applyUpdate, UserStatusUpdate{
			ID:     u.User.ID,
			Actual: u.Required,
		})
	}

	if err := n.updateStoredUsers(ctx, intentUpdate, updateState); err != nil {
		return fmt.Errorf("sync node users: %w", err)
	}
	if err := n.api.EditUsers(ctx, usersUpdate); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}
	if err := n.updateStoredUsers(ctx, applyUpdate, false); err != nil {
		return fmt.Errorf("sync node users: %w", err)
	}
	return nil
}

func (n *Node) updateStoredUsers(ctx context.Context,
	u []UserStatusUpdate, updateStatus bool) error {
	uow := n.storage.GetWriteUOW()
	uow.SetActualUserStates(u)
	if updateStatus {
		uow.SetActualStatus(models.NodeOn)
	}
	if err := uow.Do(ctx); err != nil {
		return fmt.Errorf("update stored users: %w", err)
	}
	return nil
}
