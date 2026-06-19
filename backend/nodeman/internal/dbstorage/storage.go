package dbstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/migrations"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/dynconfig"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
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

	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Storage struct {
	db         DB
	ExplainLog *zap.Logger
}

var _ users.Storage = (*Storage)(nil)
var _ nodes.Storage = (*Storage)(nil)
var _ subscr.Storage = (*Storage)(nil)
var _ auth.Storage = (*Storage)(nil)
var _ poolsync.Storage = (*Storage)(nil)
var _ dynconfig.Storage = (*Storage)(nil)

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
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}
	return nil
}

type TxFn = func(ctx context.Context) error

type txCtxKeyType struct{}

var txCtxKey = txCtxKeyType{}

func (s *Storage) DoTx(ctx context.Context, fn TxFn) (err error) {
	tx, ok := ctx.Value(txCtxKey).(TX)
	if ok {
		return xerr.New("tx in tx, Danila, are you crazy?")
	}

	tx, err = s.db.BeginTx(ctx)
	if err != nil {
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				rbErr = xerr.WrapWithStack(rbErr)
				rbErr = sqlerr.TranslatePgErr(rbErr)
				err = xerr.Join(err, rbErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			commitErr = xerr.WrapWithStack(commitErr)
			commitErr = sqlerr.TranslatePgErr(commitErr)
			err = commitErr
		}
	}()

	ctx = context.WithValue(ctx, txCtxKey, tx)
	return fn(ctx)
}

type voidFn = func(context.Context, *queries.Queries) error

func doVoid(ctx context.Context, s *Storage, fn voidFn) error {
	q := queries.New(s.db)
	if tx, ok := ctx.Value(txCtxKey).(TX); ok {
		q = queries.New(tx)
	}
	if err := fn(ctx, q); err != nil {
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}
	return nil
}

type anyFn[T any] = func(context.Context, *queries.Queries) (T, error)

func doAny[T any](ctx context.Context, s *Storage, fn anyFn[T]) (t T, err error) {
	err = doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) (err error) {
		t, err = fn(ctx, q)
		return
	})
	return
}
