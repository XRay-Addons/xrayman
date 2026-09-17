package gx

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MakeTickJob[T zapcore.ObjectMarshaler](
	name string,
	fn func(context.Context) (*T, error),
	interval time.Duration,
	log *zap.Logger,
) Job {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	iterFn := func() {
		iterCtx, iterCancel := context.WithTimeout(ctx, interval)
		defer iterCancel()
		result, err := fn(iterCtx)
		if err != nil {
			log.Warn(name, zap.Error(err))
		} else if result != nil {
			log.Info(name, zap.Inline(*result))
		}
	}

	runFn := func() error {
		defer close(done)

		iterFn()

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-timer.C:
				iterFn()
				timer.Reset(interval)
			}
		}
	}

	shutdownFn := func(shutdownCtx context.Context) error {
		cancel()

		select {
		case <-done:
			return nil
		case <-shutdownCtx.Done():
			return shutdownCtx.Err()
		}
	}

	return Job{
		Name:     name,
		Run:      runFn,
		Shutdown: shutdownFn,
	}
}
