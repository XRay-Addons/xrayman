package main

import (
	"fmt"
	stdlog "log"

	"github.com/XRay-Addons/xrayman/common/logging"
	"github.com/XRay-Addons/xrayman/node/internal/app"
	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/version"
	"go.uber.org/zap"
)

func main() {
	cli, err := config.LoadCLI()
	if err != nil {
		stdlog.Fatal(err)
	}

	if cli.Version {
		fmt.Println(version.String())
		return
	}

	log := logging.New(cli.LogLevel)
	defer func() {
		if err := log.Sync(); err != nil {
			stdlog.Print(err)
		}
	}()
	log.Warn("app config", zap.Inline(cli))

	cfg, err := config.NewConfig(cli)
	if err != nil {
		log.Error("config loading", zap.Error(err))
		return
	}

	if err = config.Validate(cfg); err != nil {
		log.Error("config validation", zap.Error(err))
		return
	}

	log.Warn(fmt.Sprintf("Starting app with config: %v...", cfg))

	app, err := app.New(cfg, log)
	if err != nil {
		log.Error("app init", zap.Error(err))
		return
	}
	err = app.Run()
	if err != nil {
		log.Error("app run", zap.Error(err))
		return
	}
}
