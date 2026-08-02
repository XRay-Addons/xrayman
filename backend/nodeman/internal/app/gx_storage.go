package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqldb"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
)

var db = gx.ProvideAnnotated(
	func(lc gx.Lifecycle, cfg *config.Config) (*sqldb.SqlDB, error) {
		db, err := sqldb.New(cfg.DBConn)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser("db", func(context.Context) error {
			return db.Close()
		})
		return db, nil
	},
	gx.As(new(dbstorage.DB)),
)

var storage = gx.ProvideAnnotated(
	dbstorage.New,
	gx.As(new(users.Storage)),
	gx.As(new(nodes.Storage)),
	gx.As(new(subscr.Storage)),
	gx.As(new(poolsync.Storage)),
	gx.As(new(poolstats.Storage)),
	gx.As(new(settings.Storage)),
	gx.As(new(auth.Storage)),
)

var Storage = gx.Module("storage",
	db,
	storage,
)
