package app

import (
	"context"
	"net/http"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/api"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	genapi "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
	"go.uber.org/zap"
)

var apiHandler = gx.Options(
	gx.ProvideAnnotated(
		handler.New,
		gx.As(new(genapi.Handler)),
	),
	gx.ProvideAnnotated(
		security.New,
		gx.As(new(genapi.SecurityHandler)),
	),
	gx.ProvideNamed(
		api.NewHandler,
		"api-handler",
	),
)

var userPage = gx.ProvideNamed(
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
	"user-page",
)

var adminPage = gx.ProvideNamed(
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
	"admin-page",
)

type HttpRouterParams struct {
	gx.In
	Cfg        *config.Config
	Log        *zap.Logger
	ApiHandler http.Handler `name:"api-handler"`
	UserPage   *pages.Page  `name:"user-page"`
	AdminPage  *pages.Page  `name:"admin-page"`
}

var httpRouter = gx.ProvideNamed(
	func(p HttpRouterParams) (http.Handler, error) {
		return router.New(
			router.WithHandler(p.Cfg.ApiServicePath, p.ApiHandler),
			router.WithSPA(p.Cfg.UserSpaPath, p.UserPage),
			router.WithSPA(p.Cfg.AdminSpaPath, p.AdminPage),
			router.WithCrossOrigin(p.Cfg.AllowedOrigins),
			router.WithLogger(p.Log))
	},
	"http-router",
)

type HttpServerParams struct {
	gx.In
	Cfg    *config.Config
	Router http.Handler `name:"http-router"`
}

var httpServer = gx.ProvideNamed(
	func(p HttpServerParams) (*server.HttpServer, error) {
		return server.New(p.Cfg.Endpoint, p.Router)
	},
	"http-server",
)

type HttpServerJobParams struct {
	gx.In
	S *server.HttpServer `name:"http-server"`
}

var httpServerJob = gx.Invoke(
	func(p HttpServerJobParams, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
			Name: "http server",
			OnStart: func(context.Context) error {
				return p.S.Listen()
			},
			OnStop: func(ctx context.Context) error {
				return p.S.Shutdown(ctx)
			},
		})
	},
)

var HttpServer = gx.Module("server",
	apiHandler,
	userPage,
	adminPage,
	httpRouter,
	httpServer,
	httpServerJob,
)
