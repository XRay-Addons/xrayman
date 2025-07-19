package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handlers"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	"github.com/XRay-Addons/xrayman/node/internal/infra/httpserver"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xraycfg"
	"github.com/XRay-Addons/xrayman/node/internal/xray/xrayservice"

	"go.uber.org/zap"
)

type Closer = func(ctx context.Context) error

type App struct {
	server  *httpserver.Server
	log     *zap.Logger
	closers []Closer
}

func New(cfg config.Config, log *zap.Logger) (app *App, err error) {
	if log == nil {
		return nil, fmt.Errorf("%w: app init: logger", errdefs.ErrNilArgPassed)
	}

	var closers []Closer
	defer func() {
		if err == nil {
			return
		}
		if closerErr := execClosers(context.TODO(), closers); closerErr != nil {
			err = fmt.Errorf("%w; close error: %w", err, closerErr)
		}
	}()

	// init srv config
	serverCfg, err := xraycfg.NewServerCfg(cfg.XRayServer())
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	// init client config
	clientCfg, err := xraycfg.NewClientCfg(cfg.XRayClient())
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	// init service and add service closer
	xrayService, err := xrayservice.New(cfg.XRayExec(), cfg.XRayServer(), log)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}
	closers = append(closers, func(ctx context.Context) error {
		return xrayService.Close(ctx)
	})

	// init api and add api closer
	xrayAPI, err := xrayapi.New(serverCfg.GetApiURL(), serverCfg.GetInbounds(), log)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}
	closers = append(closers, func(ctx context.Context) error {
		return xrayAPI.Close(ctx)
	})

	// init service
	service, err := service.New(serverCfg, clientCfg, xrayService, xrayAPI)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	// init handlers
	handlers, err := handlers.New(service)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	// init router
	router, err := router.New(cfg.AccessKey, handlers, log)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	// init server
	server := httpserver.New(cfg.Endpoint, router)
	closers = append(closers, func(ctx context.Context) error {
		return server.Shutdown(ctx)
	})

	return &App{
		server:  server,
		closers: closers,
		log:     log,
	}, nil
}

func (app *App) Close(ctx context.Context) error {
	if app == nil {
		return nil
	}
	if err := execClosers(ctx, app.closers); err != nil {
		return fmt.Errorf("app close: %w", err)
	}
	return nil
}

func (a *App) Run() error {
	if a == nil {
		return fmt.Errorf("%w: app", errdefs.ErrNilObjectCall)
	}
	return a.server.Start()
}

func execClosers(ctx context.Context, closers []Closer) error {
	var closeErrs []error
	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i](ctx); err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("closer: %w", err))
		}
	}
	if len(closeErrs) > 0 {
		return errors.Join(closeErrs...)
	}
	return nil
}
