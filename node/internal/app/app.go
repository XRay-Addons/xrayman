package app

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/api"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/infra/auth/jwt"
	appcore "github.com/XRay-Addons/xrayman/node/internal/infra/common/app"
	"github.com/XRay-Addons/xrayman/node/internal/infra/common/http/router"
	"github.com/XRay-Addons/xrayman/node/internal/infra/common/http/server"
	"github.com/XRay-Addons/xrayman/node/internal/infra/secrets"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tlscfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xraycfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayservice"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"go.uber.org/zap"
)

type App struct {
	base *appcore.App
}

func New(cfg config.Config, log *zap.Logger) (app *App, err error) {
	if log == nil {
		return nil, errdefs.NilArg("log")
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

	// secrets config
	sec, err := secrets.Init(cfg.PersistentDir)
	if err != nil {
		return
	}
	log.Info("node access", zap.String("key", sec.AccessKey.String()))

	// server config
	srvCfg, err := xraycfg.NewServerCfg(cfg.XRayServer())
	if err != nil {
		return
	}

	// client config
	clientCfg, err := xraycfg.NewClientConfig(cfg.XRayClient())
	if err != nil {
		return
	}

	// TLS config
	tlsCfg, err := tlscfg.Load(sec.Cert, sec.Key)
	if err != nil {
		return
	}

	// xray service
	xrayService, err := xrayservice.New(log)
	if err != nil {
		return
	}
	app.base.AddCloser(func(ctx context.Context) error {
		return xrayService.Close(ctx)
	})

	// xray api
	xrayAPI, err := xrayapi.New(srvCfg.GetApiURL(), srvCfg.GetInbounds(),
		xrayapi.WithLogger(log))
	if err != nil {
		return
	}
	app.base.AddCloser(func(ctx context.Context) error {
		return xrayAPI.Close(ctx)
	})

	// service
	s, err := service.New(srvCfg, clientCfg, xrayService, xrayAPI)
	if err != nil {
		return
	}

	// handler
	h, err := handler.New(s, log)
	if err != nil {
		return
	}

	// jwt
	jwt, err := jwt.New(sec.AccessKey.AccessSecret)
	if err != nil {
		return
	}

	// security
	apiSec, err := security.New(jwt)
	if err != nil {
		return
	}

	// api handler
	apiHandler, err := api.NewHandler(h, apiSec)
	if err != nil {
		return
	}

	// router
	r, err := router.New(
		router.WithHandler("/", apiHandler),
		router.WithLogger(log))
	if err != nil {
		return
	}
	if err != nil {
		return
	}

	// http server
	httpServer, err := server.New(cfg.Endpoint, r, server.WithTLS(tlsCfg))
	if err != nil {
		return
	}

	app.base.AddRunner("http server",
		func() (err error) {
			return httpServer.Listen()
		},
		func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	)

	return
}

func (app *App) Run() error {
	if app == nil {
		return errdefs.NilCall()
	}

	if err := app.base.Bootstrap(); err != nil {
		return err
	}

	return app.base.Run()
}
