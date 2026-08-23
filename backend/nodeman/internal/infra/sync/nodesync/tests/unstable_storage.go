package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
)

type UserStatus struct {
	User    models.UserProfile
	Target  models.UserStatus
	Current models.UserStatus
}

type StableStorage interface {
	poolsync.Storage
	nodes.Storage

	NewUser(context.Context, *models.User) error
	SetTargetNodeStatus(context.Context, models.NodeID, models.NodeStatus) error
	SetTargetUserStatus(context.Context, models.UserID, models.UserStatus) error
}

type UnstableStorage struct {
	s           StableStorage
	nodeID      models.NodeID
	rand        *rand.Rand
	Instability float32
}

var _ nodesync.Storage = (*UnstableStorage)(nil)

func WithInstability(s StableStorage, nUsers int, instability float32) (*UnstableStorage, error) {
	storage := UnstableStorage{
		s:           s,
		rand:        rand.New(rand.NewPCG(0, 0)), // #nosec
		Instability: instability,
	}
	node := models.Node{
		CurrentStatus: models.NodeStatusUnknown,
		TargetStatus:  models.NodeStatusRunning,
	}
	err := storage.s.NewNode(context.TODO(), &node)
	if err != nil {
		return nil, err
	}
	storage.nodeID = node.ID

	for i := range nUsers {
		storage.s.NewUser(context.TODO(), &models.User{
			Profile: models.UserProfile{
				Name: fmt.Sprintf("user %d", i),
			},
			TargetStatus: models.UserStatusEnabled,
		})
	}

	return &storage, nil
}

func (s *UnstableStorage) RandomExternalOperation(context.Context) error {
	return s.randomOp()
}

func (s *UnstableStorage) randomOp() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// some times this method returns error
	if s.rand.Float32() < s.Instability {
		return xerr.New("unstable storage")
	}
	// some time nothing happens
	if s.rand.Float32() >= s.Instability {
		return nil
	}
	// some times states changes from external:
	if s.rand.IntN(3) == 0 {
		// node state switched
		return s.s.DoTx(ctx, func(ctx context.Context) error {
			node, err := s.s.GetNode(ctx, s.nodeID)
			if err != nil {
				return err
			}
			target := node.TargetStatus
			target = (models.NodeStatusRunning + models.NodeStatusStopped) - target
			if err = s.s.SetTargetNodeStatus(ctx, s.nodeID, target); err != nil {
				return err
			}
			return nil
		})
	} else if s.rand.IntN(2) == 0 {
		// or user state switched
		return s.s.DoTx(ctx, func(ctx context.Context) error {
			users, err := s.s.ListUsers(ctx)
			if err != nil {
				return err
			}
			user := users[s.rand.IntN(len(users))]
			target := (models.UserStatusEnabled + models.UserStatusDisabled) - user.TargetStatus
			if err = s.s.SetTargetUserStatus(ctx, user.Profile.ID, target); err != nil {
				return err
			}
			return nil
		})
	} else {
		// or new user added
		return s.s.DoTx(ctx, func(ctx context.Context) error {
			u := models.User{
				Profile: models.UserProfile{
					Name: fmt.Sprintf("user %d", s.rand.IntN(877)),
				},
				TargetStatus: models.UserStatusEnabled,
			}
			if err := s.s.NewUser(ctx, &u); err != nil {
				return err
			}

			return nil
		})
	}
}

func (s *UnstableStorage) DoTx(ctx context.Context, fn nodesync.TxFn) error {
	if err := s.randomOp(); err != nil {
		return err
	}
	return s.s.DoTx(ctx, fn)
}

func (s *UnstableStorage) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	if err := s.randomOp(); err != nil {
		return nil, err
	}
	return s.s.FindPendingSyncs(ctx, s.nodeID)
}

func (s *UnstableStorage) SetCurrentNodeStatus(ctx context.Context, cs models.NodeStatus) error {
	if err := s.randomOp(); err != nil {
		return err
	}
	return s.s.SetCurrentNodeStatus(ctx, s.nodeID, cs)
}

func (s *UnstableStorage) SetNodeSettings(ctx context.Context, cfg *models.NodeSettings) error {
	if err := s.randomOp(); err != nil {
		return err
	}
	return s.s.SetNodeSettings(ctx, s.nodeID, cfg)
}

func (s *UnstableStorage) GetNodeStatus(ctx context.Context) (target models.NodeStatus, current models.NodeStatus, err error) {
	if err = s.randomOp(); err != nil {
		return
	}
	node, err := s.s.GetNode(ctx, s.nodeID)
	if err != nil {
		return
	}
	return node.TargetStatus, node.CurrentStatus, nil
}

func (s *UnstableStorage) ListUsers(ctx context.Context) ([]models.User, error) {
	if err := s.randomOp(); err != nil {
		return nil, err
	}
	return s.s.ListUsers(ctx)
}

func (s *UnstableStorage) SetNodeRev(ctx context.Context, rev models.Revision) (err error) {
	if err = s.randomOp(); err != nil {
		return
	}
	err = s.s.SetNodeRev(ctx, s.nodeID, rev)
	if err != nil {
		return
	}
	return
}
