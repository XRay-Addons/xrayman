package app

import (
	"context"
	"database/sql"
	"net/http"

	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/api"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/router"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/spa"
	a "github.com/XRay-Addons/xrayman/nodeman/internal/infra/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/sqldb"
	"github.com/XRay-Addons/xrayman/nodeman/internal/subscrman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/syncman"

	"go.uber.org/zap"
)

type App struct {
	app *a.App
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}

	var db *sql.DB
	var storage *dbstorage.Storage

	var httpClient *httpclient.ClientFactory
	var poolClient *client.PoolClient

	var poolSyncer poolsyncer.Syncer

	var syncJob *syncman.SyncMan

	var subscrMan subscrman.SubscrMan

	var s *service.Service

	var h *handler.Handler
	var apiHandler http.Handler
	var userHandler http.Handler
	var r http.Handler
	var httpServer *server.HttpServer

	app := a.New(
		// db
		a.WithComponent("db",
			func(ctx context.Context) (err error ) {
				db, err = sqldb.New(ctx, cfg.DBConn)
				return
			},
			func(context.Context) (err error) {
				err = sqldb.Close(db)
				return
			},
		),
		// storage
		a.WithComponent("storage",
			func(ctx context.Context) (err error) {
				storage, err = dbstorage.New(ctx, db,
					dbstorage.WithLogger(log), dbstorage.WithMigration())
				return
			}, nil,
		),

		// http client
		a.WithComponent("http client",
			func(context.Context) (err error) {
				httpClient = httpclient.NewClientFactory()
				return
			},
			func(context.Context) error {
				httpClient.Close()
				return nil
			},
		),
		// pool client
		a.WithComponent("pool client",
			func(context.Context) (err error) {
				poolClient, err = client.NewPoolClient(client.WithHTTPClient(httpClient))
				return
			}, nil,
		),

		// pool syncer
		a.WithComponent("pool syncer",
			func(context.Context) (err error) {
				poolSyncer, err = poolsyncer.New(poolClient, storage.PoolSyncStorage())
				return
			}, nil,
		),

		// background syncer
		a.WithComponent("background sync job",
			func(context.Context) (err error) {
				syncJob, err = syncman.New(poolSyncer, syncman.WithLog(log))
				return
			},
			func(ctx context.Context) (err error) {
				err = syncJob.Close()
				return
			},
		),

		// subscr manager
		a.WithComponent("subscr manager",
			func(context.Context) (err error) {
				subscrMan, err = subscrman.New(storage.SubscrmanStorage(), subscrman.WithLog(log))
				return
			}, nil,
		),

		// service
		a.WithComponent("service",
			func(context.Context) (err error) {
				s, err = service.New(poolSyncer, subscrMan, storage.ServiceStorage())
				return
			}, nil,
		),

		// handler
		a.WithComponent("handler",
			func(context.Context) (err error) {
				h, err = handler.New(s, handler.WithLogger(log))
				return
			}, nil,
		),

		// api handler
		a.WithComponent("api handler",
			func(context.Context) (err error) {
				apiHandler, err = api.NewHandler(h)
				return
			}, nil,
		),

		// user spa handler
		a.WithComponent("spa handler",
			func(context.Context) (err error) {
				userHandler, err = spa.NewHandler(cfg.UserSPAPrefix, cfg.APIPrefix)
				return
			}, nil,
		),

		// router
		a.WithComponent("router",
			func(context.Context) (err error) {
				r, err = router.New(
					router.WithHandler(cfg.APIPrefix, apiHandler),
					router.WithHandler(cfg.UserSPAPrefix, userHandler),
					router.WithLogger(log))
				return
			}, nil,
		),

		// http server
		a.WithComponent("http server",
			func(context.Context) (err error) {
				httpServer, err = server.New(cfg.Endpoint, r)
				return
			}, nil,
		),

		a.WithRunner("http server",
			func() (err error) {
				return httpServer.Listen()
			},
			func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
		),

		// logger
		a.WithLogger(log),
	)

	return &App{
		app: app,
	}, nil
}

func (app *App) Run() error {
	if app == nil {
		return errdefs.NewNilCall()
	}
	return app.app.Run()
}
