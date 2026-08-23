package tests

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type storage struct {
	currentStatus     models.NodeStatus
	targetStatus      models.NodeStatus
	users             []models.User
	currentUserStatus []models.UserStatus
	node              models.Node
}

/* node ops: only one node */
func (s *storage) NewNode(ctx context.Context, node *models.Node) error {
	node.ID = 0
	s.node = *node
	return nil
}

func (s *storage) SetCurrentNodeStatus(ctx context.Context, id models.NodeID, st models.NodeStatus) error {
	s.node.CurrentStatus = st
	return nil
}

func (s *storage) SetTargetNodeStatus(ctx context.Context, id models.NodeID, st models.NodeStatus) error {
	s.node.TargetStatus = st
	return nil
}

func (*storage) SetNodeSettings(ctx context.Context, id models.NodeID, s *models.NodeSettings) error {
	return nil
}

func (s *storage) GetNode(ctx context.Context, id models.NodeID) (*models.Node, error) {
	node := s.node
	return &node, nil
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
		s.currentUserStatus = append(s.currentUserStatus, models.UserStatusUnknown)
		return nil
	})
}

func (s *storage) SetTargetUserStatus(ctx context.Context, id models.UserID, st models.UserStatus) error {
	return s.do(ctx, func(s *storage) error {
		s.users[id].TargetStatus = st
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
func (s *storage) FindPendingSyncs(ctx context.Context, id models.NodeID) (
	pending []models.UserSyncStatus, err error,
) {
	pending = make([]models.UserSyncStatus, 0, len(s.users))
	err = s.do(ctx, func(s *storage) error {
		for i, u := range s.users {
			if u.TargetStatus == s.currentUserStatus[i] {
				continue
			}
			pending = append(pending, models.UserSyncStatus{
				User:          u,
				CurrentStatus: s.currentUserStatus[i],
			})
		}
		return nil
	})

	return
}

func (s *storage) SetNodeUsers(ctx context.Context, id models.NodeID, patch []models.UserStatusPatch) error {
	return s.do(ctx, func(s *storage) error {
		for i := range s.currentUserStatus {
			s.currentUserStatus[i] = models.UserStatusUnknown
		}
		for _, p := range patch {
			s.currentUserStatus[p.UserID] = p.Status
		}

		return nil
	})
}

func (s *storage) UpdateNodeUsers(ctx context.Context, id models.NodeID, patch []models.UserStatusPatch) error {
	return s.do(ctx, func(s *storage) error {
		for _, u := range patch {
			s.currentUserStatus[u.UserID] = u.Status
		}
		return nil
	})
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

func copyStorage(from *storage, to *storage) {
	to.currentStatus = from.currentStatus
	to.targetStatus = from.targetStatus
	to.users = append([]models.User{}, from.users...)
	to.currentUserStatus = append([]models.UserStatus{}, from.currentUserStatus...)
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
