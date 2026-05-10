package memstorage

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
)

type Storage struct {
	lock sync.Mutex

	nodes      []models.Node
	users      []models.User
	syncStatus [][]models.UserStatus

	adminID   int
	adminPass []byte
}

func New() *Storage {
	return &Storage{}
}

// nodes storage proxy
func (s *Storage) NodesStorage() nodes.Storage {
	return &nodesStorage{storage: s}
}

type nodesStorage struct {
	storage *Storage
}

var _ nodes.Storage = (*nodesStorage)(nil)

func (s *nodesStorage) DoUoW(ctx context.Context, fn nodes.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

// users storage proxy
func (s *Storage) UsersStorage() users.Storage {
	return &usersStorage{storage: s}
}

type usersStorage struct {
	storage *Storage
}

var _ users.Storage = (*usersStorage)(nil)

func (s *usersStorage) DoUoW(ctx context.Context, fn users.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

// subscr storage proxy
func (s *Storage) SubscrStorage() subscr.Storage {
	return &subscrStorage{storage: s}
}

type subscrStorage struct {
	storage *Storage
}

var _ subscr.Storage = (*subscrStorage)(nil)

func (s *subscrStorage) DoUoW(ctx context.Context, fn subscr.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

// poolsync storage proxy
func (s *Storage) PoolSyncStorage() poolsync.Storage {
	return &poolsyncStorage{storage: s}
}

type poolsyncStorage struct {
	storage *Storage
}

var _ poolsync.Storage = (*poolsyncStorage)(nil)

func (s *poolsyncStorage) DoUoW(ctx context.Context, fn poolsync.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

// password storage proxy
func (s *Storage) PasswordStorage() password.Storage {
	return &passwordStorage{storage: s}
}

type passwordStorage struct {
	storage *Storage
}

var _ password.Storage = (*passwordStorage)(nil)

func (s *passwordStorage) DoUoW(ctx context.Context, fn password.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

func (s *Storage) doLocked(ctx context.Context, fn func() error) error {
	if s == nil {
		return errdefs.NilCall()
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return fn()
}
