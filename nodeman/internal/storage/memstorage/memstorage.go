package memstorage

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/service"
)

type Storage struct {
	lock sync.Locker

	nodes      []models.Node
	users      []models.User
	syncStatus [][]models.UserStatus
}

func New() *Storage {
	return &Storage{}
}

var _ service.Storage = (*Storage)(nil)

func (s *Storage) NewUoW() service.UoW {
	panic("unimplemented")
}

var _ service.UoW = (*Storage)(nil)

func (s *Storage) Do(ctx context.Context, fn service.UoWFn) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return fn(s)
}

var _ service.UoWContext = (*Storage)(nil)

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
