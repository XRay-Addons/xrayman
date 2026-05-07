package sqldb

import (
	"database/sql"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type option func(options *options)

type options struct {
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

func New(dbConn string, opts ...option) (sqldb *sql.DB, err error) {
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
		return nil, xerr.WrapWithStack(err)
	}

	// apply options
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.maxConnLifetime)
	db.SetConnMaxIdleTime(cfg.maxConnIdletime)

	return db, nil
}

func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	if err := db.Close(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
