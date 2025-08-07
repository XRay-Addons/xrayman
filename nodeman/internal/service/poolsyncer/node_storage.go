package poolsyncer

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
)

type NodeUoW struct {
	base   UoW
	nodeID models.NodeID
}

var _ nodesyncer.UoW = (*NodeUoW)(nil)

func (uow *NodeUoW) Do(ctx context.Context, fn nodesyncer.UoWFn) error {
	return uow.base.Do(ctx, func(uowctx UoWContext) error {
		nodeUoWCtx := &NodeUoWContext{
			base:   uowctx,
			nodeID: uow.nodeID,
		}
		if err := fn(nodeUoWCtx); err != nil {
			return fmt.Errorf("node %v: %w", uow.nodeID, err)
		}
		return fn(nodeUoWCtx)
	})
}

type NodeUoWContext struct {
	base   UoWContext
	nodeID models.NodeID
}

var _ nodesyncer.UoWContext = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return c.base.FetchNodeStatus(ctx, c.nodeID)
}

func (c *NodeUoWContext) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	return c.base.FindPendingSyncs(ctx, c.nodeID)
}

func (c *NodeUoWContext) ListUsers(ctx context.Context) ([]models.User, error) {
	return c.base.ListUsers(ctx)
}

func (c *NodeUoWContext) PatchPendingSyncs(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return c.base.PatchPendingSyncs(ctx, c.nodeID, patch)
}

func (c *NodeUoWContext) UpdateClientConfig(ctx context.Context,
	cfg models.ClientConfig,
) error {
	return c.base.UpdateClientConfig(ctx, c.nodeID, cfg)
}

func (c *NodeUoWContext) UpdateCurrentStatus(ctx context.Context,
	s models.NodeStatus,
) error {
	return c.base.UpdateCurrentStatus(ctx, c.nodeID, s)
}
