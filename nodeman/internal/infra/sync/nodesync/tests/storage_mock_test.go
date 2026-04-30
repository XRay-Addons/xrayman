package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// simple storage mock with random extrnal operations emulation
type StorageMock struct {
	CurrentStatus     models.NodeStatus
	TargetStatus      models.NodeStatus
	Users             []models.User
	CurrentUserStatus []models.UserStatus

	rand *rand.Rand
}

func NewStorageMock(nUsers int) *StorageMock {
	users := make([]models.User, 0, nUsers)
	usersStatus := make([]models.UserStatus, 0, nUsers)
	rnd := rand.New(rand.NewPCG(0, 0)) // #nosec

	for i := range nUsers {
		u := models.User{
			Profile: models.UserProfile{
				ID:   models.UserID(i),
				Name: fmt.Sprintf("user %d", i),
			},
			TargetStatus: models.UserStatusDisabled,
		}
		if rnd.IntN(2) == 1 {
			u.TargetStatus = models.UserStatusEnabled
		}
		users = append(users, u)
		usersStatus = append(usersStatus, models.UserStatusUnknown)
	}

	return &StorageMock{
		CurrentStatus:     models.NodeStatusUnknown,
		TargetStatus:      models.NodeStatusRunning,
		Users:             users,
		CurrentUserStatus: usersStatus,
		rand:              rand.New(rand.NewPCG(0, 0)), // #nosec
	}
}

// random external operation to turn node on or off, enable or disable user
func (s *StorageMock) RandomExternalOperation() {
	switch {
	case s.rand.IntN(3) == 0:
		// switch node state
		s.TargetStatus = (models.NodeStatusRunning + models.NodeStatusStopped) - s.TargetStatus
	case s.rand.IntN(2) == 0:
		// switch user state
		userIdx := s.rand.IntN(len(s.Users))
		u := s.Users[userIdx]
		u.TargetStatus = (models.UserStatusEnabled + models.UserStatusDisabled) - u.TargetStatus
		s.Users[userIdx] = u
	default:
		// add new user
		s.Users = append(s.Users, models.User{
			Profile: models.UserProfile{
				ID:   models.UserID(len(s.Users)),
				Name: fmt.Sprintf("user %d", len(s.Users)),
			},
			TargetStatus: models.UserStatusEnabled,
		})
		s.CurrentUserStatus = append(s.CurrentUserStatus, models.UserStatusUnknown)
	}
}

func (s *StorageMock) fetchNodeStatus() (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.TargetStatus, s.CurrentStatus, nil
}

func (s *StorageMock) findPendingSyncs() (
	[]models.UserSyncStatus, error,
) {
	pending := make([]models.UserSyncStatus, 0, len(s.Users))
	for i, u := range s.Users {
		if u.TargetStatus != s.CurrentUserStatus[i] {
			pending = append(pending, models.UserSyncStatus{
				User:          u,
				CurrentStatus: s.CurrentUserStatus[i],
			})
		}
	}
	return pending, nil
}

func (s *StorageMock) listUsers() (
	[]models.User, error,
) {
	var users []models.User
	users = append(users, s.Users...)
	return users, nil
}

func (s *StorageMock) apply(patch *StorageMockPatch) error {
	if patch.statePatch != nil {
		s.CurrentStatus = *patch.statePatch
	}
	if patch.usersReplace != nil {
		for userID := range s.CurrentUserStatus {
			s.CurrentUserStatus[userID] = models.UserStatusDisabled
		}
		for _, u := range *patch.usersReplace {
			s.CurrentUserStatus[u.UserID] = u.Status
		}
	}
	if patch.usersPatch != nil {
		for _, u := range patch.usersPatch {
			s.CurrentUserStatus[u.UserID] = u.Status
		}
	}
	return nil
}

var _ nodesync.Storage = (*StorageMock)(nil)

func (s *StorageMock) DoUoW(ctx context.Context, fn nodesync.UoWFn) error {
	uow := &StorageMockPatch{
		parent: s,
	}
	if err := uow.Do(ctx, fn); err != nil {
		return err
	}
	return nil
}

type StorageMockPatch struct {
	parent       *StorageMock
	statePatch   *models.NodeStatus
	usersPatch   []models.UserStatusPatch
	usersReplace *[]models.UserStatusPatch
}

var _ nodesync.UoWContext = (*StorageMockPatch)(nil)

func (s *StorageMockPatch) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.parent.fetchNodeStatus()
}

func (s *StorageMockPatch) SetCurrentNodeStatus(ctx context.Context, nodeStatus models.NodeStatus) error {
	s.statePatch = &nodeStatus
	return nil
}

func (s *StorageMockPatch) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	return s.parent.findPendingSyncs()
}

func (s *StorageMockPatch) UpdateNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	s.usersPatch = append(s.usersPatch, patch...)
	return nil
}

func (s *StorageMockPatch) SetNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	var r []models.UserStatusPatch
	r = append(r, patch...)
	s.usersReplace = &r
	return nil
}

func (s *StorageMockPatch) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.parent.listUsers()
}

func (s *StorageMockPatch) SetClientConfig(ctx context.Context, cfg models.ClientConfigTemplate) error {
	return nil
}

func (s *StorageMockPatch) Do(ctx context.Context, fn nodesync.UoWFn) error {
	if err := fn(s); err != nil {
		return err
	}
	err := s.parent.apply(s)
	return err
}

// storage mock with external faults or edit state modifications
type UnstableStorageMock struct {
	BaseStorage *StorageMock
	Instability float32
}

func NewUnstableStorageMock(nUsers int) *UnstableStorageMock {
	return &UnstableStorageMock{
		BaseStorage: NewStorageMock(nUsers),
	}
}

func (s *UnstableStorageMock) DoUoW(ctx context.Context, fn nodesync.UoWFn) error {
	// some times this method returns error
	if s.BaseStorage.rand.Float32() < s.Instability {
		return errdefs.New("unstable storage")
	}
	// some times states changes from external
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}

	uow := &StorageMockPatch{
		parent: s.BaseStorage,
	}
	if err := uow.Do(ctx, fn); err != nil {
		return err
	}
	return nil
}

func (s *UnstableStorageMock) RandomExternalOperation() {
	s.BaseStorage.RandomExternalOperation()
}
