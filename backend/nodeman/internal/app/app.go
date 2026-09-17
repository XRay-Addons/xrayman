package app

import (
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"

	"go.uber.org/zap"
)

type App struct {
	core *gx.App
}

const JWTIssuer = "nodeman"

const cancelTimeout = 5 * time.Second

func New(cfg *config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, errdefs.NilArg("log")
	}
	srcProvider := gx.Options(
		gx.Supply(cfg),
		gx.WithLogger(log),
		gx.WithShutdownTimeout(cancelTimeout))

	appcore := gx.New(
		srcProvider,

		Config,
		Security,
		Storage,
		Nodes,
		Services,
		HttpServer,
		MetricsServer,

		Bootstrap,
		TickJobs,
		Startup,
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
