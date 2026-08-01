package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/common/app"
	fx "github.com/XRay-Addons/xrayman/common/app"
	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/http/server"
	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	genapi "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"

	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqldb"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/api"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/jwt"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/zap"
)

var CfgProvider = config.Init

var DbProvider = app.Annotate(
	func(lc app.Lifecycle, cfg *config.Config) (*sqldb.SqlDB, error) {
		db, err := sqldb.New(cfg.DBConn)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser(app.Closer{
			Name: "db",
			OnStop: func(context.Context) error {
				return db.Close()
			},
		})
		return db, nil
	},
	app.As(new(dbstorage.DB)),
)

var StorageProvider = app.Annotate(
	dbstorage.New,
	app.As(new(users.Storage)),
	app.As(new(nodes.Storage)),
	app.As(new(subscr.Storage)),
	app.As(new(poolsync.Storage)),
	app.As(new(poolstats.Storage)),
	app.As(new(settings.Storage)),
	app.As(new(auth.Storage)),
)

var JwtProvider = app.Annotate(
	func(cfg *config.Config) (*jwt.JWT, error) {
		return jwt.New(cfg.JwtSecret, jwt.WithIssuer(JWTIssuer))
	},
	app.As(new(auth.JWT)),
	app.As(new(security.JWT)),
)

var HttpClientProvider = app.Annotate(
	func(lc app.Lifecycle, cfg *config.Config, log *zap.Logger) (*httpclient.ClientFactory, error) {
		clientFactory, err := httpclient.NewClientFactory(cfg.NodeCallTimeout, log)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser(app.Closer{
			Name: "client factory",
			OnStop: func(context.Context) error {
				clientFactory.Close()
				return nil
			},
		})
		return clientFactory, nil
	},
	app.As(new(client.HTTPClientFactory)),
	app.As(app.Self()),
)

var PoolClientProvider = client.NewPoolClient

func SyncClientProvider(poolclient *client.PoolClient) poolsync.Client {
	return poolclient.PoolSyncClient()
}

var PoolSyncProvider = app.Annotate(
	poolsync.New,
	app.As(new(users.Syncer)),
	app.As(new(nodes.Syncer)),
	app.As(new(syncman.PoolSyncer)),
)

func StatsClientProvider(poolclient *client.PoolClient) poolstats.Client {
	return poolclient.PoolStatsClient()
}

var PoolStatsProvider = app.Annotate(
	poolstats.New,
	app.As(new(statsman.StatsUpdater)),
)

var SyncTimeoutProvider = app.Annotate(
	func(cfg *config.Config) time.Duration {
		return cfg.StorageCallTimeout + cfg.NodeCallTimeout
	},
	app.ResultTags(`name:"sync-timeout"`),
)

var NodesServiceProvider = app.Annotate(
	nodes.New,
	app.ParamTags(``, ``, `name:"sync-timeout"`, ``),
	app.As(new(handler.NodesService)),
)
var UsersServiceProvider = app.Annotate(
	users.New,
	app.ParamTags(``, ``, `name:"sync-timeout"`, ``),
	app.As(new(handler.UsersService)),
)
var SubscrServiceProvider = app.Annotate(
	subscr.New,
	app.As(new(handler.SubscrService)),
	app.As(app.Self()),
)
var SettingsServiceProvider = app.Annotate(
	settings.New,
	app.As(new(handler.SettingsService)),
)
var AuthServiceProvider = app.Annotate(
	auth.New,
	app.As(new(handler.AuthService)),
)

var HttpHandlerProvider = app.Annotate(
	handler.New,
	app.As(new(genapi.Handler)),
)
var HttpSecurityProvider = app.Annotate(
	security.New,
	app.As(new(genapi.SecurityHandler)),
)
var ApiHandlerProvider = app.Annotate(
	api.NewHandler,
	app.ResultTags(`name:"api-handler"`),
)

var UserPageProvider = app.Annotate(
	func(cfg *config.Config, s settings.Storage) (*pages.Page, error) {
		pageConfigHandler := func(ctx context.Context) (*pagecfg.UserPageCfg, error) {
			settings, err := s.GetSettings(ctx)
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
		return pages.NewUserPage(pageConfigHandler)
	},
	app.ResultTags(`name:"user-page"`),
)

var AdminPageProvider = app.Annotate(
	func(cfg *config.Config, s *subscr.Service) (*pages.Page, error) {
		pageConfigHandler := func(ctx context.Context) (*pagecfg.AdminPageCfg, error) {
			return &pagecfg.AdminPageCfg{
				ApiPrefix:    cfg.ApiServiceUrl,
				AdminPrefix:  cfg.AdminSpaUrl,
				UserPrefix:   cfg.UserSpaUrl,
				SettingsTags: s.SubHeadersPlaceholders(),
			}, nil
		}
		return pages.NewAdmPage(pageConfigHandler)
	},
	app.ResultTags(`name:"admin-page"`),
)

type RouterProviderParams struct {
	app.In

	ApiHandler http.Handler `name:"api-handler"`
	UserPage   *pages.Page  `name:"user-page"`
	AdminPage  *pages.Page  `name:"admin-page"`
}

var RouterProvider = app.Annotate(
	func(cfg *config.Config, p RouterProviderParams, l *zap.Logger) (http.Handler, error) {
		return router.New(
			router.WithHandler(cfg.ApiServicePath, p.ApiHandler),
			router.WithSPA(cfg.UserSpaPath, p.UserPage),
			router.WithSPA(cfg.AdminSpaPath, p.AdminPage),
			router.WithCrossOrigin(cfg.AllowedOrigins),
			router.WithLogger(l))
	},
	app.ResultTags(`name:"router"`),
)

type ServerProviderParams struct {
	app.In
	Router http.Handler `name:"router"`
}

var ServerProvider = app.Annotate(
	func(cfg *config.Config, p ServerProviderParams) (*server.HttpServer, error) {
		return server.New(cfg.Endpoint, p.Router)
	},
)

var SyncJobProvider = app.Provide(
	func(ps syncman.PoolSyncer, cfg *config.Config, l *zap.Logger) (*syncman.SyncMan, error) {
		return syncman.New(ps, cfg.StateSyncInterval, syncman.WithLogger(l))
	},
)

var StatsJobProvider = app.Provide(
	func(ps statsman.StatsUpdater, cfg *config.Config, l *zap.Logger) (*statsman.StatsMan, error) {
		return statsman.New(ps, cfg.StatsSyncInterval, statsman.WithLogger(l))
	},
)

var HelloMessageInvoker = fx.Invoke(
	func(cfg *config.Config, log *zap.Logger) {
		log.Warn(fmt.Sprintf("api available on %s via %s",
			cfg.ApiServicePath, cfg.ApiServiceUrl))
		log.Warn(fmt.Sprintf("user page available on %s via %s",
			cfg.UserSpaPath, cfg.UserSpaUrl))
		log.Warn(fmt.Sprintf("admin page available on %s via %s",
			cfg.AdminSpaPath, cfg.AdminSpaUrl))
	},
)

var HttpServerJob = fx.Invoke(
	func(s *server.HttpServer, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
			Name: "http server",
			OnStart: func(context.Context) error {
				return s.Listen()
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		})
	},
)

var BackgroundSyncJob = fx.Invoke(
	func(s *syncman.SyncMan, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
			Name: "background sync",
			OnStart: func(context.Context) error {
				return s.Run()
			},
			OnStop: func(context.Context) error {
				s.Stop()
				return nil
			},
		})
	},
)

var BackgroundStatsJob = fx.Invoke(
	func(s *statsman.StatsMan, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
			Name: "background stats",
			OnStart: func(context.Context) error {
				return s.Run()
			},
			OnStop: func(context.Context) error {
				s.Stop()
				return nil
			},
		})
	},
)
