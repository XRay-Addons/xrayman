package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SqlDB struct {
	*sql.DB
}

type SqlTX struct {
	*sql.Tx
}

var _ dbstorage.DB = (*SqlDB)(nil)
var _ dbstorage.TX = (*SqlTX)(nil)

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

func New(dbConn string, opts ...option) (*SqlDB, error) {
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
	db, err := initDB(dbConn, cfg)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &SqlDB{DB: db}, nil
}

func initDB(dbConn string, cfg options) (sqldb *sql.DB, err error) {
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

func (db *SqlDB) Raw() *sql.DB {
	return db.DB
}

func (db *SqlDB) BeginTx(ctx context.Context) (dbstorage.TX, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &SqlTX{Tx: tx}, nil
}

func (db *SqlDB) Close() error {
	if db == nil {
		return nil
	}
	if err := db.DB.Close(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
