package main

import (
	"fmt"
	stdlog "log"

	"github.com/XRay-Addons/xrayman/node/internal/app"
	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/logging"
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

func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading: %v", err)
	}
	if err = config.Validate(*cfg); err != nil {
		return nil, fmt.Errorf("config validation: %v", err)
	}
	return cfg, nil
}
