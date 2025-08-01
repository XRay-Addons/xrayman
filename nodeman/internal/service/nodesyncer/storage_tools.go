package nodesyncer

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func fetchStatus(ctx context.Context, storage NodeStorage) (
	target, current models.NodeStatus, err error,
) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		target, current, err = uowctx.FetchNodeStatus(ctx)
		return
	}); err != nil {
		err = fmt.Errorf("fetch node status: %w", err)
	}
	return
}

func listUsers(ctx context.Context, storage NodeStorage) (
	users []models.User, err error,
) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		users, err = uowctx.ListUsers(ctx)
		return
	}); err != nil {
		err = fmt.Errorf("list managed users: %w", err)
	}
	return
}

func updateCurrentStatus(ctx context.Context, storage NodeStorage,
	status models.NodeStatus) (err error,
) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		return uowctx.UpdateCurrentStatus(ctx, status)
	}); err != nil {
		err = fmt.Errorf("update current node status: %w", err)
	}
	return
}

func findPendingSyncs(ctx context.Context,
	storage NodeStorage,
) (syncs []models.UserSyncStatus, err error) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		syncs, err = uowctx.FindPendingSyncs(ctx)
		return
	}); err != nil {
		err = fmt.Errorf("find pending syncs: %w", err)
	}
	return
}

func patchPendingSyncs(ctx context.Context,
	storage NodeStorage, patch []models.UserStatusPatch,
) (err error) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		return uowctx.PatchPendingSyncs(ctx, patch)
	}); err != nil {
		err = fmt.Errorf("patch pending syncs: %w", err)
		return
	}

	return
}

func updateNode(ctx context.Context, storage NodeStorage,
	patch []models.UserStatusPatch,
	status models.NodeStatus,
	cfg *models.ClientConfig,
) (err error) {
	if err = storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		if status > 0 { // TODO: invalid status?
			if err = uowctx.UpdateCurrentStatus(ctx, status); err != nil {
				return
			}
		}
		if patch != nil {
			if err = uowctx.PatchPendingSyncs(ctx, patch); err != nil {
				return
			}
		}
		if cfg != nil {
			if err = uowctx.UpdateClientConfig(ctx, *cfg); err != nil {
				return
			}
		}
		return
	}); err != nil {
		err = fmt.Errorf("node update: %w", err)
	}
	return
}
