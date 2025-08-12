package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
)

type Storage struct {
	lock sync.Mutex

	nodes      []models.Node
	users      []models.User
	syncStatus [][]models.UserStatus
}

type serviceUoW struct {
	storage *Storage
}

type poolsyncUoW struct {
	storage *Storage
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) ServiceUoW() service.UoW {
	return &serviceUoW{storage: s}
}

func (s *Storage) PoolSyncUoW() poolsync.UoW {
	return &poolsyncUoW{storage: s}
}

func (s *Storage) DoService(ctx context.Context, fn service.UoWFn) error {
	if s == nil {
		return fmt.Errorf("storage: do service: %w", errdefs.ErrNilObjectCall)
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return fn(s)
}

func (s *Storage) DoPoolSync(ctx context.Context, fn poolsync.UoWFn) error {
	if s == nil {
		return fmt.Errorf("storage: do pool sync: %w", errdefs.ErrNilObjectCall)
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return fn(s)
}

var _ service.UoWContext = (*Storage)(nil)
var _ poolsync.UoWContext = (*Storage)(nil)
var _ service.UoW = (*serviceUoW)(nil)
var _ poolsync.UoW = (*poolsyncUoW)(nil)

func (s *serviceUoW) Do(ctx context.Context, fn service.UoWFn) error {
	return s.storage.DoService(ctx, fn)
}

func (s *poolsyncUoW) Do(ctx context.Context, fn poolsync.UoWFn) error {
	return s.storage.DoPoolSync(ctx, fn)
}

// UsersStorage impl
func (s *Storage) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	users = append(users, s.users...)
	return users, nil
}

// NodeStatesStorage impl
func (s *Storage) ListNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	nodes = append(nodes, s.nodes...)
	return nodes, nil
}

func (s *Storage) UpdateClientConfig(ctx context.Context,
	id models.NodeID, cfg models.ClientConfig,
) error {
	s.nodes[id].Config.ClientConfig = cfg
	return nil
}

func (s *Storage) FetchNodeStatus(ctx context.Context, id models.NodeID) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	node := s.nodes[id]
	return node.TargetStatus, node.CurrentStatus, nil
}

func (s *Storage) UpdateCurrentStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].CurrentStatus = status
	return nil
}

// UserSyncsStorage impl
func (s *Storage) FindPendingSyncs(ctx context.Context, id models.NodeID) (
	[]models.UserSyncStatus, error,
) {
	var syncStatus []models.UserSyncStatus
	for userId, user := range s.users {
		if user.TargetStatus == s.syncStatus[id][userId] {
			continue
		}
		syncStatus = append(syncStatus, models.UserSyncStatus{
			User:          user,
			CurrentStatus: s.syncStatus[id][userId],
		})
	}
	return syncStatus, nil
}

func (s *Storage) PatchPendingSyncs(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for _, p := range patch {
		s.syncStatus[id][p.UserID] = p.Status
	}
	return nil
}

func (s *Storage) NewNode(ctx context.Context, node *models.Node) error {
	node.ID = models.NodeID(len(s.nodes))
	s.nodes = append(s.nodes, *node)
	return nil
}

func (s *Storage) SetTargetNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].TargetStatus = status
	return nil
}
