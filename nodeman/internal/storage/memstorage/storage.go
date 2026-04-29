package memstorage

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
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
func (s *Storage) PoolSyncStorage() poolsyncer.Storage {
	return &poolsyncStorage{storage: s}
}

type poolsyncStorage struct {
	storage *Storage
}

var _ poolsyncer.Storage = (*poolsyncStorage)(nil)

func (s *poolsyncStorage) DoUoW(ctx context.Context, fn poolsyncer.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
	})
}

// auth storage proxy
func (s *Storage) AuthStorage() auth.Storage {
	return &authStorage{storage: s}
}

type authStorage struct {
	storage *Storage
}

var _ auth.Storage = (*authStorage)(nil)

func (s *authStorage) DoUoW(ctx context.Context, fn auth.UoWFn) error {
	return s.storage.doLocked(ctx, func() error {
		return fn(s.storage)
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
