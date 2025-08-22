package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
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
	rnd := rand.New(rand.NewPCG(0, 0))

	for i := range nUsers {
		u := models.User{
			ID:           models.UserID(i),
			Profile:      models.UserProfile{Name: fmt.Sprintf("user %d", i)},
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
		rand:              rand.New(rand.NewPCG(0, 0)),
	}
}

// random external operation to turn node on or off, enable or disable user
func (s *StorageMock) RandomExternalOperation() {
	if s.rand.IntN(3) == 0 {
		// switch node state
		s.TargetStatus = (models.NodeStatusRunning + models.NodeStatusStopped) - s.TargetStatus
	} else if s.rand.IntN(2) == 0 {
		// switch user state
		userIdx := s.rand.IntN(len(s.Users))
		u := s.Users[userIdx]
		u.TargetStatus = (models.UserStatusEnabled + models.UserStatusDisabled) - u.TargetStatus
		s.Users[userIdx] = u
	} else {
		// add new user
		s.Users = append(s.Users, models.User{
			ID:           models.UserID(len(s.Users)),
			Profile:      models.UserProfile{Name: fmt.Sprintf("user %d", len(s.Users))},
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
	if patch.usersPatch != nil {
		for _, u := range patch.usersPatch {
			s.CurrentUserStatus[u.UserID] = u.Status
		}
	}
	return nil
}

var _ pool.NodeUoW = (*StorageMock)(nil)

func (s *StorageMock) Do(ctx context.Context, fn pool.NodeUoWFn) error {
	uow, err := s.NewUoW()
	if err != nil {
		return fmt.Errorf("init uow: %w", err)
	}
	if err = uow.Do(ctx, fn); err != nil {
		return fmt.Errorf("do uow: %w", err)
	}
	return nil
}

func (s *StorageMock) NewUoW() (pool.NodeUoW, error) {
	return &StorageMockPatch{
		parent: s,
	}, nil
}

type StorageMockPatch struct {
	parent     *StorageMock
	statePatch *models.NodeStatus
	usersPatch []models.UserStatusPatch
}

var _ pool.NodeUoWContext = (*StorageMockPatch)(nil)
var _ pool.NodeUoW = (*StorageMockPatch)(nil)

func (s *StorageMockPatch) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.parent.fetchNodeStatus()
}

func (s *StorageMockPatch) UpdateCurrentStatus(ctx context.Context, su models.NodeStatus) error {
	s.statePatch = &su
	return nil
}

func (s *StorageMockPatch) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	return s.parent.findPendingSyncs()
}

func (s *StorageMockPatch) PatchPendingSyncs(ctx context.Context, patch []models.UserStatusPatch) error {
	s.usersPatch = append(s.usersPatch, patch...)
	return nil
}

func (s *StorageMockPatch) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.parent.listUsers()
}

func (s *StorageMockPatch) UpdateClientConfig(ctx context.Context, cfg models.ClientConfig) error {
	return nil
}

func (s *StorageMockPatch) Do(ctx context.Context, fn pool.NodeUoWFn) error {
	if err := fn(s); err != nil {
		return fmt.Errorf("patch cfg error: %w", err)
	}
	return s.parent.apply(s)
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

func (s *UnstableStorageMock) DoUoW(ctx context.Context, fn pool.NodeUoWFn) error {
	// some times this method returns error
	if s.BaseStorage.rand.Float32() < s.Instability {
		return fmt.Errorf("unstable storage")
	}
	// some times states changes from external
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}

	uow, err := s.BaseStorage.NewUoW()
	if err != nil {
		return fmt.Errorf("init uow: %w", err)
	}
	if err = uow.Do(ctx, fn); err != nil {
		return fmt.Errorf("do uow: %w", err)
	}
	return nil
}

func (s *UnstableStorageMock) NewUoW() (pool.NodeUoW, error) {
	return &StorageMockPatch{
		parent: s.BaseStorage,
	}, nil
}

func (s *UnstableStorageMock) RandomExternalOperation() {
	s.BaseStorage.RandomExternalOperation()
}
