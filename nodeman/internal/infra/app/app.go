package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/tx"
	"github.com/oklog/run"
	"go.uber.org/zap"
)

type InitFn = func(ctx context.Context) error
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

const (
	defaultCancelTimeout = 5 * time.Second
)

func New(opts ...Option) *App {
	app := &App{
		log:           zap.NewNop(),
		cancelTimeout: defaultCancelTimeout,
	}
	for _, o := range opts {
		o(app)
	}
	return app
}

func (app *App) Run() (err error) {
	// CTRL + C - cancelable context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel();
	app.log.Info("App starting, press Ctrl+C to cancel...")

	// init app components
	if err := app.init(ctx); err != nil {
		return err
	}
	defer func() {
		if closeErr := app.close(); closeErr != nil {
			err = errors.Join(closeErr, err)
		}
	}()

	app.log.Info("App started, press Ctrl+C to cancel...")
	if err := app.run(ctx); err != nil {
		return err
	}

	return nil
}

type component struct {
	init  func(ctx context.Context) error
	close func(ctx context.Context) error
}

type runner struct {
	run   func() error
	close func() error
}

func (app *App) initComponenFn(init InitFn, name string) func(ctx context.Context) error {
	return func(ctx context.Context) (err error) {
		if init == nil {
			return nil
		}
		defer func() { logOp(app.log, name, "init", err) }()
		return init(ctx)
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

func (app *App) stopRunnerFn(close CloseFn, name string) func() error {
	return func() (err error) {
		if close == nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), app.cancelTimeout)
		defer cancel()

		defer func() { logOp(app.log, name, "close", err) }()
		return close(ctx)
	}
}

func logOp(log *zap.Logger, name, op string, err error) {
	if err == nil {
		log.Info("app operation", zap.String("item", name), zap.String("op", op))
	} else {
		log.Error("app operation", zap.String("item", name), zap.String("op", op), zap.Error(err))
	}
}

func (app *App) init(ctx context.Context) error {
	// init app as transaction.
	// tx rollback time limited, init - unlimited
	initTx := tx.New(
		tx.WithRollbackTimeout(app.cancelTimeout),
	)
	for _, c := range app.components {
		initTx.AddItem(
			func(ctx context.Context) error { return c.init(ctx) },
			func(ctx context.Context) error { return c.close(ctx) },
		)
	}

	if err := initTx.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (app *App) close() error {
	var closeErrs []error
	for i := len(app.components) - 1; i >= 0; i-- {
		ctx, cancel := context.WithTimeout(context.Background(), app.cancelTimeout)
		err := app.components[i].close(ctx)
		cancel()

		if err != nil {
			closeErrs = append(closeErrs, err)
		}
	}

	if len(closeErrs) == 0 {
		return nil
	}

	return errors.Join(closeErrs...)
}

func (app *App) run(ctx context.Context) error {
	var g run.Group

	// collect runners errors
	var runErrs = make([]error, len(app.runners))
	var closeErrs = make([]error, len(app.runners))
	for i, r := range app.runners {
		g.Add(
			func() error {
				runErrs[i] = r.run()
				return runErrs[i]
			},
			func(err error) {
				closeErrs[i] = r.close()
			},
		)
	}


	done := make(chan struct{})
	g.Add(
		func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		},
		func(error) {
			close(done)
		},
	)

	// all errors collected, ignore this (it's first of runner errors)
	_ = g.Run()

	// join all errors
	allErrors := append(runErrs, closeErrs...)
	return errors.Join(allErrors...)
}
