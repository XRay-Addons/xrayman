package main

import (
	"fmt"
	"log"
	stdlog "log"

	"github.com/XRay-Addons/xrayman/common/logging"
	"github.com/XRay-Addons/xrayman/nodeman/internal/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/version"
	"go.uber.org/zap"
)

func main() {
	cli, err := config.LoadCLI()
	if err != nil {
		log.Fatal(err)
	}

	if cli.Version {
		fmt.Println(version.String())
		return
	}

	cfg, err := config.NewConfig(cli)
	if err != nil {
		stdlog.Printf("config loading: %+v", err)
		return
	}
	if err = config.Validate(cfg); err != nil {
		stdlog.Printf("config validation: %+v", err)
		return
	}

	log, err := logging.New(cfg.LogLevel)
	if err != nil {
		stdlog.Print(err)
		return
	}
	defer func() {
		if err := log.Sync(); err != nil {
			stdlog.Print(err)
		}
	}()

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
