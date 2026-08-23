package nodesync

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type syncer struct {
	storage Storage
	client  Client
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
	if s == nil || s.storage == nil || s.client == nil {
		return errdefs.NilCall()
	}

	// get current node state
	curr, prev, target, err := s.fetchNodeStatus(ctx)

	// required node and user states
	// we have 4 options:
	//  - mark node as unavailable
	//  - start/stop node
	//  - sync out of sync users.
	// when sync node users, change node state if it differs
	// from current stored state.
	// be careful: state when current status is unknown and
	// target status is stopped ignored (but node may work)
	switch {
	case err != nil && prev != models.NodeStatusUnknown:
		// set additional time to this fallback (original ctx could be)
		// already cancelled due to unavailable node
		fallbackTO := time.Second
		fallbackCtx, fallbackCancel := context.WithTimeout(context.Background(), fallbackTO)
		defer func() { fallbackCancel() }()
		err = xerr.Join(err, s.markAsUnavailable(fallbackCtx))
	case target == models.NodeStatusRunning && curr == models.NodeStatusStopped:
		err = s.startNode(ctx)
	case target == models.NodeStatusStopped && curr == models.NodeStatusRunning:
		err = s.stopNode(ctx)
	case target == models.NodeStatusRunning && curr == models.NodeStatusRunning:
		err = s.syncNodeUsers(ctx, curr != prev)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *syncer) fetchNodeStatus(ctx context.Context) (
	curr, prev, target models.NodeStatus, err error,
) {
	// fetch stored node status
	target, prev, err = s.storage.GetNodeStatus(ctx)
	if err != nil {
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
			return
		}
	}
	return
}

func (s *syncer) markAsUnavailable(ctx context.Context) (err error) {
	return s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusUnknown)
}

const InitialRevision = models.Revision(0)

func (s *syncer) startNode(ctx context.Context) (err error) {
	// pre-edit state
	var users []models.UserSyncStatus
	if err = s.storage.DoTx(ctx, func(ctx context.Context) (err error) {
		if err = s.storage.SetNodeRev(ctx, InitialRevision); err != nil {
			return
		}
		if users, err = s.storage.FindPendingSyncs(ctx); err != nil {
			return
		}
		if err = s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusUnknown); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	// start node
	enabled := make([]models.UserProfile, 0, len(users))
	rev := InitialRevision
	for _, u := range users {
		if u.User.TargetStatus == models.UserStatusEnabled {
			enabled = append(enabled, u.User.Profile)
		}
		rev = max(rev, u.Revision)
	}

	nodeSettings, err := s.client.Start(ctx, enabled)
	if err != nil {
		return err
	}

	// post-edit state
	if err := s.storage.DoTx(ctx, func(ctx context.Context) (err error) {
		if err = s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusRunning); err != nil {
			return
		}
		if err = s.storage.SetNodeSettings(ctx, nodeSettings); err != nil {
			return
		}
		if err = s.storage.SetNodeRev(ctx, rev); err != nil {
			return
		}
		return
	}); err != nil {
		return err
	}

	return nil
}

func (s *syncer) stopNode(ctx context.Context) (err error) {
	// pre-edit state
	if err = s.storage.DoTx(ctx, func(ctx context.Context) (err error) {
		if err = s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusUnknown); err != nil {
			return
		}
		if err = s.storage.SetNodeRev(ctx, InitialRevision); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	// stop node
	err = s.client.Stop(ctx)
	if err != nil {
		return err
	}

	// post-edit state
	if err = s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusStopped); err != nil {
		return
	}

	return nil
}

func (s *syncer) syncNodeUsers(ctx context.Context, updateNodeStatus bool) (err error) {
	var users []models.UserSyncStatus
	rev := InitialRevision
	if err = s.storage.DoTx(ctx, func(ctx context.Context) (err error) {
		if users, err = s.storage.FindPendingSyncs(ctx); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	if len(users) == 0 && !updateNodeStatus {
		return nil
	}

	// edit users node
	update := models.NodeUsersUpdate{
		Add:    make([]models.UserProfile, 0, len(users)),
		Remove: make([]models.UserProfile, 0, len(users)),
	}

	for _, u := range users {
		switch u.User.TargetStatus {
		case models.UserStatusEnabled:
			update.Add = append(update.Add, u.User.Profile)
		case models.UserStatusDisabled:
			update.Remove = append(update.Remove, u.User.Profile)
		}
		rev = max(rev, u.Revision)
	}

	if err := s.client.UpdateUsers(ctx, update); err != nil {
		return err
	}

	// update node state
	if err := s.storage.DoTx(ctx, func(ctx context.Context) (err error) {
		if updateNodeStatus {
			if err = s.storage.SetCurrentNodeStatus(ctx, models.NodeStatusRunning); err != nil {
				return
			}
		}
		if len(users) > 0 {
			if err = s.storage.SetNodeRev(ctx, rev); err != nil {
				return
			}
		}
		return
	}); err != nil {
		return err
	}

	return nil
}
