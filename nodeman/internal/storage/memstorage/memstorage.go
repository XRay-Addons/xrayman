package memstorage

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/subscrman"
)

type Storage struct {
	lock sync.Mutex

	nodes      []models.Node
	users      []models.User
	syncStatus [][]models.UserStatus
}

type serviceStorage struct {
	storage *Storage
}

type poolsyncStorage struct {
	storage *Storage
}

type subscrmanStorage struct {
	storage *Storage
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) ServiceStorage() service.Storage {
	return &serviceStorage{storage: s}
}

func (s *Storage) PoolSyncStorage() poolsyncer.Storage {
	return &poolsyncStorage{storage: s}
}

func (s *Storage) SubscrmanStorage() subscrman.Storage {
	return &subscrmanStorage{storage: s}
}

func (s *Storage) doService(ctx context.Context, fn service.UoWFn) error {
	return s.doLocked(ctx, func() error {
		return fn(s)
	})
}

func (s *Storage) doPoolSync(ctx context.Context, fn poolsyncer.UoWFn) error {
	return s.doLocked(ctx, func() error {
		return fn(s)
	})
}

func (s *Storage) doSubscrman(ctx context.Context, fn subscrman.UoWFn) error {
	return s.doLocked(ctx, func() error {
		return fn(s)
	})
}

func (s *Storage) doLocked(ctx context.Context, fn func() error) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return fn()
}

var _ service.UoWContext = (*Storage)(nil)
var _ poolsyncer.UoWContext = (*Storage)(nil)
var _ subscrman.UoWContext = (*Storage)(nil)

var _ service.Storage = (*serviceStorage)(nil)
var _ poolsyncer.Storage = (*poolsyncStorage)(nil)
var _ subscrman.Storage = (*subscrmanStorage)(nil)

func (s *serviceStorage) DoUoW(ctx context.Context, fn service.UoWFn) error {
	return s.storage.doService(ctx, fn)
}

func (s *poolsyncStorage) DoUoW(ctx context.Context, fn poolsyncer.UoWFn) error {
	return s.storage.doPoolSync(ctx, fn)
}

func (s *subscrmanStorage) DoUoW(ctx context.Context, fn subscrman.UoWFn) error {
	return s.storage.doSubscrman(ctx, fn)
}

// NodeStatesStorage impl
func (s *Storage) ListNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	nodes = append(nodes, s.nodes...)
	return nodes, nil
}

// GetNode implements poolsyncer.UoWContext.
func (s *Storage) GetNode(ctx context.Context, id models.NodeID) (*models.Node, bool, error) {
	return &s.nodes[id], true, nil
}

func (s *Storage) SetClientConfig(ctx context.Context,
	id models.NodeID, cfg models.ClientConfigTemplate,
) error {
	s.nodes[id].Config.ClientConfigTemplate = cfg
	return nil
}

func (s *Storage) SetCurrentNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].CurrentStatus = status
	return nil
}

func (s *Storage) DeleteNode(ctx context.Context,
	id models.NodeID,
) error {
	s.nodes[id].CurrentStatus = 0;
	s.nodes[id].TargetStatus = 0
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

func (s *Storage) UpdateNodeUsers(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for _, p := range patch {
		s.syncStatus[id][p.UserID] = p.Status
	}
	return nil
}

func (s *Storage) SetNodeUsers(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for userID := range s.syncStatus[id] {
		s.syncStatus[id][userID] = models.UserStatusDisabled
	}
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
	user.Profile.ID = models.UserID(len(s.users))
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

func (s *Storage) GetUser(ctx context.Context, id models.UserID) (*models.User, bool, error) {
	return &s.users[id], true, nil
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

func (s *Storage) DeleteUser(ctx context.Context,
	id models.UserID,
) error {
	s.users[id].TargetStatus = 0;
	return nil
}