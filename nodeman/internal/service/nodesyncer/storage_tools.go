package nodesyncer

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func fetchStatus(ctx context.Context, storage Storage) (
	target, current models.NodeStatus, err error,
) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		var innerErr error
		target, current, innerErr = uowctx.NodeStatusStorage().FetchNodeStatus(ctx)
		return innerErr
	})

	if err != nil {
		err = fmt.Errorf("fetch node status: %w", err)
		return
	}

	return
}

func listUsers(ctx context.Context, storage Storage) (
	users []models.UserTargetState, err error,
) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		var innerErr error
		users, innerErr = uowctx.UsersStorage().ListUsers(ctx)
		return innerErr
	})

	if err != nil {
		err = fmt.Errorf("list managed users: %w", err)
		return
	}

	return
}

func updateStatus(ctx context.Context, storage Storage,
	status models.NodeStatus) (err error,
) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		return uowctx.NodeStatusStorage().UpdateCurrentStatus(ctx, status)
	})

	if err != nil {
		err = fmt.Errorf("update current node status: %w", err)
		return
	}

	return
}

func findPendingSyncs(ctx context.Context,
	storage Storage,
) (syncs []models.UserSyncStatus, err error) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		var innerErr error
		syncs, innerErr = uowctx.PendingSyncsStorage().FindPendingSyncs(ctx)
		return innerErr
	})

	if err != nil {
		err = fmt.Errorf("find pending syncs: %w", err)
		return
	}

	return
}

func patchPendingSyncs(ctx context.Context,
	storage Storage, patch []models.UserStatusPatch,
) (err error) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		return uowctx.PendingSyncsStorage().PatchPendingSyncs(ctx, patch)
	})

	if err != nil {
		err = fmt.Errorf("patch pending syncs: %w", err)
		return
	}

	return
}

func patchAndUpdateStatus(ctx context.Context,
	storage Storage, patch []models.UserStatusPatch,
	status models.NodeStatus,
) (err error) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		if e := uowctx.NodeStatusStorage().UpdateCurrentStatus(ctx, status); e != nil {
			return e
		}
		if e := uowctx.PendingSyncsStorage().PatchPendingSyncs(ctx, patch); e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		err = fmt.Errorf("patch and update status: %w", err)
		return
	}

	return
}

func nodeFullUpdate(ctx context.Context, storage Storage,
	patch []models.UserStatusPatch,
	status models.NodeStatus,
	cfg *models.ClientConfig,
) (err error) {
	err = storage.DoUoW(ctx, func(uowctx UoWContext) error {
		if e := uowctx.NodeStatusStorage().UpdateCurrentStatus(ctx, status); e != nil {
			return e
		}
		if e := uowctx.PendingSyncsStorage().PatchPendingSyncs(ctx, patch); e != nil {
			return e
		}
		if e := uowctx.NodeConfigStorage().UpdateClientConfig(ctx, cfg); e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		err = fmt.Errorf("node full update: %w", err)
		return
	}

	return
}
