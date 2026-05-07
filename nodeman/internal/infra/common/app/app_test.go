package app

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestApp_New(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),
	)
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("close", zap.Error(err))
		}
	}()

	require.NotNil(t, app)
	require.Len(t, app.closers, 0)
	require.Len(t, app.runners, 0)
}

func TestApp_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),
	)
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("close", zap.Error(err))
		}
	}()

	// init
	app.AddCloser(func(ctx context.Context) error {
		logger.Info("init closer 1")
		return nil
	})

	// bootstrap
	app.AddBootstrap("bs1", func(ctx context.Context) error {
		logger.Info("bootstrap 1")
		return nil
	}, nil)
	app.AddCloser(func(ctx context.Context) error {
		logger.Info("bootstrap closer 1")
		return nil
	})
	app.AddBootstrap("bs2", func(ctx context.Context) error {
		logger.Info("bootstrap 2")
		return nil
	}, nil)

	// runners
	app.AddRunner("run1", func() error {
		logger.Info("runner 1")
		return nil
	}, func(context.Context) error {
		logger.Info("runner closer 1")
		return nil
	})
	app.AddRunner("run2", func() error {
		logger.Info("runner 2")
		return nil
	}, func(context.Context) error {
		logger.Info("runner closer 2")
		return nil
	})
	err := app.Run()

	require.NoError(t, err)
}

func TestApp_BootstrapFail(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),
	)
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("close", zap.Error(err))
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGINT)
	}()

	bootstrapErr := xerr.New("init timeout error")
	bootstrapFn := func(ctx context.Context) error {
		select {
		case <-time.After(20 * time.Second):
			return nil
		case <-ctx.Done():
			logger.Error("bootstrap error")
			return bootstrapErr
		}
	}

	app.AddBootstrap("bs1", bootstrapFn, nil)
	err := app.Bootstrap()
	require.ErrorIs(t, err, bootstrapErr)
}

func TestApp_RunFail(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),
	)
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("close", zap.Error(err))
		}
	}()

	// cancel (after 4 seconds)
	go func() {
		time.Sleep(4 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGINT)
	}()

	// bootstrap (2 seconds)
	bootstrapErr := xerr.New("init timeout error")
	bootstrapFn := func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return nil
		case <-ctx.Done():
			logger.Error("bootstrap error")
			return bootstrapErr
		}
	}

	app.AddBootstrap("bs1", bootstrapFn, nil)

	// run (10 seconds)
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()
	runFn := func() error {
		select {
		case <-time.After(10 * time.Second):
			return nil
		case <-runCtx.Done():
			logger.Info("run op cancelled")
			return nil
		}
	}
	stopFn := func(context.Context) error {
		logger.Info("stop fn")
		runCancel()
		return nil
	}
	app.AddRunner("run1", runFn, stopFn)

	err := app.Run()
	logger.Error("run error", zap.Error(err))
}
