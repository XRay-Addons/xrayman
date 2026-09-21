package app

import (
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"go.uber.org/zap"
)

type App struct {
	core *gx.App
}

const ShutdownTimeout = 5 * time.Second

func New(cfg *config.Config, log *zap.Logger) (app *App, err error) {
	if log == nil {
		return nil, errdefs.NilArg("log")
	}

	srcProvider := gx.Options(
		gx.Supply(cfg),
		gx.WithLogger(log),
		gx.WithShutdownTimeout(ShutdownTimeout),
	)

	appcore := gx.New(
		srcProvider,

		Config,
		Security,
		Services,
		XRay,
		Performance,
		Server,
		Bootstrap,
	)

	return &App{
		core: &appcore,
	}, nil
}

func (app *App) Run() error {
	if app == nil || app.core == nil {
		return xerr.NilCall()
	}
	return app.core.Run()
}
