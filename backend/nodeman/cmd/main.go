package main

import (
	stdlog "log"

	"github.com/XRay-Addons/xrayman/common/logging"
	"github.com/XRay-Addons/xrayman/nodeman/internal/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"go.uber.org/zap"
)

func main() {
	// load and validate config
	cfg, err := loadConfig()
	if err != nil {
		stdlog.Printf("config loading: %v", err)
		return
	}

	// create log. use std log to log log errors,
	// because who log the log
	log, err := logging.New()
	if err != nil {
		stdlog.Print(err)
		return
	}
	defer func() {
		if err := log.Sync(); err != nil {
			stdlog.Print(err)
		}
	}()

	app, err := app.New(*cfg, log)
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

func loadConfig() (*config.RawConfig, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	if err = config.Validate(*cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
