package app

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/xray/api"
	"github.com/XRay-Addons/xrayman/node/internal/xray/clientcfg"
	"github.com/XRay-Addons/xrayman/node/internal/xray/servercfg"
	"github.com/XRay-Addons/xrayman/node/internal/xray/servicectl"
	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: app init: logger", errdefs.ErrNilArgPassed)
	}

	srvCfg, err := servercfg.New(cfg.XRayServer())
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}
	clientCfg, err := clientcfg.New(cfg.XRayClient(), "", "")
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}
	xrayService, err := servicectl.New(cfg.XRayExec(), cfg.XRayServer(), log)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}
	xrayAPI, err := api.New(srvCfg.GetApiURL(), srvCfg.GetInbounds(), log)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	return &App{log: log}
}

func (app *App) Close() error {
	return nil
}

func (a *App) Run(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("%w: app", errdefs.ErrNilObjectCall)
	}

	/*testExecPath := "/usr/local/bin/xrayman/xray"
	testCfgPath := "/usr/local/bin/xrayman/server.json"

	// init service
	xrayCtl, err := launchctl.New(testExecPath, testCfgPath, a.log)
	if err != nil {
		return fmt.Errorf("run app: %w", err)
	}
	defer xrayCtl.Close(ctx)

	// status 5 times
	for range 2 {
		status, err := xrayCtl.Status(ctx)
		if err != nil {
			a.log.Warn("status attempt", zap.Error(err))
		} else {
			a.log.Info(fmt.Sprintf("status request: %s", status))
		}
		time.Sleep(1 * time.Second)
	}

	// start 5 times
	for range 2 {
		if err := xrayCtl.Start(ctx); err != nil {
			a.log.Warn("start attempt", zap.Error(err))
		} else {
			a.log.Info("service started")
		}
		time.Sleep(1 * time.Second)
	}

	// status 5 times
	for range 2 {
		status, err := xrayCtl.Status(ctx)
		if err != nil {
			a.log.Warn("status attempt", zap.Error(err))
		} else {
			a.log.Info(fmt.Sprintf("status request: %s", status))
		}
		time.Sleep(1 * time.Second)
	}

	// stop 5 times
	for range 2 {
		err := xrayCtl.Stop(ctx)
		if err != nil {
			a.log.Warn("stop attempt", zap.Error(err))
		} else {
			a.log.Info("service stopped")
		}
		time.Sleep(1 * time.Second)
	}

	// status 5 times
	for range 2 {
		status, err := xrayCtl.Status(ctx)
		if err != nil {
			a.log.Warn("status attempt", zap.Error(err))
		} else {
			a.log.Info(fmt.Sprintf("status request: %s", status))
		}
		time.Sleep(1 * time.Second)
	}

	<-ctx.Done()*/
	return nil
}
