package gx

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/logging"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type CallsCounter struct {
	current atomic.Int32
}

func (cc *CallsCounter) Call() {
	cc.current.Add(1)
}

func (cc *CallsCounter) Check(t *testing.T, name string, target int32) {
	require.Equal(t, target, cc.current.Load(),
		fmt.Sprintf("%s calls count check", name))
}

func sendSIGINT(t *testing.T, log *zap.Logger, delay time.Duration) <-chan struct{} {
	t.Helper()

	done := make(chan struct{})

	go func() {
		defer close(done)

		time.Sleep(delay)

		p, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)
		log.Warn("sent SIGINT")
		err = p.Signal(syscall.SIGINT)
		require.NoError(t, err)
	}()

	return done
}

type A struct{}
type B struct{}
type C struct{}

func TestGx_Simple(t *testing.T) {
	logger := logging.New(zapcore.InfoLevel)

	var aStart, bStart, cStart, cStop, iI, iIstop, iII, iIIstop CallsCounter
	defer aStart.Check(t, "a", 1)
	defer bStart.Check(t, "b", 1)
	defer cStart.Check(t, "c", 1)
	defer cStop.Check(t, "c-stop", 1)
	defer iI.Check(t, "iI", 1)
	defer iIstop.Check(t, "iI-stop", 1)
	defer iII.Check(t, "iII", 1)
	defer iIIstop.Check(t, "iII-stop", 1)

	app := New(
		// logger
		WithLogger(logger),

		// a, b, c - long providers
		fx.Provide(func(lc Lifecycle, l *zap.Logger) (A, error) {
			time.Sleep(3 * time.Second)
			lc.AppendJob(Job{
				Name: "A Hook",
				OnStart: func(context.Context) error {
					aStart.Call()
					time.Sleep(3 * time.Second)
					return nil
				},
			})
			return A{}, nil
		}),
		fx.Provide(func(lc Lifecycle, l *zap.Logger, a A) (B, error) {
			time.Sleep(3 * time.Second)
			lc.AppendCloser(Closer{
				Name: "B Hook",
				OnClose: func(context.Context) error {
					bStart.Call()
					time.Sleep(3 * time.Second)
					return nil
				},
			})
			return B{}, nil
		}),
		fx.Provide(func(lc Lifecycle, l *zap.Logger, a A, b B) (C, error) {
			time.Sleep(3 * time.Second)
			lc.AppendJob(Job{
				Name: "C Hook",
				OnStart: func(context.Context) error {
					cStart.Call()
					time.Sleep(3 * time.Second)
					return nil
				},
				OnStop: func(context.Context) error {
					cStop.Call()
					time.Sleep(3 * time.Second)
					return nil
				},
			})
			return C{}, nil
		}),

		// long invokers
		fx.Invoke(func(lc Lifecycle, l *zap.Logger, c C) error {
			time.Sleep(3 * time.Second)
			lc.AppendJob(Job{
				Name: "Long invoke",
				OnStart: func(context.Context) error {
					iI.Call()
					return nil
				},
				OnStop: func(context.Context) error {
					iIstop.Call()
					time.Sleep(3 * time.Second)
					return nil
				},
			})
			return nil
		}),
		// long invokers
		fx.Invoke(func(lc Lifecycle, l *zap.Logger, c C) error {
			ctx, cancel := context.WithCancel(context.Background())
			lc.AppendJob(Job{
				Name: "Infinite invoke",
				OnStart: func(context.Context) error {
					<-ctx.Done()
					iII.Call()
					return nil
				},
				OnStop: func(context.Context) error {
					iIIstop.Call()
					cancel()
					return nil
				},
			})
			return nil
		}),
	)

	ch := sendSIGINT(t, logger, 30)
	defer func() { <-ch }()
	require.NoError(t, app.Run())
	logger.Warn("app stopped")
}

func TestGx_New(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(WithLogger(logger))
	require.NotNil(t, app)
}

func TestGx_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),

		// bootstrap
		fx.Invoke(func(lc Lifecycle, log *zap.Logger) {
			log.Info("bootsrtap 1")
			lc.AppendCloser(Closer{
				Name: "bs1 closer",
				OnClose: func(context.Context) error {
					log.Info("bootsrtap closer 1")
					return nil
				}})
		}),

		// jobs
		fx.Invoke(func(lc Lifecycle, log *zap.Logger) {
			lc.AppendJob(Job{
				Name: "job 1",
				OnStart: func(context.Context) error {
					log.Info("job 1")
					return nil
				},
				OnStop: func(context.Context) error {
					log.Info("job 1 stop")
					return nil
				},
			})
		}),
		fx.Invoke(func(lc Lifecycle, log *zap.Logger) {
			lc.AppendJob(Job{
				Name: "job 2",
				OnStart: func(context.Context) error {
					log.Info("job 2")
					return nil
				},
				OnStop: func(context.Context) error {
					log.Info("job 2 stop")
					return nil
				},
			})
		}),
	)

	ch := sendSIGINT(t, logger, 30)
	defer func() { <-ch }()
	require.NoError(t, app.Run())
	logger.Warn("app stopped")
}

func TestGx_ProvideFail(t *testing.T) {

	logger := zaptest.NewLogger(t)
	provideErr := xerr.New("provide timeout error")

	app := New(
		WithLogger(logger),

		// providers
		fx.Provide(func(ctx context.Context) (*A, error) {
			select {
			case <-time.After(2 * time.Second):
				return &A{}, nil
			case <-ctx.Done():
				logger.Error("provide error")
				return nil, provideErr
			}
		}),

		// jobs
		fx.Invoke(func(a *A, log *zap.Logger) {
			log.Info("job success")
		}),
	)

	ch := sendSIGINT(t, logger, 1*time.Second)
	defer func() { <-ch }()
	require.ErrorIs(t, app.Run(), provideErr)
	logger.Warn("app stopped")
}

func TestGx_BootstrapFail(t *testing.T) {

	logger := zaptest.NewLogger(t)
	bootstrapErr := xerr.New("init timeout error")

	app := New(
		WithLogger(logger),

		// bootstrap
		fx.Invoke(func(ctx context.Context, log *zap.Logger) error {
			select {
			case <-time.After(2 * time.Second):
				return nil
			case <-ctx.Done():
				logger.Error("bootstrap error")
				return bootstrapErr
			}
		}),

		// jobs
		fx.Invoke(func(lc Lifecycle, log *zap.Logger) {
			lc.AppendJob(Job{
				Name: "job 1",
				OnStart: func(context.Context) error {
					log.Info("job 1")
					return nil
				},
				OnStop: func(context.Context) error {
					log.Info("job 1 stop")
					return nil
				},
			})
		}),
		fx.Invoke(func(lc Lifecycle, log *zap.Logger) {
			lc.AppendJob(Job{
				Name: "job 2",
				OnStart: func(context.Context) error {
					log.Info("job 2")
					return nil
				},
				OnStop: func(context.Context) error {
					log.Info("job 2 stop")
					return nil
				},
			})
		}),
	)

	ch := sendSIGINT(t, logger, 1*time.Second)
	defer func() { <-ch }()
	require.ErrorIs(t, app.Run(), bootstrapErr)
	logger.Warn("app stopped")
}

func TestGx_JobFail(t *testing.T) {
	logger := zaptest.NewLogger(t)

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

	// run (10 seconds)
	runErr := xerr.New("run timeout error")
	runCtx, runCancel := context.WithCancel(context.Background())
	runFn := func(context.Context) error {
		select {
		case <-time.After(10 * time.Second):
			return nil
		case <-runCtx.Done():
			logger.Info("run op cancelled")
			return runErr
		}
	}
	stopFn := func(context.Context) error {
		logger.Info("stop fn")
		runCancel()
		return nil
	}

	app := New(
		WithLogger(logger),
		fx.Invoke(bootstrapFn),
		fx.Invoke(func(lc Lifecycle) {
			lc.AppendJob(Job{
				Name:    "job",
				OnStart: runFn,
				OnStop:  stopFn,
			})
		}),
	)

	// cancel (after 4 seconds)
	ch := sendSIGINT(t, logger, 4*time.Second)
	defer func() { <-ch }()

	require.ErrorIs(t, app.Run(), runErr)
	logger.Warn("app stopped")
}
