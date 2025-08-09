package nodesync

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/sync/poolsync"
)

type syncer struct {
	uow    poolsync.NodeUoW
	client poolsync.NodeClient
}

// sync node state between node (available via client) and uow.
// required node and user states are described in uow.
// we want to try put node to this state via client,
// and after successful or unsuccessful attempt
// update actual node state according to changes we made or not.
// the situation I hate and try to avoid is
//
//   - i update node via client (for example, remove user)
//
//   - then trying to write it to uow, and all attempts failed
//     due to database connection lost or db host limitations or whatever
//
//   - after i fix it, user marked in database as active,
//     but it's actually not. and i have no clue what is going wrong
//     and what items in database are now incorrect. moreover, and the worst,
//     user tries to made something, some parts of service use database as
//     source of data, other - communicates with node client, and inconsistency
//     between them leads to not-so-interesting errors. hate it.
//
//     to avoid it, let's mark items we are going to modify as 'Unknown value'
//     in uow, and after attempt, try to write to uow actual values.
//     the worst case is node modified but next uow update fails,
//     but now invalid values are explicitly marked as 'Unknown' in uow,
//     so it is possible to detect and handle it.
func (s *syncer) SyncNodeState(ctx context.Context) (err error) {
	if s == nil || s.uow == nil || s.client == nil {
		return fmt.Errorf("node: sync state: %w", errdefs.ErrNilObjectCall)
	}

	// get current.
	curr, prev, target, err := s.fetchNodeStatus(ctx)
	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	// required node and user states
	// we have 3 options: start/stop node, sync out of sync users.
	// when sync node users, change node state if it differs
	// from current stored state
	switch {
	case target == models.NodeStatusRunning && curr == models.NodeStatusStopped:
		err = s.startNode(ctx)
	case target == models.NodeStatusStopped && curr == models.NodeStatusRunning:
		err = s.stopNode(ctx)
	case target == models.NodeStatusRunning && curr == models.NodeStatusRunning:
		err = s.syncNodeUsers(ctx, curr != prev)
	}

	if err != nil {
		return fmt.Errorf("node: sync state: %w", err)
	}

	return nil
}

func (s *syncer) fetchNodeStatus(ctx context.Context) (
	curr, prev, target models.NodeStatus, err error,
) {
	// fetch stored node status
	if err = s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) (err error) {
		target, prev, err = uowctx.FetchNodeStatus(ctx)
		return
	}); err != nil {
		err = fmt.Errorf("node: sync state: %w", err)
		return
	}

	// fetch curr node status if required
	// if node should be running and currently it's not surely stopped, let's
	// check and update its state (node can sometimes switch
	// from any state to stop or disconnected or off due to its internal faults
	// or connection errors)
	curr = prev
	if target == models.NodeStatusRunning && prev != models.NodeStatusStopped {
		if curr, err = s.client.CheckStatus(ctx); err != nil {
			err = fmt.Errorf("node: sync state: %w", err)
			return
		}
	}
	return
}

func (s *syncer) startNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = s.updateStoredStatus(ctx, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	active, inactive, err := s.getUsers(ctx)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// start node
	clientConfig, err := s.client.Start(ctx, active)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// update stored node state
	if err := s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) (err error) {
		if err = uowctx.PatchPendingSyncs(ctx, s.getUsersPatch(active, inactive)); err != nil {
			return
		}
		if err = uowctx.UpdateCurrentStatus(ctx, models.NodeStatusRunning); err != nil {
			return
		}
		if err = uowctx.UpdateClientConfig(ctx, *clientConfig); err != nil {
			return
		}
		return
	}); err != nil {
		return fmt.Errorf("start node: update state: %w", err)
	}

	return nil
}

func (s *syncer) getUsers(ctx context.Context) (active, inactive []models.UserProfile, err error) {
	// get all users
	var users []models.User
	if err = s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) (err error) {
		users, err = uowctx.ListUsers(ctx)
		return
	}); err != nil {
		err = fmt.Errorf("update node status: %w", err)
		return
	}

	// split onto active and inactive
	for _, u := range users {
		switch u.TargetStatus {
		case models.UserStatusActive:
			active = append(active, u.Profile)
		case models.UserStatusInactive:
			inactive = append(inactive, u.Profile)
		}
	}

	return
}

func (s *syncer) getUsersPatch(active, inactive []models.UserProfile) []models.UserStatusPatch {
	patch := make([]models.UserStatusPatch, 0, len(active)+len(inactive))
	for _, u := range active {
		patch = append(patch, models.UserStatusPatch{
			UserID: u.ID,
			Status: models.UserStatusActive,
		})
	}
	for _, u := range inactive {
		patch = append(patch, models.UserStatusPatch{
			UserID: u.ID,
			Status: models.UserStatusInactive,
		})
	}
	return patch
}

func (s *syncer) stopNode(ctx context.Context) (err error) {
	// safe state-changing stuff
	if err = s.updateStoredStatus(ctx, models.NodeStatusUnknown); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	if err = s.client.Stop(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	// dont update actual status of users on disabled nodes because it has no matter.
	// it updates when node started again.
	// TODO: maybe update?
	if err = s.updateStoredStatus(ctx, models.NodeStatusStopped); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	return nil
}

func (s *syncer) syncNodeUsers(ctx context.Context, updateNodeStatus bool) error {
	pending, err := s.getPendingSyncs(ctx)
	if err != nil {
		return err
	}

	if len(pending) == 0 && !updateNodeStatus {
		return nil
	}

	usersUpdate, prePatch, postPatch := s.buildUserUpdate(pending)

	if err := s.applyNodeStatePatch(ctx, prePatch); err != nil {
		return err
	}

	if err := s.client.UpdateUsers(ctx, usersUpdate); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}

	if err := s.applyNodeStatePatch(ctx, postPatch); err != nil {
		return err
	}

	return nil
}

func (s *syncer) getPendingSyncs(ctx context.Context) (pending []models.UserSyncStatus, err error) {
	if err = s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) (err error) {
		pending, err = uowctx.FindPendingSyncs(ctx)
		return err
	}); err != nil {
		err = fmt.Errorf("get pending syncs: %w", err)
		return
	}
	return
}

func (s *syncer) buildUserUpdate(syncs []models.UserSyncStatus) (
	update models.NodeUsersUpdate, prePatch, postPatch []models.UserStatusPatch,
) {
	prePatch = make([]models.UserStatusPatch, 0, len(syncs))
	postPatch = make([]models.UserStatusPatch, 0, len(syncs))
	update.Add = make([]models.UserProfile, 0, len(syncs))
	update.Remove = make([]models.UserProfile, 0, len(syncs))

	for _, u := range syncs {
		switch u.User.TargetStatus {
		case models.UserStatusActive:
			update.Add = append(update.Add, u.User.Profile)
		case models.UserStatusInactive:
			update.Remove = append(update.Remove, u.User.Profile)
		}
		prePatch = append(prePatch, models.UserStatusPatch{
			UserID: u.User.Profile.ID,
			Status: models.UserStatusUnknown,
		})
		postPatch = append(postPatch, models.UserStatusPatch{
			UserID: u.User.Profile.ID,
			Status: u.User.TargetStatus,
		})
	}
	return
}

func (s *syncer) applyNodeStatePatch(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) error {
		if err := uowctx.PatchPendingSyncs(ctx, patch); err != nil {
			return err
		}
		if err := uowctx.UpdateCurrentStatus(ctx, models.NodeStatusRunning); err != nil {
			return err
		}
		return nil
	})
}

func (s *syncer) updateStoredStatus(ctx context.Context, status models.NodeStatus) error {
	if err := s.uow.Do(ctx, func(uowctx poolsync.NodeUoWContext) error {
		return uowctx.UpdateCurrentStatus(ctx, status)
	}); err != nil {
		return fmt.Errorf("update node status: %w", err)
	}
	return nil
}
