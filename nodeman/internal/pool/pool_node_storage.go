package pool

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type poolNodeUoW struct {
	base   UoW
	nodeID models.NodeID
}

var _ NodeUoW = (*poolNodeUoW)(nil)

func (uow *poolNodeUoW) Do(ctx context.Context, fn NodeUoWFn) error {
	return uow.base.Do(ctx, func(uowctx UoWContext) error {
		nodeUoWCtx := &PoolNodeUoWContext{
			base:   uowctx,
			nodeID: uow.nodeID,
		}
		if err := fn(nodeUoWCtx); err != nil {
			return errdefs.WrapWithf(err, "node: %v", uow.nodeID)
		}
		return fn(nodeUoWCtx)
	})
}

type PoolNodeUoWContext struct {
	base   UoWContext
	nodeID models.NodeID
}

var _ NodeUoWContext = (*PoolNodeUoWContext)(nil)

func (c *PoolNodeUoWContext) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return c.base.FetchNodeStatus(ctx, c.nodeID)
}

func (c *PoolNodeUoWContext) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	return c.base.FindPendingSyncs(ctx, c.nodeID)
}

func (c *PoolNodeUoWContext) ListUsers(ctx context.Context) ([]models.User, error) {
	return c.base.ListUsers(ctx)
}

func (c *PoolNodeUoWContext) PatchPendingSyncs(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return c.base.PatchPendingSyncs(ctx, c.nodeID, patch)
}

func (c *PoolNodeUoWContext) UpdateClientConfig(ctx context.Context,
	cfg models.ClientConfig,
) error {
	return c.base.UpdateClientConfig(ctx, c.nodeID, cfg)
}

func (c *PoolNodeUoWContext) UpdateCurrentStatus(ctx context.Context,
	s models.NodeStatus,
) error {
	return c.base.UpdateCurrentStatus(ctx, c.nodeID, s)
}
