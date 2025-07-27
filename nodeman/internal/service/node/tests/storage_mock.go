package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/node"
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

func (s *StorageMock) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.TargetState, s.CurrentState, nil
}

func (s *StorageMock) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	pending := make([]models.UserSyncStatus, 0, len(s.Users))
	for _, u := range s.Users {
		if u.CurrentStatus != u.TargetStatus {
			pending = append(pending, u)
		}
	}
	return pending, nil
}

func (s *StorageMock) ListManagedUsers(ctx context.Context) ([]models.UserTargetState, error) {
	users := make([]models.UserTargetState, 0, len(s.Users))
	for _, u := range s.Users {
		users = append(users, models.UserTargetState{
			User:   u.User,
			Target: u.TargetStatus,
		})
	}
	return users, nil
}

// random external operation to turn node on or off, enable or disable user
func (s *StorageMock) RandomExternalOperation() {
	if s.rand.IntN(2) == 0 {
		// switch node state
		s.TargetState = (models.NodeStatusRunning + models.NodeStatusStopped) - s.TargetState
	} else {
		// switch user state
		userIdx := s.rand.IntN(len(s.Users))
		u := s.Users[userIdx]
		u.TargetStatus = (models.UserStatusActive + models.UserStatusInactive) - u.TargetStatus
		s.Users[userIdx] = u
	}
}

func (s *StorageMock) BeginTx() node.StorageTx {
	return &MockStorageTx{parent: s}
}

type MockStorageTx struct {
	parent *StorageMock
	c      *models.NodeConfig
	s      *models.NodeStatus
	u      []models.UserStatusPatch
}

func (tx *MockStorageTx) UpdateNodeConfig(c models.NodeConfig) {
	tx.c = &c
}

// UpdateNodeStatus implements StorageTx.
func (tx *MockStorageTx) UpdateNodeStatus(s models.NodeStatus) {
	tx.s = &s
}

// UpdateNodeUsers implements StorageTx.
func (tx *MockStorageTx) UpdateNodeUsers(u []models.UserStatusPatch) {
	tx.u = append(tx.u, u...)
}

func (m *MockStorageTx) Commit(ctx context.Context) error {
	if m.s != nil {
		m.parent.CurrentState = *m.s
	}
	if m.u != nil {
		for _, u := range m.u {
			for i, pu := range m.parent.Users {
				if u.UserID == pu.User.ID {
					pu.CurrentStatus = u.Status
					m.parent.Users[i] = pu
				}
			}
		}
	}
	return nil
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

func (s *UnstableStorageMock) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}
	if s.BaseStorage.rand.Float32() < s.Instability {
		return models.NodeStatusUnknown, models.NodeStatusUnknown, fmt.Errorf("random storage fault")
	}
	return s.BaseStorage.FetchNodeStatus(ctx)
}

func (s *UnstableStorageMock) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}
	if s.BaseStorage.rand.Float32() < s.Instability {
		return nil, fmt.Errorf("random storage fault")
	}
	return s.BaseStorage.FindPendingSyncs(ctx)
}

func (s *UnstableStorageMock) ListManagedUsers(ctx context.Context) ([]models.UserTargetState, error) {
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}
	if s.BaseStorage.rand.Float32() < s.Instability {
		return nil, fmt.Errorf("random storage fault")
	}
	return s.BaseStorage.ListManagedUsers(ctx)
}

func (s *UnstableStorageMock) RandomExternalOperation() {
	s.BaseStorage.RandomExternalOperation()
}

func (s *UnstableStorageMock) BeginTx() node.StorageTx {
	return &UnstableMockStorageTx{
		MockStorageTx: MockStorageTx{parent: s.BaseStorage},
		rand:          s.BaseStorage.rand,
		Instability:   s.Instability,
	}
}

type UnstableMockStorageTx struct {
	MockStorageTx
	rand        *rand.Rand
	Instability float32
}

func (m *UnstableMockStorageTx) Commit(ctx context.Context) error {
	if m.MockStorageTx.parent.rand.Float32() < m.Instability {
		m.MockStorageTx.parent.RandomExternalOperation()
	}
	if m.parent.rand.Float32() < m.Instability {
		return fmt.Errorf("random storage fault")
	}
	return m.MockStorageTx.Commit(ctx)
}

/*type StorageMock struct {
	actualState   models.NodeStatus
	requiredState models.NodeStatus

	users []models.UserSyncStatus

	rand                     *rand.Rand
	failProb                 float32
	externalModificationProb float32

	Unstable bool
}

// BeginTx implements NodeStorage.
func (s *StorageMock) BeginTx() StorageTx {
	panic("unimplemented")
}

// FetchNodeStatus implements NodeStorage.
func (s *StorageMock) FetchNodeStatus(ctx context.Context) (target models.NodeStatus, current models.NodeStatus, err error) {
	panic("unimplemented")
}

// FindPendingSyncs implements NodeStorage.
func (s *StorageMock) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	panic("unimplemented")
}

// ListManagedUsers implements NodeStorage.
func (s *StorageMock) ListManagedUsers(ctx context.Context) ([]models.UserTargetState, error) {
	panic("unimplemented")
}

func NewStorageMock(nUsers int, log *zap.Logger) *StorageMock {
	users := make([]UserStateIntent, 0, nUsers)

	for i := range nUsers {
		state := models.UserEnabled
		if i%2 == 0 {
			state = models.UserDisabled
		}
		users = append(users, UserStateIntent{
			User:     models.User{ID: models.UserID(i)},
			Required: state,
			Actual:   models.UserStatusUnknown,
		})
	}

	return &StorageMock{
		actualState:              models.NodeStatusUnknown,
		requiredState:            models.NodeOn,
		users:                    users,
		rand:                     rand.New(rand.NewPCG(0, 0)),
		failProb:                 0.25,
		externalModificationProb: 0.25,
		Unstable:                 true,
	}
}

var _ NodeStorage = (*StorageMock)(nil)

func (s *StorageMock) GetNodeState(ctx context.Context) (
	required models.NodeStatus, actual models.NodeStatus, err error,
) {
	if err := s.applyUnstability(); err != nil {
		return models.NodeStatusUnknown, models.NodeStatusUnknown, fmt.Errorf("get node state: %w", err)
	}
	return s.requiredState, s.actualState, nil
}

// ListOutOfSyncUsers implements node.Storage.
func (s *StorageMock) ListOutOfSyncUsers(ctx context.Context) ([]UserStateIntent, error) {
	if err := s.applyUnstability(); err != nil {
		return nil, fmt.Errorf("get node state: %w", err)
	}

	oosUsers := make([]UserStateIntent, 0)
	for _, u := range s.Users {
		if u.Actual != u.Required {
			oosUsers = append(oosUsers, u)
		}
	}
	return oosUsers, nil
}

func (s *StorageMock) ListUsers(ctx context.Context) ([]UserState, error) {
	if err := s.applyUnstability(); err != nil {
		return nil, fmt.Errorf("get node state: %w", err)
	}

	l := make([]UserState, 0, len(s.Users))
	for _, u := range s.Users {
		l = append(l, UserState{
			User:   u.User,
			Status: u.Required,
		})
	}

	return l, nil
}

type UOW struct {
	parent *StorageMock
	s      models.NodeStatus
	p      *models.NodeProperties
	u      []UserStatusUpdate
}

var _ StorageWriteUOW = (*UOW)(nil)

func (uow *UOW) SetActualStatus(s models.NodeStatus) {
	uow.s = s
}

func (uow *UOW) SetActualUserStates(us []UserStatusUpdate) {
	for _, u := range us {
		uow.u = append(uow.u, u)
	}
}

// SetNodeProperties implements StorageWriteUOW.
func (s *UOW) SetNodeProperties(p models.NodeProperties) {
	s.p = &p
}

func (uow *UOW) Do(ctx context.Context) error {
	if err := uow.parent.applyUnstability(); err != nil {
		return fmt.Errorf("do uow: %w", err)
	}
	if uow.s > 0 {
		uow.parent.actualState = uow.s
	}
	for _, u := range uow.u {
		for i, su := range uow.parent.users {
			if su.User.ID == u.ID {
				su.Actual = u.Actual
				uow.parent.users[i] = su
			}
		}
	}
	return nil
}

func (s *StorageMock) GetWriteUOW() StorageWriteUOW {
	return &UOW{parent: s}
}

func (s *StorageMock) ApplyExternalModifications() {
	// external modifications
	if s.rand.Float32() < s.externalModificationProb {
		s.enableRandomUser()
	}
	if s.rand.Float32() < s.externalModificationProb {
		s.disableRandomUser()
	}
	if s.rand.Float32() < s.externalModificationProb {
		s.turnOn()
	}
	if s.rand.Float32() < s.externalModificationProb {
		s.turnOff()
	}
}

func (s *StorageMock) applyUnstability() error {
	if !s.Unstable {
		return nil
	}
	if s.rand.Float32() < s.failProb {
		return fmt.Errorf("storage mock fail")
	}
	s.ApplyExternalModifications()

	return nil
}

func (s *StorageMock) enableRandomUser() {
	idx := s.rand.IntN(len(s.Users))
	u := s.Users[idx]
	u.Required = models.UserEnabled
	s.Users[idx] = u
}

func (s *StorageMock) disableRandomUser() {
	idx := s.rand.IntN(len(s.Users))
	u := s.Users[idx]
	u.Required = models.UserDisabled
	s.Users[idx] = u
}

func (s *StorageMock) turnOn() {
	s.requiredState = models.NodeOn
}

func (s *StorageMock) turnOff() {
	s.requiredState = models.NodeOff
	for i, u := range s.Users {
		u.Actual = models.UserDisabled
		s.Users[i] = u
	}
}*/
