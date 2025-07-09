package app

import (
	"context"
	"fmt"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/xrayctl/launchctl"

	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
}

func New(log *zap.Logger) *App {
	if log == nil {
		log = zap.NewNop()
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

	testExecPath := "/usr/local/bin/xrayman/xray"
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

	<-ctx.Done()
	return nil
}
