package memstorage

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
)

type Storage struct {
	lock sync.Locker

	nodes         []models.NodeConfig
	currentStatus []models.NodeStatus
	targetStatus  []models.NodeStatus

	users      []models.UserProfile
	userStatus []models.UserStatus

	syncStatus [][]models.UserStatus
}

func NewMemStorage() *Storage {
	return &Storage{}
}

var _ service.Storage = (*Storage)(nil)

func (s *Storage) DoUoW(ctx context.Context, fn service.UoWFn) error {
	s.lock.Lock()
	defer s.lock.Lock()

	return fn(s)
}

func (s *Storage) NewUoW() (service.UoW, error) {
	return s, nil
}

var _ service.UoW = (*Storage)(nil)

func (s *Storage) Do(ctx context.Context, fn service.UoWFn) error {
	return s.DoUoW(ctx, fn)
}

var _ service.UoWContext = (*Storage)(nil)

func (s *Storage) NodeConfigStorage() service.NodeConfigStorage {
	return s
}

func (s *Storage) NodeStatusStorage() service.NodeStatusStorage {
	return s
}

func (s *Storage) PendingSyncsStorage() service.PendingSyncsStorage {
	return s
}

func (s *Storage) UserStatusStorage() service.UserStatusStorage {
	return s
}

func (s *Storage) UsersStorage() service.UsersStorage {
	return s
}

var _ service.NodeConfigStorage = (*Storage)(nil)

func (s *Storage) AddNode(ctx context.Context, node *models.NodeConfig) error {
	node.ID = models.NodeID(len(s.nodes))
	s.nodes = append(s.nodes, *node)
	s.currentStatus = append(s.currentStatus, models.NodeStatusUnknown)
	s.targetStatus = append(s.targetStatus, models.NodeStatusStopped)
	return nil
}

func (s *Storage) ListNodes(ctx context.Context) ([]models.NodeConfig, error) {
	nodes := make([]models.NodeConfig, 0, len(s.nodes))
	nodes = append(nodes, s.nodes...)
	return nodes, nil
}

func (s *Storage) UpdateConnectionInfo(ctx context.Context,
	id models.NodeID, connInfo *models.NodeConnectionInfo,
) error {
	s.nodes[id].ConnectionInfo = *connInfo
	return nil
}

func (s *Storage) GetConnectionInfo(ctx context.Context, id models.NodeID) (
	*models.NodeConnectionInfo, error,
) {
	return &s.nodes[id].ConnectionInfo, nil
}

func (s *Storage) GetClientConfig(ctx context.Context, id models.NodeID) (
	*models.ClientConfig, error,
) {
	return &s.nodes[id].ClientConfig, nil
}

func (s *Storage) UpdateClientConfig(ctx context.Context,
	id models.NodeID, cfg *models.ClientConfig,
) error {
	s.nodes[id].ClientConfig = *cfg
	return nil
}

var _ service.NodeStatusStorage = (*Storage)(nil)

func (s *Storage) FetchNodeStatus(ctx context.Context, id models.NodeID) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.targetStatus[id], s.currentStatus[id], nil
}

func (s *Storage) UpdateCurrentStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.currentStatus[id] = status
	return nil
}

func (s *Storage) UpdateTargetStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.targetStatus[id] = status
	return nil
}

var _ service.PendingSyncsStorage = (*Storage)(nil)

func (s *Storage) FindPendingSyncs(ctx context.Context, id models.NodeID) (
	[]models.UserSyncStatus, error,
) {
	syncs := make([]models.UserSyncStatus, 0, len(s.users))
	for i, u := range s.users {
		if s.userStatus[i] != s.syncStatus[id][i] {
			syncs = append(syncs, models.UserSyncStatus{
				User:          u,
				TargetStatus:  s.userStatus[i],
				CurrentStatus: s.syncStatus[id][i],
			})
		}
	}
	return syncs, nil
}

func (s *Storage) PatchPendingSyncs(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for _, u := range patch {
		s.syncStatus[id][u.UserID] = u.Status
	}
	return nil
}

var _ service.UserStatusStorage = (*Storage)(nil)

func (s *Storage) GetUserStatus(ctx context.Context, id models.UserID) (
	models.UserStatus, error,
) {
	return s.userStatus[id], nil
}

func (s *Storage) SetUserStatus(ctx context.Context, id models.UserID,
	status models.UserStatus,
) error {
	s.userStatus[id] = status
	return nil
}

var _ service.UsersStorage = (*Storage)(nil)

func (s *Storage) AddUser(ctx context.Context, user *models.UserProfile) error {
	user.ID = models.UserID(len(s.users))
	s.users = append(s.users, *user)
	s.userStatus = append(s.userStatus, models.UserStatusInactive)
	for i := range s.userStatus {
		s.syncStatus[i] = append(s.syncStatus[i], models.UserStatusInactive)
	}
	return nil
}

func (s *Storage) ListUsers(ctx context.Context) ([]models.UserTargetState, error) {
	users := make([]models.UserTargetState, 0, len(s.users))
	for i, u := range s.users {
		users = append(users, models.UserTargetState{
			User:   u,
			Target: s.userStatus[i],
		})
	}
	return users, nil
}
