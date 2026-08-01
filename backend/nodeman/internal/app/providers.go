package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
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

var DbProvider = gx.Annotate(
	func(lc gx.Lifecycle, cfg *config.Config) (*sqldb.SqlDB, error) {
		db, err := sqldb.New(cfg.DBConn)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser(gx.Closer{
			Name: "db",
			OnStop: func(context.Context) error {
				return db.Close()
			},
		})
		return db, nil
	},
	gx.As(new(dbstorage.DB)),
)

var StorageProvider = gx.Annotate(
	dbstorage.New,
	gx.As(new(users.Storage)),
	gx.As(new(nodes.Storage)),
	gx.As(new(subscr.Storage)),
	gx.As(new(poolsync.Storage)),
	gx.As(new(poolstats.Storage)),
	gx.As(new(settings.Storage)),
	gx.As(new(auth.Storage)),
)

var JwtProvider = gx.Annotate(
	func(cfg *config.Config) (*jwt.JWT, error) {
		return jwt.New(cfg.JwtSecret, jwt.WithIssuer(JWTIssuer))
	},
	gx.As(new(auth.JWT)),
	gx.As(new(security.JWT)),
)

var HttpClientProvider = gx.Annotate(
	func(lc gx.Lifecycle, cfg *config.Config, log *zap.Logger) (*httpclient.ClientFactory, error) {
		clientFactory, err := httpclient.NewClientFactory(cfg.NodeCallTimeout, log)
		if err != nil {
			return nil, err
		}
		lc.AppendCloser(gx.Closer{
			Name: "client factory",
			OnStop: func(context.Context) error {
				clientFactory.Close()
				return nil
			},
		})
		return clientFactory, nil
	},
	gx.As(new(client.HTTPClientFactory)),
	gx.As(gx.Self()),
)

var PoolClientProvider = client.NewPoolClient

func SyncClientProvider(poolclient *client.PoolClient) poolsync.Client {
	return poolclient.PoolSyncClient()
}

var PoolSyncProvider = gx.Annotate(
	poolsync.New,
	gx.As(new(users.Syncer)),
	gx.As(new(nodes.Syncer)),
	gx.As(new(syncman.PoolSyncer)),
)

func StatsClientProvider(poolclient *client.PoolClient) poolstats.Client {
	return poolclient.PoolStatsClient()
}

var PoolStatsProvider = gx.Annotate(
	poolstats.New,
	gx.As(new(statsman.StatsUpdater)),
)

var SyncTimeoutProvider = gx.Annotate(
	func(cfg *config.Config) time.Duration {
		return cfg.StorageCallTimeout + cfg.NodeCallTimeout
	},
	gx.ResultTags(`name:"sync-timeout"`),
)

var NodesServiceProvider = gx.Annotate(
	nodes.New,
	gx.ParamTags(``, ``, `name:"sync-timeout"`, ``),
	gx.As(new(handler.NodesService)),
)
var UsersServiceProvider = gx.Annotate(
	users.New,
	gx.ParamTags(``, ``, `name:"sync-timeout"`, ``),
	gx.As(new(handler.UsersService)),
)
var SubscrServiceProvider = gx.Annotate(
	subscr.New,
	gx.As(new(handler.SubscrService)),
	gx.As(gx.Self()),
)
var SettingsServiceProvider = gx.Annotate(
	settings.New,
	gx.As(new(handler.SettingsService)),
)
var AuthServiceProvider = gx.Annotate(
	auth.New,
	gx.As(new(handler.AuthService)),
)

var HttpHandlerProvider = gx.Annotate(
	handler.New,
	gx.As(new(genapi.Handler)),
)
var HttpSecurityProvider = gx.Annotate(
	security.New,
	gx.As(new(genapi.SecurityHandler)),
)
var ApiHandlerProvider = gx.Annotate(
	api.NewHandler,
	gx.ResultTags(`name:"api-handler"`),
)

var UserPageProvider = gx.Annotate(
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
	gx.ResultTags(`name:"user-page"`),
)

var AdminPageProvider = gx.Annotate(
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
	gx.ResultTags(`name:"admin-page"`),
)

type RouterProviderParams struct {
	gx.In

	ApiHandler http.Handler `name:"api-handler"`
	UserPage   *pages.Page  `name:"user-page"`
	AdminPage  *pages.Page  `name:"admin-page"`
}

var RouterProvider = gx.Annotate(
	func(cfg *config.Config, p RouterProviderParams, l *zap.Logger) (http.Handler, error) {
		return router.New(
			router.WithHandler(cfg.ApiServicePath, p.ApiHandler),
			router.WithSPA(cfg.UserSpaPath, p.UserPage),
			router.WithSPA(cfg.AdminSpaPath, p.AdminPage),
			router.WithCrossOrigin(cfg.AllowedOrigins),
			router.WithLogger(l))
	},
	//gx.ResultTags(`name:"router"`),
)

type ServerProviderParams struct {
	gx.In
	Router http.Handler //`name:"router"`
}

var ServerProvider = gx.Annotate(
	func(cfg *config.Config, p ServerProviderParams) (*server.HttpServer, error) {
		return server.New(cfg.Endpoint, p.Router)
	},
)

var SyncJobProvider = gx.Provide(
	func(ps syncman.PoolSyncer, cfg *config.Config, l *zap.Logger) (*syncman.SyncMan, error) {
		return syncman.New(ps, cfg.StateSyncInterval, syncman.WithLogger(l))
	},
)

var StatsJobProvider = gx.Provide(
	func(ps statsman.StatsUpdater, cfg *config.Config, l *zap.Logger) (*statsman.StatsMan, error) {
		return statsman.New(ps, cfg.StatsSyncInterval, statsman.WithLogger(l))
	},
)

var HelloMessageInvoker = gx.Invoke(
	func(cfg *config.Config, log *zap.Logger) {
		log.Warn(fmt.Sprintf("api available on %s via %s",
			cfg.ApiServicePath, cfg.ApiServiceUrl))
		log.Warn(fmt.Sprintf("user page available on %s via %s",
			cfg.UserSpaPath, cfg.UserSpaUrl))
		log.Warn(fmt.Sprintf("admin page available on %s via %s",
			cfg.AdminSpaPath, cfg.AdminSpaUrl))
	},
)

var HttpServerJob = gx.Invoke(
	func(s *server.HttpServer, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
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

var BackgroundSyncJob = gx.Invoke(
	func(s *syncman.SyncMan, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
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

var BackgroundStatsJob = gx.Invoke(
	func(s *statsman.StatsMan, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
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
