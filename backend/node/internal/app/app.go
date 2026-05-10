package app

import (
	"context"
	"crypto/tls"
	"net/http"

	appcore "github.com/XRay-Addons/xrayman/common/app"
	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/api"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/infra/auth/jwt"
	"github.com/XRay-Addons/xrayman/node/internal/infra/secrets"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tlscfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/clientcfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/servercfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayservice"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"go.uber.org/zap"
)

type App struct {
	core *appcore.App
}

func New(cfg config.Config, log *zap.Logger) (app *App, err error) {
	if log == nil {
		return nil, errdefs.NilArg("log")
	}

	app = &App{
		core: appcore.New(appcore.WithLogger(log)),
	}

	defer func() {
		if err != nil {
			err = xerr.Join(err, app.core.Close())
		}
	}()

	// configs
	configs, err := app.initConfigs(cfg)
	if err != nil {
		return
	}
	log.Warn("node access", zap.String("key",
		configs.accessKey.String()))
	log.Warn("node access", zap.String("key",
		configs.accessKey.String()))
	log.Warn("node access", zap.String("key",
		configs.accessKey.String()))

	// xray service
	xrayService, err := xrayservice.New(cfg.XRayDataDir,
		xrayservice.WithLogger(log))
	if err != nil {
		return
	}
	app.core.AddCloser(func(ctx context.Context) error {
		return xrayService.Close(ctx)
	})

	// xray api
	xrayAPI, err := xrayapi.New(configs.serverCfg.GetApiURL(),
		configs.serverCfg.GetInbounds(), xrayapi.WithLogger(log))
	if err != nil {
		return
	}
	app.core.AddCloser(func(ctx context.Context) error {
		return xrayAPI.Close(ctx)
	})

	// service
	s, err := service.New(configs.serverCfg, configs.clientCfg,
		xrayService, xrayAPI)
	if err != nil {
		return
	}

	// jwt
	jwt, err := jwt.New(configs.accessKey.AccessSecret)
	if err != nil {
		return
	}

	// http
	httpServer, err := app.initHttpServer(cfg, s, jwt, configs.tlsCfg, log)
	if err != nil {
		return
	}

	///////////////////////////////////////////////////////////////////////////
	// bootstrap - check xray config ok
	app.core.AddBootstrap("xray check", func(ctx context.Context) (err error) {
		log.Info("bootstrap...")
		// start service
		if _, err = s.Start(ctx, models.StartParams{}); err != nil {
			return
		}
		defer func() {
			_, closeErr := s.Stop(ctx, models.StopParams{})
			err = xerr.Join(err, closeErr)
		}()

		// ping xray
		if err = xrayAPI.Ping(ctx); err != nil {
			return err
		}
		log.Info("service ping ok!")
		return nil
	}, func(error) bool {
		return false
	})

	///////////////////////////////////////////////////////////////////////////
	// run app components

	// http server
	app.core.AddRunner("http server",
		func() (err error) {
			return httpServer.Listen()
		},
		func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	)

	return
}

type configs struct {
	accessKey models.AccessKey
	serverCfg *servercfg.Config
	clientCfg *clientcfg.Config
	tlsCfg    *tls.Config
}

func (a *App) initConfigs(
	cfg config.Config,
) (cfgs *configs, err error) {
	cfgs = &configs{}

	// secrets config
	secrets, err := secrets.Init(cfg.PersistentDir)
	if err != nil {
		return
	}
	cfgs.accessKey = secrets.AccessKey

	// server config
	cfgs.serverCfg, err = servercfg.New(cfg.XRayServer())
	if err != nil {
		return
	}

	// client config
	cfgs.clientCfg, err = clientcfg.New(cfg.XRayClient())
	if err != nil {
		return
	}

	// TLS config
	cfgs.tlsCfg, err = tlscfg.Load(secrets.Cert, secrets.Key)
	if err != nil {
		return
	}

	return
}

func (a *App) initHttpServer(
	cfg config.Config,
	s *service.Service,
	authJWT *jwt.JWT,
	tlsCfg *tls.Config,
	log *zap.Logger,
) (h *server.HttpServer, err error) {
	// api handler
	apiHandler, err := a.initHandler(s, authJWT, log)
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

	// http server
	if h, err = server.New(cfg.Endpoint, r, server.WithTLS(tlsCfg)); err != nil {
		return
	}

	return
}

func (a *App) initHandler(s *service.Service,
	authJWT *jwt.JWT, log *zap.Logger,
) (h http.Handler, err error) {
	// requests handler
	reqH, err := handler.New(s, handler.WithLogger(log))
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

	if err := app.core.Bootstrap(); err != nil {
		return err
	}

	return app.core.Run()
}
