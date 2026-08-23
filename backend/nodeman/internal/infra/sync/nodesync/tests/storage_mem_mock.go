package tests

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type storage struct {
	globalRev models.Revision

	users    []models.User
	usersRev []models.Revision

	currentStatus models.NodeStatus
	targetStatus  models.NodeStatus
	node          models.Node
	nodeRev       models.Revision
}

/* node ops: only one node */
func (s *storage) NewNode(ctx context.Context, node *models.Node) error {
	node.ID = 0
	s.node = *node
	return nil
}

func (s *storage) SetCurrentNodeStatus(ctx context.Context, id models.NodeID, st models.NodeStatus) error {
	return s.do(ctx, func(s *storage) error {
		s.node.CurrentStatus = st
		return nil
	})
}

func (s *storage) SetTargetNodeStatus(ctx context.Context, id models.NodeID, st models.NodeStatus) error {
	return s.do(ctx, func(s *storage) error {
		s.node.TargetStatus = st
		return nil
	})
}

func (*storage) SetNodeSettings(ctx context.Context, id models.NodeID, s *models.NodeSettings) error {
	return nil
}

func (s *storage) GetNode(ctx context.Context, id models.NodeID) (node *models.Node, err error) {
	err = s.do(ctx, func(s *storage) error {
		node = &s.node
		return nil
	})
	return
}

func (s *storage) ListNodes(ctx context.Context) ([]models.Node, error) {
	panic("unimplemented")
}

func (s *storage) DeleteNode(ctx context.Context, id models.NodeID) error {
	panic("unimplemented")
}

/* user ops */
func (s *storage) NewUser(ctx context.Context, u *models.User) error {
	return s.do(ctx, func(s *storage) error {
		u.Profile.ID = len(s.users)
		s.users = append(s.users, *u)
		s.globalRev++
		s.usersRev = append(s.usersRev, s.globalRev)
		return nil
	})
}

func (s *storage) SetTargetUserStatus(ctx context.Context, id models.UserID, st models.UserStatus) error {
	return s.do(ctx, func(s *storage) error {
		s.users[id].TargetStatus = st
		s.globalRev++
		s.usersRev[id] = s.globalRev
		return nil
	})
}

func (s *storage) ListUsers(ctx context.Context) (users []models.User, err error) {
	err = s.do(ctx, func(s *storage) error {
		users = append(users, s.users...)
		return nil
	})
	return
}

func (s *storage) DeleteUser(ctx context.Context, id models.UserID) error {
	panic("unimplemented")
}

/* usernode ops */
func (s *storage) SetNodeRev(ctx context.Context, id models.NodeID, rev models.Revision) (err error) {
	err = s.do(ctx, func(s *storage) error {
		s.nodeRev = rev
		return nil
	})
	return
}

func (s *storage) FindPendingSyncs(ctx context.Context, id models.NodeID) (
	pending []models.UserSyncStatus, err error,
) {
	pending = make([]models.UserSyncStatus, 0, len(s.users))
	err = s.do(ctx, func(s *storage) error {
		for i, u := range s.users {
			if s.usersRev[i] <= s.nodeRev {
				continue
			}
			pending = append(pending, models.UserSyncStatus{
				User:     u,
				Revision: s.usersRev[i],
			})
		}
		return nil
	})

	return
}

var _ StableStorage = (*storage)(nil)

func NewMemoryMockStorage() *storage {
	return &storage{
		currentStatus: models.NodeStatusUnknown,
		targetStatus:  models.NodeStatusRunning,
	}
}

// //////////////////////////////////////////////////////////////////////////////
// infra
type TxFn = func(ctx context.Context) error

type txCtxKeyType struct{}

var txCtxKey = txCtxKeyType{}

func (s *storage) DoTx(ctx context.Context, fn TxFn) (err error) {
	var stx storage
	copyStorage(s, &stx)
	ctx = context.WithValue(ctx, txCtxKey, &stx)

	defer func() {
		if err != nil {
			return
		}
		copyStorage(&stx, s)
	}()

	return fn(ctx)
}

/*
globalRev models.Revision

users    []models.User
usersRev []models.Revision

currentStatus models.NodeStatus
targetStatus  models.NodeStatus
node          models.Node
nodeRev       models.Revision
*/
func copyStorage(from *storage, to *storage) {
	to.globalRev = from.globalRev

	to.users = append([]models.User{}, from.users...)
	to.usersRev = append([]models.Revision{}, from.usersRev...)

	to.currentStatus = from.currentStatus
	to.targetStatus = from.targetStatus
	to.node = from.node
	to.nodeRev = from.nodeRev
}

func (s *storage) do(ctx context.Context, fn func(*storage) error) error {
	f := func(s *storage) error {
		return fn(s)
	}

	if stx, ok := ctx.Value(txCtxKey).(*storage); ok {
		return f(stx)
	}
	return f(s)
}
