package bootstrap

import (
	"context"
	"errors"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type Config struct {
	AdminPassword string
}

func Bootstrap(ctx context.Context, cfg Config, auth AuthService, log *zap.Logger) error {
	if cfg.AdminPassword != "" {
		if err := withRetry(ctx, func(ctx context.Context) error {
			err := auth.SetAdmin(ctx, cfg.AdminPassword)
			if errors.Is(err, errdefs.ErrTemporaryUnavailable) {
				return retry.RetryableError(err)
			}
			return err
		}, log); err != nil {
			return err
		}
	}
	return nil
}

func withRetry(ctx context.Context, fn func(context.Context) error, log *zap.Logger) error {
	// retry with policy till success or cancel
	const inintalRetry = 100 * time.Millisecond
	const maxRetry = 2 * time.Second
	backoff := retry.NewFibonacci(inintalRetry)
	return retry.Do(ctx, backoff, func(ctx context.Context) error {
		if err := fn(ctx); err != nil {
			log.Warn("bootstrap error, retry", zap.Error(err))
			return err
		}
		return nil
	})
}
