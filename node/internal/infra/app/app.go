package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/infra/tx"
	"github.com/oklog/run"
	"go.uber.org/zap"
)

type InitFn = func() error
type CloseFn = func(ctx context.Context) error
type RunFn = func() error

type App struct {
	components []component
	runners    []runner
	log        *zap.Logger

	cancelTimeout time.Duration
}

type Option func(app *App)

func WithComponent(name string, init InitFn, close CloseFn) Option {
	return func(app *App) {
		app.components = append(app.components, component{
			init:  app.initComponenFn(init, name),
			close: app.closeComponentFn(close, name),
		})
	}
}

func WithRunner(name string, run RunFn, close CloseFn) Option {
	return func(app *App) {
		app.runners = append(app.runners, runner{
			run:   app.runRunnerFn(run, name),
			close: app.stopRunnerFn(close, name),
		})
	}
}

func WithSignalCancel() Option {
	return func(app *App) {
		sigCh := make(chan os.Signal, 1)
		app.runners = append(app.runners, runner{
			run: func() error {
				app.log.Info("Press Ctrl+C to stop the server...")
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				if sig := <-sigCh; sig != nil {
					return fmt.Errorf("received signal: %s", sig)
				}
				return nil
			},
			close: func(error) {
				signal.Stop(sigCh)
				close(sigCh)
			},
		})
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(app *App) {
		if l != nil {
			app.log = l
		}
	}
}

func WithCancelTimeout(timeout time.Duration) Option {
	return func(app *App) {
		app.cancelTimeout = timeout
	}
}

func New(opts ...Option) *App {
	app := &App{
		log:           zap.NewNop(),
		cancelTimeout: 5 * time.Second,
	}
	for _, o := range opts {
		o(app)
	}
	return app
}

func (app *App) Run() error {
	// init app components
	if err := app.init(); err != nil {
		return fmt.Errorf("app run: %w", err)
	}
	defer app.close()

	// run runners
	if err := app.run(); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}

type component struct {
	init  func() error
	close func(ctx context.Context) error
}

type runner struct {
	run   func() error
	close func(error)
}

func (app *App) initComponenFn(init InitFn, name string) func() error {
	return func() (err error) {
		if init == nil {
			return nil
		}
		defer func() { logOp(app.log, name, "init", err) }()
		return init()
	}
}

func (app *App) closeComponentFn(close CloseFn, name string) func(ctx context.Context) error {
	return func(ctx context.Context) (err error) {
		if close == nil {
			return nil
		}
		defer func() { logOp(app.log, name, "close", err) }()

		return close(ctx)
	}
}

func (app *App) runRunnerFn(runner RunFn, name string) func() error {
	return func() (err error) {
		logOp(app.log, name, "start", nil)
		defer func() { logOp(app.log, name, "stopped", err) }()
		return runner()
	}
}

func (app *App) stopRunnerFn(close CloseFn, name string) func(error) {
	return func(error) {
		if close == nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), app.cancelTimeout)
		defer cancel()

		err := close(ctx)
		defer func() { logOp(app.log, name, "close", err) }()
	}
}

func logOp(log *zap.Logger, name, op string, err error) {
	if err == nil {
		log.Info("app operation", zap.String("item", name), zap.String("op", op))
	} else {
		log.Error("app operation", zap.String("item", name), zap.String("op", op), zap.Error(err))
	}
}

func (app *App) init() error {
	// init app as transaction.
	// tx rollback time limited, init - unlimited
	initTx := tx.New(
		tx.WithRollbackTimeout(app.cancelTimeout),
	)
	for _, c := range app.components {
		// love it again
		c := c
		initTx.AddItem(
			func(context.Context) error { return c.init() },
			func(ctx context.Context) error { return c.close(ctx) },
		)
	}

	if err := initTx.Run(context.Background()); err != nil {
		return fmt.Errorf("init app: %w", err)
	}

	return nil
}

func (app *App) close() error {
	var closeErrs []error
	for i := len(app.components) - 1; i >= 0; i-- {
		ctx, cancel := context.WithTimeout(context.Background(), app.cancelTimeout)
		defer cancel()

		if err := app.components[i].close(ctx); err != nil {
			closeErrs = append(closeErrs, err)
		}
	}

	if len(closeErrs) == 0 {
		return nil
	}

	return fmt.Errorf("close app error: %w", errors.Join(closeErrs...))
}

func (app *App) run() error {

	var g run.Group
	for _, runner := range app.runners {
		g.Add(runner.run, runner.close)
	}

	if err := g.Run(); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}
