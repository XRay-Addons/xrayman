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
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
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
	db      DB
	timeout time.Duration
	log     *zap.Logger
}

var _ users.Storage = (*Storage)(nil)
var _ nodes.Storage = (*Storage)(nil)
var _ subscr.Storage = (*Storage)(nil)
var _ auth.Storage = (*Storage)(nil)
var _ poolsync.Storage = (*Storage)(nil)
var _ settings.Storage = (*Storage)(nil)

type option func(o *options)

type options struct {
	timeout time.Duration
	log     *zap.Logger
}

func WithTimeout(t time.Duration) option {
	return func(o *options) {
		o.timeout = t
	}
}

func WithLogger(l *zap.Logger) option {
	return func(o *options) {
		if l != nil {
			o.log = l
		}
	}
}

func New(db DB, opts ...option) (s *Storage, err error) {
	if db == nil {
		return nil, errdefs.NilArg("db")
	}
	o := options{
		timeout: 5 * time.Second,
		log:     zap.NewNop(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &Storage{
		db:      db,
		timeout: o.timeout,
		log:     o.log,
	}, nil
}

func (s *Storage) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := migrations.ApplyMigrations(ctx, s.db.Raw(), s.log); err != nil {
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}
	return nil
}

type TxFn = func(ctx context.Context) error

type txCtxKeyType struct{}

var txCtxKey = txCtxKeyType{}

const (
	beginSavepointTx   = "SAVEPOINT sptx"
	cmSavepointTx      = "RELEASE SAVEPOINT sptx"
	rbSavepointTx      = "ROLLBACK TO SAVEPOINT sptx"
	rbSavepointTimeout = 1 * time.Second
)

func (s *Storage) DoTx(ctx context.Context, fn TxFn) (err error) {
	if tx, ok := ctx.Value(txCtxKey).(TX); ok {
		return s.doSavepointTx(ctx, fn, tx)
	}
	return s.doPureTx(ctx, fn)
}

func (s *Storage) doPureTx(ctx context.Context, fn TxFn) (err error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}
	ctx = context.WithValue(ctx, txCtxKey, tx)
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

	return fn(ctx)
}

func (s *Storage) doSavepointTx(ctx context.Context, fn TxFn, tx TX) (err error) {
	if _, err = tx.ExecContext(ctx, beginSavepointTx); err != nil {
		err = sqlerr.TranslatePgErr(err)
		err = xerr.WrapWithStack(err)
		return err
	}
	defer func() {
		if err != nil {
			rbCtx, rbCtxCancel := context.WithTimeout(context.Background(), rbSavepointTimeout)
			defer rbCtxCancel()
			if _, rbErr := tx.ExecContext(rbCtx, rbSavepointTx); rbErr != nil {
				rbErr = xerr.WrapWithStack(rbErr)
				rbErr = sqlerr.TranslatePgErr(rbErr)
				err = xerr.Join(err, rbErr)
			}
			return
		}
		if _, commitErr := tx.ExecContext(ctx, cmSavepointTx); commitErr != nil {
			commitErr = xerr.WrapWithStack(commitErr)
			commitErr = sqlerr.TranslatePgErr(commitErr)
			err = commitErr
		}
	}()

	return fn(ctx)
}

type voidFn = func(context.Context, *queries.Queries) error

func doVoid(ctx context.Context, s *Storage, fn voidFn) error {

	var q *queries.Queries

	if tx, ok := ctx.Value(txCtxKey).(TX); ok {
		q = queries.New(tx)
	} else {
		q = queries.New(s.db)
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := fn(ctx, q); err != nil {
		err = xerr.WrapWithStack(err)
		err = sqlerr.TranslatePgErr(err)
		return err
	}
	return nil
}

type anyFn[T any] = func(context.Context, *queries.Queries) (T, error)

func doAny[T any](ctx context.Context, s *Storage, fn anyFn[T]) (t T, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) (err error) {
		t, err = fn(ctx, q)
		return
	})
	return
}
