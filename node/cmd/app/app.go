package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/http/handlers"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	"github.com/XRay-Addons/xrayman/node/internal/http/server"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"github.com/XRay-Addons/xrayman/node/internal/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/xraycfg"
	"github.com/XRay-Addons/xrayman/node/internal/xrayctl"

	"go.uber.org/zap"
)

// RAII-like stuff
type (
	Closer = func(ctx context.Context) error
	Initer = func(ctx context.Context) (Closer, error)
)

type App struct {
	cfg config.Config

	xrayCfg XRayCfg
	xrayCtl XRayCtl
	xrayAPI XRayApi

	service handlers.Service

	handler http.Handler
	server  *server.Server

	log *zap.Logger

	closers []Closer
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		log = zap.NewNop()
	}

	app := App{cfg: cfg, log: log}

	initers := []Initer{
		app.initXRayCfg,
		app.initXRayCtl,
		app.initXRayAPI,
		app.initService,
		app.initHandlers,
		app.initServer,
	}

	if err := app.init(initers); err != nil {
		return nil, err
	}

	return &app, nil
}

func (app *App) Close() error {
	// close with timeout
	closeTimeout := 15 * time.Second
	ctx, close := context.WithTimeout(context.TODO(), closeTimeout)
	defer close()

	var closeErrs []error
	for i := len(app.closers) - 1; i >= 0; i-- {
		if app.closers[i] == nil {
			continue
		}
		if err := app.closers[i](ctx); err != nil {
			closeErrs = append(closeErrs, err)
		}
	}

	if len(closeErrs) == 0 {
		return nil
	}

	return fmt.Errorf("close app: %w", errors.Join(closeErrs...))
}

func (a *App) Run(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("app not exists")
	}
	if a.server == nil {
		return fmt.Errorf("server not exists")
	}

	// start server
	serverErrCh := a.server.Start()

	// wait for cancel or server error
	select {
	case <-ctx.Done():
		a.log.Info("app received cancel signal")
		return nil
	case srvErr := <-serverErrCh:
		return fmt.Errorf("server running: %v", srvErr)
	}
}

func (app *App) init(initers []Initer) error {
	for _, initer := range initers {
		closer, initErr := initer(context.TODO())
		if initErr == nil {
			app.closers = append(app.closers, closer)
			continue
		}
		if closeErr := app.Close(); closeErr != nil {
			return errors.Join(initErr, closeErr)
		}
		return initErr
	}
	return nil
}

func (app *App) initXRayCfg(ctx context.Context) (Closer, error) {
	xrayCfg, err := xraycfg.New(app.cfg.XRayServerConfigPath)
	if err != nil {
		return nil, fmt.Errorf("init xray cfg: %w", err)
	}
	app.xrayCfg = xrayCfg
	return nil, nil
}

func (app *App) initXRayCtl(ctx context.Context) (Closer, error) {
	xrayCtl, err := xrayctl.New(
		app.cfg.XRayExecPath,
		app.cfg.XRayServerConfigPath,
		app.log,
	)
	if err != nil {
		return nil, fmt.Errorf("init xray ctl: %w", err)
	}
	app.xrayCtl = xrayCtl
	return app.closeXRayCtl, nil
}

func (app *App) closeXRayCtl(ctx context.Context) error {
	if err := app.xrayCtl.Close(ctx); err != nil {
		return fmt.Errorf("close xray ctl: %w", err)
	}
	return nil
}

func (app *App) initXRayAPI(ctx context.Context) (Closer, error) {
	xrayAPI, err := xrayapi.New(
		app.xrayCfg.GetApiURL(),
		app.xrayCfg.GetInbounds(),
	)
	if err != nil {
		return nil, fmt.Errorf("init xray api: %w", err)
	}
	app.xrayAPI = xrayAPI
	return app.closeXRayAPI, nil
}

func (app *App) closeXRayAPI(ctx context.Context) error {
	if err := app.xrayAPI.Close(); err != nil {
		return fmt.Errorf("close xray api: %w", err)
	}
	return nil
}

func (app *App) initService(ctx context.Context) (Closer, error) {
	service, err := service.New(app.xrayCfg, app.xrayAPI, app.xrayCtl)
	if err != nil {
		return nil, fmt.Errorf("init service: %w", err)
	}
	app.service = service
	return nil, nil
}

func (app *App) initHandlers(ctx context.Context) (Closer, error) {
	handler, err := router.New(app.log, app.service)
	if err != nil {
		return nil, fmt.Errorf("init handler: %v", err)
	}
	app.handler = handler
	return nil, nil
}

func (app *App) initServer(ctx context.Context) (Closer, error) {
	app.server = server.New(app.cfg.Endpoint, app.handler)
	return app.closeServer, nil
}

func (app *App) closeServer(ctx context.Context) error {
	if err := app.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("close server: %w", err)
	}
	return nil
}
