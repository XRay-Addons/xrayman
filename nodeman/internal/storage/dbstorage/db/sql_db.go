package db

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLDB struct {
	db            *sql.DB
	migrate       atomic.Bool
	migrationLock sync.Mutex
}

type option func(options *options)

func WithMigration() option {
	return func(options *options) {
		options.migrate = true
	}
}

type options struct {
	migrate         bool
	maxOpenConns    int
	maxIdleConns    int
	maxConnLifetime time.Duration
	maxConnIdletime time.Duration
}

const (
	defaultMaxOpenConns    = 16
	defaultMaxIdleConns    = 8
	defaultMaxConnLifetime = 30 * time.Minute
	defaultMaxConnIdletime = 300 * time.Minute
)

func NewSQLDB(dbConn string, opts ...option) (*SQLDB, error) {
	cfg := options{
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		maxConnLifetime: defaultMaxConnLifetime,
		maxConnIdletime: defaultMaxConnIdletime,
	}
	for _, o := range opts {
		o(&cfg)
	}

	// sql.open not actually conntects, just check dbConn string
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	// apply options
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.maxConnLifetime)
	db.SetConnMaxIdleTime(cfg.maxConnIdletime)

	sqldb := &SQLDB{
		db: db,
	}
	// set up migration flag
	sqldb.migrate.Store(cfg.migrate)

	return sqldb, nil
}

func (d *SQLDB) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	if err := d.db.Close(); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}

func (d *SQLDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if d == nil || d.db == nil {
		return nil, errdefs.NewNilCall()
	}

	// apply migration if required before every tx till success
	if err := d.applyMigration(ctx); err != nil {
		return nil, err
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	return tx, nil
}

func (d *SQLDB) applyMigration(ctx context.Context) error {
	if !d.migrate.Load() {
		return nil
	}

	d.migrationLock.Lock()
	defer d.migrationLock.Unlock()
	if !d.migrate.Load() {
		return nil
	}
	if err := migrate(ctx, d.db); err != nil {
		return err
	}
	d.migrate.Store(false)

	return nil
}
