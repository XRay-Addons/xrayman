package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
)

// simple storage mock with random extrnal operations emulation
type StorageMock struct {
	CurrentState models.NodeStatus
	TargetState  models.NodeStatus
	Users        []models.UserSyncStatus

	rand *rand.Rand
}

func NewStorageMock(nUsers int) *StorageMock {
	users := make([]models.UserSyncStatus, 0, nUsers)

	rnd := rand.New(rand.NewPCG(0, 0))

	for i := range nUsers {
		u := models.UserSyncStatus{
			User:          models.UserProfile{ID: models.UserID(i)},
			CurrentStatus: models.UserStatusUnknown,
			TargetStatus:  models.UserStatusInactive,
		}
		if rnd.IntN(2) == 1 {
			u.TargetStatus = models.UserStatusActive
		}
		users = append(users, u)
	}

	return &StorageMock{
		CurrentState: models.NodeStatusUnknown,
		TargetState:  models.NodeStatusRunning,
		Users:        users,
		rand:         rand.New(rand.NewPCG(0, 0)),
	}
}

// random external operation to turn node on or off, enable or disable user
func (s *StorageMock) RandomExternalOperation() {
	if s.rand.IntN(3) == 0 {
		// switch node state
		s.TargetState = (models.NodeStatusRunning + models.NodeStatusStopped) - s.TargetState
	} else if s.rand.IntN(2) == 0 {
		// switch user state
		userIdx := s.rand.IntN(len(s.Users))
		u := s.Users[userIdx]
		u.TargetStatus = (models.UserStatusActive + models.UserStatusInactive) - u.TargetStatus
		s.Users[userIdx] = u
	} else {
		// add new user
		s.Users = append(s.Users, models.UserSyncStatus{
			User:          models.UserProfile{ID: models.UserID(len(s.Users))},
			CurrentStatus: models.UserStatusUnknown,
			TargetStatus:  models.UserStatusActive,
		})
	}
}

func (s *StorageMock) fetchNodeStatus() (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.TargetState, s.CurrentState, nil
}

func (s *StorageMock) findPendingSyncs() (
	[]models.UserSyncStatus, error,
) {
	pending := make([]models.UserSyncStatus, 0, len(s.Users))
	for _, u := range s.Users {
		if u.CurrentStatus != u.TargetStatus {
			pending = append(pending, u)
		}
	}
	return pending, nil
}

func (s *StorageMock) listUsers() (
	[]models.UserTargetState, error,
) {
	users := make([]models.UserTargetState, 0, len(s.Users))
	for _, u := range s.Users {
		users = append(users, models.UserTargetState{
			User:   u.User,
			Target: u.TargetStatus,
		})
	}
	return users, nil
}

func (s *StorageMock) apply(patch *StorageMockPatch) error {
	if patch.statePatch != nil {
		s.CurrentState = *patch.statePatch
	}
	if patch.usersPatch != nil {
		for _, u := range patch.usersPatch {
			s.Users[u.UserID].CurrentStatus = u.Status
		}
	}
	return nil
}

var _ nodesyncer.Storage = (*StorageMock)(nil)

// DoUoW implements nodesyncer.Storage.
func (s *StorageMock) DoUoW(ctx context.Context, fn nodesyncer.UoWFn) error {
	uow, err := s.NewUoW()
	if err != nil {
		return fmt.Errorf("init uow: %w", err)
	}
	if err = uow.Do(ctx, fn); err != nil {
		return fmt.Errorf("do uow: %w", err)
	}
	return nil
}

func (s *StorageMock) NewUoW() (nodesyncer.UoW, error) {
	return &StorageMockPatch{
		parent: s,
	}, nil
}

type StorageMockPatch struct {
	parent     *StorageMock
	statePatch *models.NodeStatus
	usersPatch []models.UserStatusPatch
}

var _ nodesyncer.NodeConfigStorage = (*StorageMockPatch)(nil)
var _ nodesyncer.NodeStatusStorage = (*StorageMockPatch)(nil)
var _ nodesyncer.UsersStorage = (*StorageMockPatch)(nil)
var _ nodesyncer.PendingSyncsStorage = (*StorageMockPatch)(nil)
var _ nodesyncer.UoWContext = (*StorageMockPatch)(nil)
var _ nodesyncer.UoW = (*StorageMockPatch)(nil)

// nodesyncer.UoWContext impl
func (s *StorageMockPatch) NodeConfigStorage() nodesyncer.NodeConfigStorage {
	return s
}
func (s *StorageMockPatch) NodeStatusStorage() nodesyncer.NodeStatusStorage {
	return s
}

func (s *StorageMockPatch) PendingSyncsStorage() nodesyncer.PendingSyncsStorage {
	return s
}

func (s *StorageMockPatch) UsersStorage() nodesyncer.UsersStorage {
	return s
}

// nodesyncer.NodeStatusStorage impl
func (s *StorageMockPatch) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.parent.fetchNodeStatus()
}

func (s *StorageMockPatch) UpdateCurrentStatus(ctx context.Context, su models.NodeStatus) error {
	s.statePatch = &su
	return nil
}

// nodesyncer.PendingSyncsStorage impl
func (s *StorageMockPatch) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	return s.parent.findPendingSyncs()
}

func (s *StorageMockPatch) PatchPendingSyncs(ctx context.Context, patch []models.UserStatusPatch) error {
	s.usersPatch = append(s.usersPatch, patch...)
	return nil
}

// nodesyncer.UsersStorage impl
func (s *StorageMockPatch) ListUsers(ctx context.Context) ([]models.UserTargetState, error) {
	return s.parent.listUsers()
}

// nodesyncer.NodeConfigStorage impl
func (s *StorageMockPatch) UpdateClientConfig(ctx context.Context, cfg *models.ClientConfig) error {
	return nil
}

// Do implements nodesyncer.UoW.
func (s *StorageMockPatch) Do(ctx context.Context, fn nodesyncer.UoWFn) error {
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

// DoUoW implements nodesyncer.Storage.
func (s *UnstableStorageMock) DoUoW(ctx context.Context, fn nodesyncer.UoWFn) error {
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

func (s *UnstableStorageMock) NewUoW() (nodesyncer.UoW, error) {
	return &StorageMockPatch{
		parent: s.BaseStorage,
	}, nil
}

func (s *UnstableStorageMock) RandomExternalOperation() {
	s.BaseStorage.RandomExternalOperation()
}
