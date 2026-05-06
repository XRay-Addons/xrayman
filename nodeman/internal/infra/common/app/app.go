package aapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type runFn = func() error
type closeFn = func(context.Context) error

type bootstrapFn = func(context.Context) error
type retryFn = func(error) bool

type bootstrap struct {
	name  string
	fn    bootstrapFn
	retry retryFn
}

type runner struct {
	name  string
	run   runFn
	close closeFn
}

type App struct {
	log           *zap.Logger
	cancelTimeout time.Duration

	closers    []closeFn
	bootstraps []bootstrap
	runners    []runner

	ctx    context.Context
	cancel context.CancelFunc
}

type Option func(app *App)

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

	app.ctx, app.cancel = signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	app.log.Info("App starting, press Ctrl+C to cancel...")

	return app
}

func (a *App) Close() error {
	defer a.cancel()

	var closeErrs []error
	for i := len(a.closers) - 1; i >= 0; i-- {
		ctx, cancel := context.WithTimeout(context.Background(), a.cancelTimeout)
		err := a.closers[i](ctx)
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

func (a *App) AddBootstrap(name string,
	fn func(context.Context) error, retry func(error) bool,
) {
	a.bootstraps = append(a.bootstraps, bootstrap{
		name:  name,
		fn:    fn,
		retry: retry,
	})
}

func (a *App) AddCloser(c func(context.Context) error) {
	if c != nil {
		a.closers = append(a.closers, c)
	}
}

func (a *App) AddRunner(name string,
	r func() error, c func(context.Context) error,
) {
	a.runners = append(a.runners, runner{
		name:  name,
		run:   r,
		close: c,
	})
}

func (a *App) Bootstrap() (err error) {
	const retryInterval = 1000 * time.Millisecond
	backoff := retry.NewConstant(retryInterval)

	for _, bs := range a.bootstraps {
		if err = retry.Do(a.ctx, backoff, func(ctx context.Context) error {
			err := bs.fn(a.ctx)
			if err != nil && bs.retry != nil && bs.retry(err) {
				a.log.Warn(fmt.Sprintf("bootstrap %s: retry", bs.name), zap.Error(err))
				return retry.RetryableError(err)
			}
			return err
		}); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run() (err error) {
	// CTRL + C - cancelable context
	a.log.Info("App started, press Ctrl+C to cancel...")

	var g run.Group

	// run runners, collect errors
	var runErrs = make([]error, len(a.runners))
	var closeErrs = make([]error, len(a.runners))
	for i, r := range a.runners {
		runFn := func() error {
			a.log.Info(fmt.Sprintf("run %s...", r.name))
			if r.run == nil {
				return nil
			}
			runErrs[i] = r.run()
			return runErrs[i]
		}
		closeFn := func(error) {
			a.log.Info(fmt.Sprintf("stop %s...", r.name))
			if r.close == nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), a.cancelTimeout)
			closeErrs[i] = r.close(ctx)
			cancel()
		}
		g.Add(runFn, closeFn)
	}

	// wait for cancel
	done := make(chan struct{})
	g.Add(
		func() error {
			select {
			case <-a.ctx.Done():
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
