package dbstorage

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/db"
	"github.com/XRay-Addons/xrayman/nodeman/internal/subscrman"
)

type Storage struct {
	db *db.SQLDB
}

type serviceStorage struct {
	storage *Storage
}

type poolsyncStorage struct {
	storage *Storage
}

type subscrmanStorage struct {
	storage *Storage
}

func New(dbConn string) (*Storage, error) {
	db, err := db.NewSQLDB(dbConn, db.WithMigration())
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) ServiceStorage() service.Storage {
	return &serviceStorage{storage: s}
}

func (s *Storage) PoolSyncStorage() poolsyncer.Storage {
	return &poolsyncStorage{storage: s}
}

func (s *Storage) SubscrmanStorage() subscrman.Storage {
	return &subscrmanStorage{storage: s}
}

func (s *Storage) doService(ctx context.Context, fn service.UoWFn) error {
	return s.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

func (s *Storage) doPoolSync(ctx context.Context, fn poolsyncer.UoWFn) error {
	return s.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

func (s *Storage) doSubscrman(ctx context.Context, fn subscrman.UoWFn) error {
	return s.doTx(ctx, func(uowctx *uowctx) error {
		return fn(uowctx)
	})
}

func (s *Storage) doTx(ctx context.Context, fn func(uowctx *uowctx) error) (err error) {
	if s == nil {
		return errdefs.NewNilCall()
	}

	tx, err := s.db.BeginTx(ctx)
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

var _ service.Storage = (*serviceStorage)(nil)

func (s *serviceStorage) DoUoW(ctx context.Context, fn service.UoWFn) error {
	return s.storage.doService(ctx, fn)
}

var _ poolsyncer.Storage = (*poolsyncStorage)(nil)

func (s *poolsyncStorage) DoUoW(ctx context.Context, fn poolsyncer.UoWFn) error {
	return s.storage.doPoolSync(ctx, fn)
}

var _ subscrman.Storage = (*subscrmanStorage)(nil)

func (s *subscrmanStorage) DoUoW(ctx context.Context, fn subscrman.UoWFn) error {
	return s.storage.doSubscrman(ctx, fn)
}
