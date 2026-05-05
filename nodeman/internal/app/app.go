package app

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/XRay-Addons/xrayman/nodeman/internal/app/bootstrap"
	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/api"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/router"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/jwt"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	a "github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	nodes "github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	subscr "github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	users "github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/sqldb"

	"go.uber.org/zap"
)

type App struct {
	app *a.App
}

//go:embed userpage/**
var userpageFS embed.FS

//go:embed admpage/**
var admpageFS embed.FS

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}

	var db *sql.DB
	var storage *dbstorage.Storage

	var httpClient *httpclient.ClientFactory
	var poolClient *client.PoolClient

	var pwd *password.Password
	var authJWT *jwt.JWT

	var poolSyncer poolsync.Syncer

	var syncJob *syncman.SyncMan

	var authService *auth.Service
	var usersService *users.Service
	var nodesService *nodes.Service
	var subscrService *subscr.Service

	var h *handler.Handler
	var s *security.Handler
	var apiHandler http.Handler
	var userpageSpa fs.FS
	var admpageSpa fs.FS

	var r http.Handler
	var httpServer *server.HttpServer

	app := a.New(
		// db
		a.WithComponent("db",
			func(ctx context.Context) (err error) {
				db, err = sqldb.New(cfg.DBConn)
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

		// nodes http client
		a.WithComponent("nodes http client",
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
				poolSyncer, err = poolsync.New(poolClient, storage.PoolSyncStorage())
				return
			}, nil,
		),

		// password
		a.WithComponent("password",
			func(ctx context.Context) (err error) {
				pwd, err = password.New(storage.PasswordStorage())
				return
			}, nil,
		),
		// jwt
		a.WithComponent("jwt",
			func(context.Context) (err error) {
				authJWT, err = jwt.New(cfg.JWTSecret)
				return
			}, nil,
		),

		// nodes service
		a.WithComponent("nodes service",
			func(context.Context) (err error) {
				nodesService, err = nodes.New(poolSyncer, storage.NodesStorage())
				return
			}, nil,
		),
		// users service
		a.WithComponent("users service",
			func(context.Context) (err error) {
				usersService, err = users.New(poolSyncer, storage.UsersStorage())
				return
			}, nil,
		),
		// subscr service
		a.WithComponent("subscr service",
			func(context.Context) (err error) {
				subscrService, err = subscr.New(storage.SubscrStorage(), subscr.WithLogger(log))
				return
			}, nil,
		),
		// auth service
		a.WithComponent("auth service",
			func(ctx context.Context) (err error) {
				authService, err = auth.New(pwd, authJWT)
				return
			}, nil,
		),

		// bootstrap
		a.WithComponent("bootstrap",
			func(ctx context.Context) (err error) {
				bootstrapCfg := bootstrap.Config{
					AdminPassword: cfg.AdminPassword,
				}
				err = bootstrap.Bootstrap(ctx, bootstrapCfg, pwd, log)
				return
			}, nil,
		),

		// handler
		a.WithComponent("handler",
			func(context.Context) (err error) {
				h, err = handler.New(
					usersService,
					nodesService,
					subscrService,
					authService,
					handler.WithLogger(log))
				return
			}, nil,
		),

		// security handler
		a.WithComponent("security",
			func(context.Context) (err error) {
				s, err = security.New(authJWT)
				return
			}, nil,
		),

		// api handler
		a.WithComponent("api handler",
			func(context.Context) (err error) {
				apiHandler, err = api.NewHandler(h, s)
				return
			}, nil,
		),

		// userpage spa
		a.WithComponent("userpage spa",
			func(context.Context) (err error) {
				userpageSpa, err = fs.Sub(userpageFS, "userpage")
				if err != nil {
					return errdefs.WrapWithStack(err)
				}
				return nil
			}, nil,
		),

		// adminpage spa

		a.WithComponent("adminpage spa",
			func(context.Context) (err error) {
				admpageSpa, err = fs.Sub(admpageFS, "admpage")
				if err != nil {
					return errdefs.WrapWithStack(err)
				}
				return nil
			}, nil,
		),

		// router
		a.WithComponent("router",
			func(context.Context) (err error) {
				userpageCfg := map[string]string{
					"api_prefix":  cfg.APIPrefix,
					"user_prefix": cfg.UserSpaPrefix,
				}
				admpageCfg := map[string]string{
					"api_prefix":   cfg.APIPrefix,
					"admin_prefix": cfg.AdminSpaPrefix,
				}

				r, err = router.New(
					router.WithHandler(cfg.APIPrefix, apiHandler),
					router.WithSPA(cfg.UserSpaPrefix, userpageSpa, userpageCfg),
					router.WithSPA(cfg.AdminSpaPrefix, admpageSpa, admpageCfg),
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

		// background syncer
		a.WithComponent("background sync",
			func(context.Context) (err error) {
				syncJob, err = syncman.New(poolSyncer, syncman.WithLogger(log))
				return
			}, nil,
		),

		// ------------ RUNNERS ------------- //
		a.WithRunner("http server",
			func() (err error) {
				return httpServer.Listen()
			},
			func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
		),
		// background syncer
		a.WithRunner("background sync job",
			func() (err error) {
				return syncJob.Run()
			},
			func(context.Context) (err error) {
				return syncJob.Stop()
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
