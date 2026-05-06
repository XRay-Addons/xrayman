package app

import (
	"context"
	"errors"

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
	appcore "github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/sqldb"

	"go.uber.org/zap"
)

type App struct {
	base *appcore.App
}

func New(cfg config.Config, log *zap.Logger) (app *App, err error) {
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}

	baseApp := appcore.New(appcore.WithLogger(log))
	defer func() {
		if err != nil {
			err = errors.Join(err, baseApp.Close())
		}
	}()

	///////////////////////////////////////////////////////////////////////////
	// create app components

	// db
	db, err := sqldb.New(cfg.DBConn)
	if err != nil {
		return
	}
	baseApp.AddCloser(func(context.Context) error {
		db.Close()
		return nil
	})

	// storage
	storage, err := dbstorage.New(context.TODO(), db)
	if err != nil {
		return
	}

	// nodes http client
	nodesClient := httpclient.NewClientFactory()
	baseApp.AddCloser(func(context.Context) error {
		nodesClient.Close()
		return nil
	})

	// pool client
	poolClient, err := client.NewPoolClient(
		client.WithHTTPClient(nodesClient))
	if err != nil {
		return
	}

	// pool syncer
	poolSyncer, err := poolsync.New(
		poolClient, storage.PoolSyncStorage())
	if err != nil {
		return
	}

	// password
	pwd, err := password.New(storage.PasswordStorage())
	if err != nil {
		return
	}

	// JWT
	authJWT, err := jwt.New(cfg.JWTSecret)
	if err != nil {
		return
	}

	// nodes service
	nodesService, err := nodes.New(
		poolSyncer, storage.NodesStorage())
	if err != nil {
		return
	}

	// users service
	usersService, err := users.New(
		poolSyncer, storage.UsersStorage())
	if err != nil {
		return
	}

	// subscr service
	subscrService, err := subscr.New(
		storage.SubscrStorage(), subscr.WithLogger(log))
	if err != nil {
		return
	}

	// auth service
	authService, err := auth.New(pwd, authJWT)
	if err != nil {
		return
	}

	// requests handler
	requestsHandler, err := handler.New(
		usersService,
		nodesService,
		subscrService,
		authService,
		handler.WithLogger(log))
	if err != nil {
		return
	}

	// security handler
	securityHandler, err := security.New(authJWT)
	if err != nil {
		return
	}

	// api handler
	apiHandler, err := api.NewHandler(
		requestsHandler, securityHandler)
	if err != nil {
		return
	}

	// userpage spa
	userpageSpa, err := pages.NewUserPage(
		cfg.APIPrefix, cfg.UserSpaPrefix)
	if err != nil {
		return
	}
	// admpage spa
	admpageSpa, err := pages.NewAdmPage(
		cfg.APIPrefix, cfg.AdminSpaPrefix, cfg.UserSpaPrefix)
	if err != nil {
		return
	}

	// router
	r, err := router.New(
		router.WithHandler(cfg.APIPrefix, apiHandler),
		router.WithSPA(cfg.UserSpaPrefix, userpageSpa),
		router.WithSPA(cfg.AdminSpaPrefix, admpageSpa),
		router.WithLogger(log))
	if err != nil {
		return
	}

	// http server
	httpServer, err := server.New(cfg.Endpoint, r)
	if err != nil {
		return
	}

	// background sync job
	syncJob, err := syncman.New(poolSyncer, syncman.WithLogger(log))
	if err != nil {
		return
	}

	///////////////////////////////////////////////////////////////////////////
	// bootstrap app components

	// migrate db
	baseApp.AddBootstrap("migrated db", func(ctx context.Context) error {
		return storage.Migrage(ctx, dbstorage.WithLogger(log))
	}, func(err error) bool {
		return errors.Is(err, errdefs.ErrTemporaryUnavailable)
	})

	// set password
	baseApp.AddBootstrap("set password", func(ctx context.Context) error {
		if cfg.AdminPassword == "" {
			return nil
		}
		return pwd.Update(ctx, cfg.AdminPassword)
	}, func(err error) bool {
		return errors.Is(err, errdefs.ErrTemporaryUnavailable)
	})

	///////////////////////////////////////////////////////////////////////////
	// run app components

	// http server
	baseApp.AddRunner("http server",
		func() (err error) {
			return httpServer.Listen()
		},
		func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	)

	// background syncer
	baseApp.AddRunner("background sync",
		func() (err error) {
			return syncJob.Run()
		},
		func(context.Context) error {
			return syncJob.Stop()
		},
	)

	///////////////////////////////////////////////////////////////////////////

	return &App{
		base: baseApp,
	}, nil
}

func (app *App) Run() error {
	if app == nil {
		return errdefs.NewNilCall()
	}

	if err := app.base.Bootstrap(); err != nil {
		return err
	}

	return app.base.Run()
}
