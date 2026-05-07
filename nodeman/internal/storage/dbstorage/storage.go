package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/migrations"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, db *sql.DB) (s *Storage, err error) {
	if db == nil {
		return nil, errdefs.NilArg("db")
	}

	return &Storage{
		db: db,
	}, nil
}

type option func(o *options)

type options struct {
	log *zap.Logger
}

func WithLogger(l *zap.Logger) option {
	return func(o *options) {
		if l != nil {
			o.log = l
		}
	}
}

func (s *Storage) Migrage(ctx context.Context, opts ...option) error {
	cfg := options{
		log: zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	if err := migrations.ApplyMigrations(ctx, s.db, cfg.log); err != nil {
		return translatePgErr(err)
	}

	return nil
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
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
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
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
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
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
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
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

// password storage proxy
func (s *Storage) PasswordStorage() password.Storage {
	return &passwordStorage{storage: s}
}

var _ password.Storage = (*passwordStorage)(nil)

type passwordStorage struct {
	storage *Storage
}

func (s *passwordStorage) DoUoW(ctx context.Context, fn password.UoWFn) error {
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

// main doTx impl
func (s *Storage) doTx(ctx context.Context, fn func(uowctx *uowctx) error) (err error) {
	if s == nil {
		return errdefs.NilCall()
	}
	defer func() {
		err = translatePgErr(err)
	}()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(err, xerr.WrapWithStack(rbErr))
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = xerr.WrapWithStack(commitErr)
		}
	}()

	err = fn(&uowctx{tx: tx})
	return
}
