package nodestats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"go.uber.org/zap"
)

type stats struct {
	storage Storage
	client  Client
}

func UpdateNodeStats(ctx context.Context, client Client, storage Storage, log *zap.Logger) error {
	if client == nil {
		return errdefs.NilArg("client")
	}
	if storage == nil {
		return errdefs.NilArg("storage")
	}
	stats, err := client.GetStats(ctx)
	if err != nil {
		return err
	}

	if err := storage.DoUoW(ctx, func(uowctx UoWContext) error {
		return uowctx.UpdateNodeStats(ctx, *stats)
	}); err != nil {
		return err
	}

	return err
}
