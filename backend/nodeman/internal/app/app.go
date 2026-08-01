package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"

	"go.uber.org/zap"
)

type App struct {
	core *gx.App
}

const JWTIssuer = "nodeman"

func New(rawCfg config.RawConfig, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, errdefs.NilArg("log")
	}

	///////////////////////////////////////////////////////////////////////////
	// create app components - chaotic good init order
	var ParamsProvider = gx.Options(
		gx.Provide(
			func() config.RawConfig {
				return rawCfg
			},
		),
		gx.WithLogger(log),
	)

	var InfraModule = gx.Provide(
		CfgProvider,
		DbProvider,
		StorageProvider,
		JwtProvider,
	)

	var PoolOpsModule = gx.Provide(
		HttpClientProvider,
		PoolClientProvider,
		SyncClientProvider,
		PoolSyncProvider,
		StatsClientProvider,
		PoolStatsProvider,
	)
	var ServicesModule = gx.Provide(
		SyncTimeoutProvider,
		NodesServiceProvider,
		UsersServiceProvider,
		SubscrServiceProvider,
		SettingsServiceProvider,
		AuthServiceProvider,
	)

	var HttpServerModule = gx.Provide(
		HttpHandlerProvider,
		HttpSecurityProvider,
		ApiHandlerProvider,

		UserPageProvider,
		AdminPageProvider,
		RouterProvider,
		ServerProvider,
	)
	/*
		// admpage spa
		adminCfgHandler := func(ctx context.Context) (*pagecfg.AdminPageCfg, error) {
			return &pagecfg.AdminPageCfg{
				ApiPrefix:    cfg.ApiServiceUrl,
				AdminPrefix:  cfg.AdminSpaUrl,
				UserPrefix:   cfg.UserSpaUrl,
				SettingsTags: s.subscr.SubHeadersPlaceholders(),
			}, nil
		}

		admpageSpa, err := pages.NewAdmPage(adminCfgHandler)
		if err != nil {
			return
		}

		// router
		r, err := router.New(
			router.WithHandler(cfg.ApiServicePath, apiHandler),
			router.WithSPA(cfg.UserSpaPath, userpageSpa),
			router.WithSPA(cfg.AdminSpaPath, admpageSpa),
			router.WithCrossOrigin(cfg.AllowedOrigins),
			router.WithLogger(log))
		if err != nil {
			return
		}

		// http server
		if h, err = server.New(cfg.Endpoint, r); err != nil {
			return
		}

		// log info
		log.Warn(fmt.Sprintf("api available on %s via %s",
			cfg.ApiServicePath, cfg.ApiServiceUrl))
		log.Warn(fmt.Sprintf("user page available on %s via %s",
			cfg.UserSpaPath, cfg.UserSpaUrl))
		log.Warn(fmt.Sprintf("admin page available on %s via %s",
			cfg.AdminSpaPath, cfg.AdminSpaUrl))

		return
	*/

	appcore := gx.New(
		ParamsProvider,
		InfraModule,
		PoolOpsModule,
		ServicesModule,
		HttpServerModule,
		SyncJobProvider,
		StatsJobProvider,
		HelloMessageInvoker,
		HttpServerJob,
		BackgroundSyncJob,
		BackgroundStatsJob,
	)
	return &App{
		core: &appcore,
	}, nil
}

func (app *App) Run() error {
	if app == nil || app.core == nil {
		return xerr.NilCall()
	}
	return app.core.Run()
}

/*ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gx.Start(ctx); err != nil {
		panic(err)
	}

	fmt.Println("started")

	if err := gx.Stop(ctx); err != nil {
		panic(err)
	}

	fmt.Println("stopped")
}

/*var ServiceModule = gx.Options()
var HTTPModule = gx.Options()
var JobsModule = gx.Options()

var Module = gx.Options(
	gx.Provide()
	InfraModule,
	ServiceModule,
	HTTPModule,
	JobsModule,
)

// runtime config
cfg, err := config.Init(rawCfg)
if err != nil {
	return
}

// infrasturcture
infra, err := gx.initInfra(*cfg)
if err != nil {
	return
}
infra.storage.ExplainLog = log

// pool sync, pool stats
poolSyncer, poolStats, err := gx.initPoolOps(*cfg, *infra, log)
if err != nil {
	return
}

// services
services, err := gx.initServices(cfg, poolSyncer, infra.authJWT, infra.storage, log)
if err != nil {
	return
}

// http server
httpServer, err := gx.initHttpServer(*cfg, infra.storage, *services, infra.authJWT, log)
if err != nil {
	return
}

// background sync job
syncJob, err := syncman.New(poolSyncer, cfg.StateSyncInterval, syncman.WithLogger(log))
if err != nil {
	return
}

// background stats job
statsJob, err := statsman.New(poolStats, cfg.StateSyncInterval, statsman.WithLogger(log))
if err != nil {
	return
}

///////////////////////////////////////////////////////////////////////////
// bootstrap app components

// migrate db
gx.core.AddBootstrap("migrate db", func(ctx context.Context) error {
	return infra.storage.Migrate(ctx, dbstorage.WithLogger(log))
}, func(err error) bool {
	return errors.Is(err, errdefs.ErrTemporaryUnavailable)
})

// set password
gx.core.AddBootstrap("set password", func(ctx context.Context) error {
	if cfg.AdminPassword == "" {
		return nil
	}
	return services.auth.Update(ctx, cfg.AdminPassword)
}, func(err error) bool {
	return errors.Is(err, errdefs.ErrTemporaryUnavailable)
})

// set default settings
gx.core.AddBootstrap("ensure settings", func(ctx context.Context) error {
	return services.settings.EnsureSettings(ctx)
}, func(err error) bool {
	return errors.Is(err, errdefs.ErrTemporaryUnavailable)
})

///////////////////////////////////////////////////////////////////////////
// run app components

// http server
gx.core.AddRunner("http server",
	func() (err error) {
		return httpServer.Listen()
	},
	func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	},
)

// background syncer
gx.core.AddRunner("background sync",
	func() (err error) {
		return syncJob.Run()
	},
	func(context.Context) error {
		return syncJob.Stop()
	},
)

// background stats
gx.core.AddRunner("background stats",
	func() (err error) {
		return statsJob.Run()
	},
	func(context.Context) error {
		return statsJob.Stop()
	},
)

///////////////////////////////////////////////////////////////////////////

return*/

/*type infrasturcture struct {
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
	a.core.AddCloser(func(context.Context) (err error) {
		err = db.Close()
		return
	})

	// storage
	if infra.storage, err = dbstorage.New(db); err != nil {
		return
	}

	// JWT
	if infra.authJWT, err = jwt.New(cfg.JwtSecret, jwt.WithIssuer(JWTIssuer)); err != nil {
		return
	}
	return infra, nil
}

func (a *App) initPoolOps(cfg config.Config, infra infrasturcture, log *zap.Logger) (
	poolSyncer *poolsync.Syncer, poolStats *poolstats.Stats, err error,
) {
	// nodes http client
	nc, err := httpclient.NewClientFactory(cfg.NodeCallTimeout, log)
	if err != nil {
		return
	}
	a.core.AddCloser(func(context.Context) error {
		nc.Close()
		return nil
	})

	// pool client
	pc, err := client.NewPoolClient(nc, log)
	if err != nil {
		return
	}

	// pool syncer
	poolSyncer, err = poolsync.New(pc.PoolSyncClient(), infra.storage, log)
	if err != nil {
		return
	}

	// pool stats
	poolStats, err = poolstats.New(pc.PoolStatsClient(), infra.storage, log)
	if err != nil {
		return
	}

	return
}

type services struct {
	nodes    *nodes.Service
	users    *users.Service
	subscr   *subscr.Service
	settings *settings.Service
	auth     *auth.Service
}

func (a *App) initServices(
	cfg *config.Config,
	ps *poolsync.Syncer,
	authJWT *jwt.JWT,
	s *dbstorage.Storage,
	log *zap.Logger,
) (ss *services, err error) {
	ss = &services{}

	stateSyncTimeout := cfg.NodeCallTimeout + cfg.StorageCallTimeout
	// nodes service
	if ss.nodes, err = nodes.New(ps, s, stateSyncTimeout, log); err != nil {
		return
	}

	// users service
	if ss.users, err = users.New(ps, s, stateSyncTimeout, log); err != nil {
		return
	}

	// subscr service
	if ss.subscr, err = subscr.New(s, subscr.WithLogger(log)); err != nil {
		return
	}

	// settings service
	if ss.settings, err = settings.New(s); err != nil {
		return
	}

	// auth service
	if ss.auth, err = auth.New(s, authJWT); err != nil {
		return
	}

	return
}

func (a *App) initHttpServer(
	cfg config.Config,
	storage settings.Storage,
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
	userCfgHandler := func(ctx context.Context) (*pagecfg.UserPageCfg, error) {
		settings, err := storage.GetSettings(ctx)
		if err != nil {
			return nil, err
		}
		return &pagecfg.UserPageCfg{
			ApiPrefix:   cfg.ApiServiceUrl,
			UserPrefix:  cfg.UserSpaUrl,
			SupportLink: settings.TgPage,
			AppLinks:    settings.AppLinks,
		}, nil
	}

	userpageSpa, err := pages.NewUserPage(userCfgHandler)
	if err != nil {
		return
	}

	// admpage spa
	adminCfgHandler := func(ctx context.Context) (*pagecfg.AdminPageCfg, error) {
		return &pagecfg.AdminPageCfg{
			ApiPrefix:    cfg.ApiServiceUrl,
			AdminPrefix:  cfg.AdminSpaUrl,
			UserPrefix:   cfg.UserSpaUrl,
			SettingsTags: s.subscr.SubHeadersPlaceholders(),
		}, nil
	}

	admpageSpa, err := pages.NewAdmPage(adminCfgHandler)
	if err != nil {
		return
	}

	// router
	r, err := router.New(
		router.WithHandler(cfg.ApiServicePath, apiHandler),
		router.WithSPA(cfg.UserSpaPath, userpageSpa),
		router.WithSPA(cfg.AdminSpaPath, admpageSpa),
		router.WithCrossOrigin(cfg.AllowedOrigins),
		router.WithLogger(log))
	if err != nil {
		return
	}

	// http server
	if h, err = server.New(cfg.Endpoint, r); err != nil {
		return
	}

	// log info
	log.Warn(fmt.Sprintf("api available on %s via %s",
		cfg.ApiServicePath, cfg.ApiServiceUrl))
	log.Warn(fmt.Sprintf("user page available on %s via %s",
		cfg.UserSpaPath, cfg.UserSpaUrl))
	log.Warn(fmt.Sprintf("admin page available on %s via %s",
		cfg.AdminSpaPath, cfg.AdminSpaUrl))

	return
}

func (a *App) initHandler(s services, authJWT *jwt.JWT, log *zap.Logger) (h http.Handler, err error) {
	// requests handler
	reqH, err := handler.New(
		s.users,
		s.nodes,
		s.subscr,
		s.settings,
		s.auth,
		log)
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
		return errdefs.NilCall()
	}

	if err := gx.core.Bootstrap(); err != nil {
		return err
	}

	return gx.core.Run()
}*/
