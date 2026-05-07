package app

import (
	"context"
	"errors"
	"net/http"

	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/api"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/jwt"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	appcore "github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/http/router"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/httpclient"
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

const JWTIssuer = "nodeman"

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

	app = &App{
		base: baseApp,
	}

	///////////////////////////////////////////////////////////////////////////
	// create app components - chaotic good init order

	// infrasturcture
	infra, err := app.initInfra(cfg)
	if err != nil {
		return
	}

	// pool sync
	poolSyncer, err := app.initPoolSyncer(*infra)
	if err != nil {
		return
	}

	// password
	pwd, err := password.New(infra.storage.PasswordStorage())
	if err != nil {
		return
	}

	// services
	services, err := app.initServices(poolSyncer, pwd, infra.authJWT, infra.storage, log)
	if err != nil {
		return
	}

	// http server
	httpServer, err := app.initHttpServer(cfg, *services, infra.authJWT, log)
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
		return infra.storage.Migrage(ctx, dbstorage.WithLogger(log))
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

	return
}

type infrasturcture struct {
	storage *dbstorage.Storage
	authJWT *jwt.JWT
}

func (a *App) initInfra(cfg config.Config) (infra *infrasturcture, err error) {
	infra = &infrasturcture{}

	// db
	db, err := sqldb.New(cfg.DBConn)
	if err != nil {
		return
	}
	a.base.AddCloser(func(context.Context) error {
		db.Close()
		return nil
	})

	// storage
	if infra.storage, err = dbstorage.New(context.TODO(), db); err != nil {
		return
	}

	// JWT

	if infra.authJWT, err = jwt.New(cfg.JWTSecret, jwt.WithIssuer(JWTIssuer)); err != nil {
		return
	}
	return infra, nil
}

func (a *App) initPoolSyncer(infra infrasturcture) (ps poolsync.Syncer, err error) {
	// nodes http client
	nc := httpclient.NewClientFactory()
	a.base.AddCloser(func(context.Context) error {
		nc.Close()
		return nil
	})

	// pool client
	pc, err := client.NewPoolClient(client.WithHTTPClient(nc))
	if err != nil {
		return
	}

	// pool syncer
	ps, err = poolsync.New(pc, infra.storage.PoolSyncStorage())
	if err != nil {
		return
	}

	return
}

type services struct {
	nodes  *nodes.Service
	users  *users.Service
	subscr *subscr.Service
	auth   *auth.Service
}

func (a *App) initServices(
	ps poolsync.Syncer,
	pwd *password.Password,
	authJWT *jwt.JWT,
	s *dbstorage.Storage,
	log *zap.Logger,
) (ss *services, err error) {
	ss = &services{}

	// nodes service
	if ss.nodes, err = nodes.New(ps, s.NodesStorage()); err != nil {
		return
	}

	// users service
	if ss.users, err = users.New(ps, s.UsersStorage()); err != nil {
		return
	}

	// subscr service
	if ss.subscr, err = subscr.New(s.SubscrStorage(), subscr.WithLogger(log)); err != nil {
		return
	}

	// auth service
	if ss.auth, err = auth.New(pwd, authJWT); err != nil {
		return
	}

	return
}

func (a *App) initHttpServer(
	cfg config.Config,
	s services,
	authJWT *jwt.JWT,
	log *zap.Logger,
) (h *server.HttpServer, err error) {
	// api handler
	apiHandler, err := a.initHandler(s, authJWT, log)
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
	if h, err = server.New(cfg.Endpoint, r); err != nil {
		return
	}

	return
}
func (a *App) initHandler(s services, authJWT *jwt.JWT, log *zap.Logger) (h http.Handler, err error) {
	// requests handler
	reqH, err := handler.New(
		s.users,
		s.nodes,
		s.subscr,
		s.auth,
		handler.WithLogger(log))
	if err != nil {
		return
	}

	// security handler
	secH, err := security.New(authJWT)
	if err != nil {
		return
	}

	// api handler
	if h, err = api.NewHandler(reqH, secH); err != nil {
		return
	}

	return
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
