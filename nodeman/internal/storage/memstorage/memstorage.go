package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
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

func (s *Storage) PoolUoW() pool.UoW {
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

func (s *Storage) DoPoolSync(ctx context.Context, fn pool.UoWFn) error {
	if s == nil {
		return fmt.Errorf("storage: do pool sync: %w", errdefs.ErrNilObjectCall)
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return fn(s)
}

var _ service.UoWContext = (*Storage)(nil)
var _ pool.UoWContext = (*Storage)(nil)
var _ service.UoW = (*serviceUoW)(nil)
var _ pool.UoW = (*poolsyncUoW)(nil)

func (s *serviceUoW) Do(ctx context.Context, fn service.UoWFn) error {
	return s.storage.DoService(ctx, fn)
}

func (s *poolsyncUoW) Do(ctx context.Context, fn pool.UoWFn) error {
	return s.storage.DoPoolSync(ctx, fn)
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
	nodeUsers := make([]models.UserStatus, len(s.users))
	for i := range nodeUsers {
		nodeUsers[i] = models.UserStatusDisabled
	}
	s.syncStatus = append(s.syncStatus, nodeUsers)
	return nil
}

func (s *Storage) SetTargetNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].TargetStatus = status
	return nil
}

// UsersStorage impl
func (s *Storage) NewUser(ctx context.Context, user *models.User) error {
	user.ID = models.UserID(len(s.users))
	s.users = append(s.users, *user)
	for nodeID := range s.syncStatus {
		s.syncStatus[nodeID] = append(s.syncStatus[nodeID], models.UserStatusDisabled)
	}
	return nil
}

func (s *Storage) SetTargetUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error {
	s.users[id].TargetStatus = status
	return nil
}

func (s *Storage) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	users = append(users, s.users...)
	return users, nil
}

func (s *Storage) GetUser(ctx context.Context, id models.UserID) (*models.User, error) {
	return &s.users[id], nil
}

func (s *Storage) GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error) {
	var nodes []models.Node

	for _, node := range s.nodes {
		userNode := node.CurrentStatus == models.NodeStatusRunning &&
			node.TargetStatus == models.NodeStatusRunning &&
			s.syncStatus[node.ID][id] == models.UserStatusEnabled

		if userNode {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}
