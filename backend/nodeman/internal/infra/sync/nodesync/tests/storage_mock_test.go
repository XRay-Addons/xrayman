package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type storage struct {
	currentStatus     models.NodeStatus
	targetStatus      models.NodeStatus
	users             []models.User
	currentUserStatus []models.UserStatus

	rand        *rand.Rand
	Instability float32
}

var _ nodesync.Storage = (*storage)(nil)

func NewStorage(nUsers int) *storage {
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

	return &storage{
		currentStatus:     models.NodeStatusUnknown,
		targetStatus:      models.NodeStatusRunning,
		users:             users,
		currentUserStatus: usersStatus,
		rand:              rand.New(rand.NewPCG(0, 0)), // #nosec
	}
}

// //////////////////////////////////////////////////////////////////////////////
// getters
func (s *storage) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	err = s.do(ctx, func(s *storage) error {
		target = s.targetStatus
		current = s.currentStatus
		return nil
	})
	return
}

func (s *storage) FindPendingSyncs(ctx context.Context) (
	pending []models.UserSyncStatus, err error,
) {
	pending = make([]models.UserSyncStatus, 0, len(s.users))
	err = s.do(ctx, func(s *storage) error {
		for i, u := range s.users {
			if u.TargetStatus == s.currentUserStatus[i] {
				continue
			}
			pending = append(pending, models.UserSyncStatus{
				User:          u,
				CurrentStatus: s.currentUserStatus[i],
			})
		}
		return nil
	})

	return
}

func (s *storage) ListUsers(ctx context.Context) (
	users []models.User, err error,
) {
	err = s.do(ctx, func(s *storage) error {
		users = append(users, s.users...)
		return nil
	})
	return
}

// //////////////////////////////////////////////////////////////////////////////
// setters
func (s *storage) SetNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	return s.do(ctx, func(s *storage) error {
		for i := range s.currentUserStatus {
			s.currentUserStatus[i] = models.UserStatusDisabled
		}
		for _, p := range patch {
			s.currentUserStatus[p.UserID] = p.Status
		}

		return nil
	})
}

func (s *storage) UpdateNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	return s.do(ctx, func(s *storage) error {
		for _, u := range patch {
			s.currentUserStatus[u.UserID] = u.Status
		}
		return nil
	})
}

func (s *storage) SetNodeSettings(ctx context.Context, _ *models.NodeSettings) error {
	return nil
}

func (s *storage) SetCurrentNodeStatus(ctx context.Context, st models.NodeStatus) error {
	return s.do(ctx, func(s *storage) error {
		s.currentStatus = st
		return nil
	})
}

// //////////////////////////////////////////////////////////////////////////////
// Random external operation
func (s *storage) RandomExternalOperation() {
	switch {
	case s.rand.IntN(3) == 0:
		// switch node state
		s.targetStatus = (models.NodeStatusRunning + models.NodeStatusStopped) - s.targetStatus
	case s.rand.IntN(2) == 0:
		// switch user state
		userIdx := s.rand.IntN(len(s.users))
		u := s.users[userIdx]
		u.TargetStatus = (models.UserStatusEnabled + models.UserStatusDisabled) - u.TargetStatus
		s.users[userIdx] = u
	default:
		// add new user
		s.users = append(s.users, models.User{
			Profile: models.UserProfile{
				ID:   models.UserID(len(s.users)),
				Name: fmt.Sprintf("user %d", len(s.users)),
			},
			TargetStatus: models.UserStatusEnabled,
		})
		s.currentUserStatus = append(s.currentUserStatus, models.UserStatusUnknown)
	}
}

// //////////////////////////////////////////////////////////////////////////////
// infra
type TxFn = func(ctx context.Context) error

type txCtxKeyType struct{}

var txCtxKey = txCtxKeyType{}

func (s *storage) DoTx(ctx context.Context, fn TxFn) (err error) {
	var stx storage
	copyStorage(s, &stx)
	ctx = context.WithValue(ctx, txCtxKey, &stx)

	defer func() {
		if err != nil {
			return
		}
		copyStorage(&stx, s)
	}()

	return fn(ctx)
}

func copyStorage(from *storage, to *storage) {
	to.currentStatus = from.currentStatus
	to.targetStatus = from.targetStatus
	to.users = append([]models.User{}, from.users...)
	to.currentUserStatus = append([]models.UserStatus{}, from.currentUserStatus...)
	to.rand = rand.New(rand.NewPCG(0, 0)) // #nosec
	to.Instability = from.Instability
}

func (s *storage) do(ctx context.Context, fn func(*storage) error) error {
	f := func(s *storage) error {
		// some times this method returns error
		if s.rand.Float32() < s.Instability {
			return xerr.New("unstable storage")
		}
		// some times states changes from external
		if s.rand.Float32() < s.Instability {
			s.RandomExternalOperation()
		}
		return fn(s)
	}

	if stx, ok := ctx.Value(txCtxKey).(*storage); ok {
		return f(stx)
	}
	return f(s)
}

// //////////////////////////////////////////////////////////////////////////////

/*type patch struct {
	parent       *StorageMock
	statePatch   *models.NodeStatus
	usersPatch   []models.UserStatusPatch
	usersReplace *[]models.UserStatusPatch
}

func (s *BaseStorage) applyPatch(p *patch) error {
	if p.statePatch != nil {
		s.currentStatus = *p.statePatch
	}
	if p.usersReplace != nil {
		for userID := range s.currentUserStatus {
			s.currentUserStatus[userID] = models.UserStatusDisabled
		}
		for _, u := range *p.usersReplace {
			s.currentUserStatus[u.UserID] = u.Status
		}
	}
	if p.usersPatch != nil {
		for _, u := range p.usersPatch {
			s.currentUserStatus[u.UserID] = u.Status
		}
	}
	return nil
}

// random external operation to turn node on or off, enable or disable user
func (s *BaseStorage) randomExternalOperation() {
	switch {
	case s.rand.IntN(3) == 0:
		// switch node state
		s.targetStatus = (models.NodeStatusRunning + models.NodeStatusStopped) - s.targetStatus
	case s.rand.IntN(2) == 0:
		// switch user state
		userIdx := s.rand.IntN(len(s.users))
		u := s.users[userIdx]
		u.TargetStatus = (models.UserStatusEnabled + models.UserStatusDisabled) - u.TargetStatus
		s.users[userIdx] = u
	default:
		// add new user
		s.users = append(s.users, models.User{
			Profile: models.UserProfile{
				ID:   models.UserID(len(s.users)),
				Name: fmt.Sprintf("user %d", len(s.users)),
			},
			TargetStatus: models.UserStatusEnabled,
		})
		s.currentUserStatus = append(s.currentUserStatus, models.UserStatusUnknown)
	}
}

// simple storage mock with random extrnal operations emulation
type StorageMock struct {
	base BaseStorage
}

var _ nodesync.Storage = (*StorageMock)(nil)

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
		base: BaseStorage{
			currentStatus:     models.NodeStatusUnknown,
			targetStatus:      models.NodeStatusRunning,
			users:             users,
			currentUserStatus: usersStatus,
			rand:              rand.New(rand.NewPCG(0, 0)), // #nosec
		},
	}
}

// random external operation to turn node on or off, enable or disable user
func (s *StorageMock) RandomExternalOperation() {
	s.base.randomExternalOperation()
}

// getters
func (s *StorageMock) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.base.targetStatus, s.base.currentStatus, nil
}

func (s *StorageMock) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	pending := make([]models.UserSyncStatus, 0, len(s.base.users))
	for i, u := range s.base.users {
		if u.TargetStatus != s.base.currentUserStatus[i] {
			pending = append(pending, models.UserSyncStatus{
				User:          u,
				CurrentStatus: s.base.currentUserStatus[i],
			})
		}
	}
	return pending, nil
}

func (s *StorageMock) ListUsers(ctx context.Context) (
	[]models.User, error,
) {
	var users []models.User
	users = append(users, s.base.users...)
	return users, nil
}

// setters
func (s *StorageMock) SetClientConfig(ctx context.Context, cfg models.ClientConfigTemplate) error {
	return nil
}

func (s *StorageMock) SetCurrentNodeStatus(ctx context.Context, status models.NodeStatus) error {
	s.base.applyPatch(&patch{
		statePatch: &status,
	})
	return nil
}

func (s *StorageMock) SetNodeUsers(ctx context.Context, up []models.UserStatusPatch) error {
	s.base.applyPatch(&patch{
		usersPatch: append([]models.UserStatusPatch{}, up...),
	})
	return nil
}

func (s *StorageMock) UpdateNodeUsers(ctx context.Context, up []models.UserStatusPatch) error {
	var r []models.UserStatusPatch
	r = append(r, up...)
	s.base.applyPatch(&patch{
		usersReplace: &r,
	})
	return nil
}

func (s *StorageMock) DoTx(ctx context.Context, fn nodesync.TxFn) error {
	patch := &StorageMockTx{
		parent: s,
	}
	if err := patch.DoTx(ctx, fn); err != nil {
		return err
	}
	return nil
}

var _ nodesync.Storage = (*StorageMockTx)(nil)

func (s *StorageMockTx) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return s.parent.GetNodeStatus(ctx)
}

func (s *StorageMockTx) SetCurrentNodeStatus(ctx context.Context, nodeStatus models.NodeStatus) error {
	s.statePatch = &nodeStatus
	return nil
}

func (s *StorageMockTx) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	return s.parent.FindPendingSyncs(ctx)
}

func (s *StorageMockTx) UpdateNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	s.usersPatch = append(s.usersPatch, patch...)
	return nil
}

func (s *StorageMockTx) SetNodeUsers(ctx context.Context, patch []models.UserStatusPatch) error {
	var r []models.UserStatusPatch
	r = append(r, patch...)
	s.usersReplace = &r
	return nil
}

func (s *StorageMockTx) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.parent.ListUsers(ctx)
}

func (s *StorageMockTx) SetClientConfig(ctx context.Context, cfg models.ClientConfigTemplate) error {
	return nil
}

func (s *StorageMockTx) DoTx(ctx context.Context, fn nodesync.TxFn) error {
	if err := fn(ctx); err != nil {
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

func (s *UnstableStorageMock) DoTx(ctx context.Context, fn nodesync.TxFn) error {
	// some times this method returns error
	if s.BaseStorage.rand.Float32() < s.Instability {
		return xerr.New("unstable storage")
	}
	// some times states changes from external
	if s.BaseStorage.rand.Float32() < s.Instability {
		s.RandomExternalOperation()
	}

	uow := &StorageMockTx{
		parent: s.BaseStorage,
	}
	if err := uow.DoTx(ctx, fn); err != nil {
		return err
	}
	return nil
}

func (s *UnstableStorageMock) RandomExternalOperation() {
	s.BaseStorage.RandomExternalOperation()
}*/
