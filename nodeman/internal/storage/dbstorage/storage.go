package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/migrations"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
}

type option func(o *options)

type options struct {
	migrate bool
	log     *zap.Logger
}

func WithMigration() option {
	return func(o *options) {
		o.migrate = true
	}
}

func WithLogger(l *zap.Logger) option {
	return func(o *options) {
		if l != nil {
			o.log = l
		}
	}
}

func New(ctx context.Context, db *sql.DB, opts ...option) (s *Storage, err error) {
	if db == nil {
		return nil, errdefs.NewNilArg("db")
	}

	cfg := options{
		migrate: false,
		log:     zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.migrate {
		if err = applyMigrations(ctx, db, cfg.log); err != nil {
			return
		}
	}

	return &Storage{
		db: db,
	}, nil
}

func applyMigrations(ctx context.Context, db *sql.DB, log *zap.Logger) error {
	const retryInterval = 100 * time.Millisecond
	b := retry.NewConstant(retryInterval)

	return retry.Do(ctx, b, func(ctx context.Context) error {
		err := migrations.ApplyMigrations(ctx, db, log)
		if err == nil {
			log.Warn("migration successed")
			return nil
		}
		err = translatePgErr(err)
		if errors.Is(err, errdefs.ErrTemporaryUnavailable) {
			log.Warn("migrations retryable fail", zap.Error(err))
			return retry.RetryableError(err)
		}
		log.Warn("migrations unretryable fail", zap.Error(err))
		return nil
	})
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

// auth storage proxy
func (s *Storage) AuthStorage() auth.Storage {
	return &authStorage{storage: s}
}

var _ auth.Storage = (*authStorage)(nil)

type authStorage struct {
	storage *Storage
}

func (s *authStorage) DoUoW(ctx context.Context, fn auth.UoWFn) error {
	return s.storage.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

// main doTx impl
func (s *Storage) doTx(ctx context.Context, fn func(uowctx *uowctx) error) (err error) {
	if s == nil {
		return errdefs.NewNilCall()
	}
	defer func() {
		err = translatePgErr(err)
	}()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errdefs.WrapWithStack(err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(err, errdefs.WrapWithStack(rbErr))
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = errdefs.WrapWithStack(commitErr)
		}
	}()

	err = fn(&uowctx{tx: tx})
	return
}
