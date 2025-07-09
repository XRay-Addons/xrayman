package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/app"
	"github.com/XRay-Addons/xrayman/node/internal/logging"
	"go.uber.org/zap"

	stdlog "log"
)

func main() {
	logger, err := logging.New()
	if err != nil {
		stdlog.Printf("logger init: %v", err)
		return
	}
	defer logger.Sync()

	app := app.New(logger)
	defer func() {
		logger.Info("close app")
		if err := app.Close(); err != nil {
			logger.Error("close app", zap.Error(err))
		}
	}()

	// create gentle cancelling to context
	ctx := gracefulCancellingCtx(logger)

	time.Sleep(10 * time.Second)
	// run app
	if err := app.Run(ctx); err != nil {
		logger.Error("run app", zap.Error(err))
		return
	}

}

func gracefulCancellingCtx(logger *zap.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		logger.Info("Press Ctrl+C to stop the server...")
		sig := <-sigChan
		logger.Info(fmt.Sprintf("interruption signal received: %v, shutting down server...", sig))
		cancel()
	}()
	return ctx
}
