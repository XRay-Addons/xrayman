package gx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

/*
Package 'gx' provides better wrapper around 'fx'.
(g > f)

How it works when you call [Run]:
    1. Init
		- builds the dependency graph;
		- executes all Invoke functions;
		- collects jobs and closers registered in the custom Lifecycle.
    2. Starts all registered jobs by calling Job.OnStart in parallel.
       Execution continues until one of the following happens:
         - all jobs exit successfully;
         - any job returns an error;
         - the application receives an interrupt signal (Ctrl+C, SIGINT, SIGTERM).

    3. Calls Job.OnStop for every registered job (running or already finished),
       then waits until every Job.OnStart function has returned.

    4. Returns a single error created with xerr.Join containing all job start
       and stop errors.

Differences from the original fx lifecycle:
  - Context is provided so you can use it even in Invoke and Provide functions
  - Uses a custom Lifecycle instead of Lifecycle.
  - Lifecycle contains named jobs and closers instead of generic hooks.
  - Jobs have explicit start and stop callbacks.
  - Shutdown is coordinated by the wrapper rather than by App.
*/

// reexport fx tools
type (
	Option     = fx.Option
	In         = fx.In
	Out        = fx.Out
	Annotation = fx.Annotation
)

var (
	Provide    = fx.Provide
	Invoke     = fx.Invoke
	Supply     = fx.Supply
	Options    = fx.Options
	Annotate   = fx.Annotate
	Module     = fx.Module
	As         = fx.As
	Self       = fx.Self
	ParamTags  = fx.ParamTags
	ResultTags = fx.ResultTags
)

// upgraded fx tools
func WithLogger(log *zap.Logger) Option {
	if log == nil {
		return nil
	}
	return Options(
		Invoke(func(lc *lifecycle) { lc.log = log }),
		Supply(log),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}

func WithCancelTimeout(to time.Duration) Option {
	return Invoke(func(lc *lifecycle) {
		lc.closeTimeout = to
	})
}

// app impl
type App struct {
	options []Option
}

func New(opts ...Option) App {
	return App{options: opts}
}

func (a *App) Run() (err error) {
	if a == nil {
		return xerr.NilCall()
	}

	// collect all options
	options := append([]Option{}, a.options...)

	// lifecycle option
	lc := &lifecycle{
		log: zap.NewNop(),
	}
	options = append(options, Supply(Annotate(
		lc,
		As(new(Lifecycle)),
		As(Self()),
	)))

	// interruptable context option
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	options = append(options, Supply(Annotate(
		ctx,
		As(new(context.Context)),
	)))

	app := fx.New(options...)
	defer func() {
		err = xerr.Join(err, lc.Close())
	}()

	if err = app.Err(); err != nil {
		return
	}

	return lc.Run(ctx)
}
