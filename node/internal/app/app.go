package app

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/http/server"
	"github.com/XRay-Addons/xrayman/node/internal/http/tlscfg"
	a "github.com/XRay-Addons/xrayman/node/internal/infra/app"
	"github.com/XRay-Addons/xrayman/node/internal/secrets"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xraycfg"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xrayservice"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"

	"go.uber.org/zap"
)

type App struct {
	app *a.App
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}

	var sec *secrets.Secrets
	var srvCfg *xraycfg.ServerCfg
	var clientCfg *xraycfg.ClientCfg
	var tlsCfg *tls.Config
	var xrayService *xrayservice.XRayService
	var xrayAPI *xrayapi.XRayApi

	var s *service.Service
	var h *handler.Handler
	var apiSec api.SecurityHandler
	var r http.Handler

	var httpServer *server.HttpServer

	app := a.New(
		// secrets config
		a.WithComponent("secrets",
			func() (err error) {
				sec, err = secrets.Init(cfg.PersistentDir)
				if err != nil {
					return
				}
				log.Info("node access", zap.String("key", sec.AccessKey.String()))
				return
			}, nil,
		),
		// server config
		a.WithComponent("server cfg",
			func() (err error) {
				srvCfg, err = xraycfg.NewServerCfg(cfg.XRayServer())
				return
			}, nil,
		),
		// client config
		a.WithComponent("client cfg",
			func() (err error) {
				clientCfg, err = xraycfg.NewClientCfg(cfg.XRayClient())
				return
			}, nil,
		),
		// TLS config
		a.WithComponent("tls cfg",
			func() (err error) {
				tlsCfg, err = tlscfg.Load(sec.Cert, sec.Key)
				return
			}, nil,
		),
		// xray service
		a.WithComponent("xray service",
			func() (err error) {
				xrayService, err = xrayservice.New(cfg.XRayExec(), cfg.XRayServer(), log)
				return err
			},
			func(ctx context.Context) error {
				return xrayService.Close(ctx)
			},
		),
		// xray api
		a.WithComponent("xray api",
			func() (err error) {
				xrayAPI, err = xrayapi.New(srvCfg.GetApiURL(), srvCfg.GetInbounds(), log)
				return
			},
			func(ctx context.Context) error {
				return xrayAPI.Close(ctx)
			},
		),
		// service
		a.WithComponent("service",
			func() (err error) {
				s, err = service.New(srvCfg, clientCfg, xrayService, xrayAPI)
				return
			}, nil,
		),
		// handler
		a.WithComponent("handler",
			func() (err error) {
				h, err = handler.New(s, log)
				return
			}, nil,
		),
		// security
		a.WithComponent("security",
			func() (err error) {
				jwtsec := sec.AccessKey.AccessSecret
				apiSec = security.New(jwtsec)
				return
			}, nil,
		),
		// router
		a.WithComponent("router",
			func() (err error) {
				r, err = router.New(h, apiSec, router.WithLogger(log))
				return
			}, nil,
		),

		// http server
		a.WithComponent("http server",
			func() (err error) {
				httpServer, err = server.New(cfg.Endpoint, r, tlsCfg)
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

		// cancel by Ctrl+C
		a.WithSignalCancel(),
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
