package node

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/models"
	"go.uber.org/zap"
)

type StorageMock struct {
	actualState   models.NodeStatus
	requiredState models.NodeStatus

	users []OutOfSyncUser

	rand                     *rand.Rand
	failProb                 float32
	externalModificationProb float32

	Unstable bool
}

func NewStorageMock(nUsers int, log *zap.Logger) *StorageMock {
	users := make([]OutOfSyncUser, 0, nUsers)

	for i := range nUsers {
		state := models.UserEnabled
		if i%2 == 0 {
			state = models.UserDisabled
		}
		users = append(users, OutOfSyncUser{
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

var _ Storage = (*StorageMock)(nil)

func (s *StorageMock) GetNodeState(ctx context.Context) (
	required models.NodeStatus, actual models.NodeStatus, err error,
) {
	if err := s.applyUnstability(); err != nil {
		return models.NodeStatusUnknown, models.NodeStatusUnknown, fmt.Errorf("get node state: %w", err)
	}
	return s.requiredState, s.actualState, nil
}

// ListOutOfSyncUsers implements node.Storage.
func (s *StorageMock) ListOutOfSyncUsers(ctx context.Context) ([]OutOfSyncUser, error) {
	if err := s.applyUnstability(); err != nil {
		return nil, fmt.Errorf("get node state: %w", err)
	}

	oosUsers := make([]OutOfSyncUser, 0)
	for _, u := range s.users {
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

	l := make([]UserState, 0, len(s.users))
	for _, u := range s.users {
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
	idx := s.rand.IntN(len(s.users))
	u := s.users[idx]
	u.Required = models.UserEnabled
	s.users[idx] = u
}

func (s *StorageMock) disableRandomUser() {
	idx := s.rand.IntN(len(s.users))
	u := s.users[idx]
	u.Required = models.UserDisabled
	s.users[idx] = u
}

func (s *StorageMock) turnOn() {
	s.requiredState = models.NodeOn
	/*for i, u := range s.users {
		u.Required = models.UserDisabled
		s.users[i] = u
	}*/
}

func (s *StorageMock) turnOff() {
	s.requiredState = models.NodeOff
	for i, u := range s.users {
		u.Actual = models.UserDisabled
		s.users[i] = u
	}
}
