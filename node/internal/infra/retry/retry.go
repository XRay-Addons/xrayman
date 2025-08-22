package retry

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

type Fn = func(ctx context.Context) error

func Retry(ctx context.Context, fn Fn, delays []time.Duration) error {
	var err error
	if err = fn(ctx); err == nil {
		return nil
	}

	for _, delay := range delays {
		select {
		case <-time.After(delay):
			if err = fn(ctx); err == nil {
				return nil
			}
		case <-ctx.Done():
			return err
		}
	}

	return err
}

func RetryInfinite(ctx context.Context, fn Fn, delay time.Duration) (err error) {
	if err = fn(ctx); err == nil {
		return
	}

	for {
		select {
		case <-time.After(delay):
			if err = fn(ctx); err == nil {
				return
			}
		case <-ctx.Done():
			return errdefs.WrapWith(err, "retrying cancelled, this is the last error")
		}
	}
}
