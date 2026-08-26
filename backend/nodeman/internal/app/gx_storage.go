package app

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqldb"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/metrics"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var db = gx.ProvideAnnotated(
	func(lc gx.Lifecycle, cfg *config.Config) (*sqldb.SqlDB, error) {
		db, err := sqldb.New(cfg.DBConn)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser(gx.Closer{
			Name: "db",
			OnClose: func(context.Context) error {
				return db.Close()
			},
		})
		return db, nil
	},
	gx.As(new(dbstorage.DB)),
)

type StorageParams struct {
	fx.In
	DB      dbstorage.DB
	Timeout time.Duration `name:"storage-call-timeout"`
	Log     *zap.Logger
}

var storage = gx.ProvideAnnotated(
	func(p StorageParams) (*dbstorage.Storage, error) {
		return dbstorage.New(p.DB,
			dbstorage.WithTimeout(p.Timeout),
			dbstorage.WithLogger(p.Log))
	},
	gx.As(new(users.Storage)),
	gx.As(new(nodes.Storage)),
	gx.As(new(subscr.Storage)),
	gx.As(new(poolsync.Storage)),
	gx.As(new(poolstats.Storage)),
	gx.As(new(settings.Storage)),
	gx.As(new(auth.Storage)),
	gx.As(new(metrics.Storage)),
	gx.As(gx.Self()),
)

var Storage = gx.Module("storage",
	db,
	storage,
)
