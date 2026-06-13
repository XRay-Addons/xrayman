package dbstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/migrations"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subheaders"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/zap"
)

type TX interface {
	Commit() error
	Rollback() error

	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB interface {
	Raw() *sql.DB
	BeginTx(ctx context.Context) (TX, error)
}

type Storage struct {
	db         DB
	ExplainLog *zap.Logger
}

func New(db DB) (s *Storage, err error) {
	if db == nil {
		return nil, errdefs.NilArg("db")
	}

	return &Storage{
		db: db,
	}, nil
}

type option func(o *options)

type options struct {
	log     *zap.Logger
	timeout time.Duration
}

func WithLogger(l *zap.Logger) option {
	return func(o *options) {
		if l != nil {
			o.log = l
		}
	}
}

func WithTimeout(t time.Duration) option {
	return func(o *options) {
		o.timeout = t
	}
}

func (s *Storage) Migrate(ctx context.Context, opts ...option) error {
	cfg := options{
		log:     zap.NewNop(),
		timeout: 5 * time.Second,
	}
	for _, o := range opts {
		o(&cfg)
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()
	if err := migrations.ApplyMigrations(ctx, s.db.Raw(), cfg.log); err != nil {
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

// subscr headers storage proxy
func (s *Storage) SubHeadersStorage() subheaders.Storage {
	return &subHeadersStorage{storage: s}
}

type subHeadersStorage struct {
	storage *Storage
}

var _ subheaders.Storage = (*subHeadersStorage)(nil)

func (s *subHeadersStorage) DoUoW(ctx context.Context, fn subheaders.UoWFn) error {
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

// stats storage
func (s *Storage) StatsStorage() poolstats.Storage {
	return &poolstatsStorage{storage: s}
}

var _ poolstats.Storage = (*poolstatsStorage)(nil)

type poolstatsStorage struct {
	storage *Storage
}

func (s *poolstatsStorage) DoUoW(ctx context.Context, fn poolstats.UoWFn) error {
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

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = xerr.Join(err, xerr.WrapWithStack(rbErr))
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = xerr.WrapWithStack(commitErr)
		}
	}()

	q := queries.New(tx)
	uowctx := uowctx{q: *q}

	err = fn(&uowctx)

	return
}
