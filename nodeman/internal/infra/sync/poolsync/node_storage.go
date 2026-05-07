package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	node "github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type nodeStorage struct {
	base   Storage
	nodeID models.NodeID
}

var _ node.Storage = (*nodeStorage)(nil)

func (s *nodeStorage) DoUoW(ctx context.Context, fn node.UoWFn) error {
	return s.base.DoUoW(ctx, func(uowctx UoWContext) error {
		nodeUoWCtx := &PoolNodeUoWContext{
			base:   uowctx,
			nodeID: s.nodeID,
		}
		if err := fn(nodeUoWCtx); err != nil {
			return xerr.WrapWithf(err, "node: %v", s.nodeID)
		}
		return nil
	})
}

type PoolNodeUoWContext struct {
	base   UoWContext
	nodeID models.NodeID
}

var _ node.UoWContext = (*PoolNodeUoWContext)(nil)

func (c *PoolNodeUoWContext) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	node, exists, err := c.base.GetNode(ctx, c.nodeID)
	if err != nil {
		return
	}
	if !exists {
		err = xerr.New("node not exists")
		return
	}
	target = node.TargetStatus
	current = node.CurrentStatus
	return
}

func (c *PoolNodeUoWContext) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	return c.base.FindPendingSyncs(ctx, c.nodeID)
}

func (c *PoolNodeUoWContext) ListUsers(ctx context.Context) ([]models.User, error) {
	return c.base.ListUsers(ctx)
}

func (c *PoolNodeUoWContext) UpdateNodeUsers(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return c.base.UpdateNodeUsers(ctx, c.nodeID, patch)
}

func (c *PoolNodeUoWContext) SetNodeUsers(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return c.base.SetNodeUsers(ctx, c.nodeID, patch)
}

func (c *PoolNodeUoWContext) SetClientConfig(ctx context.Context,
	cfg models.ClientConfigTemplate,
) error {
	return c.base.SetClientConfig(ctx, c.nodeID, cfg)
}

func (c *PoolNodeUoWContext) SetCurrentNodeStatus(ctx context.Context,
	s models.NodeStatus,
) error {
	return c.base.SetCurrentNodeStatus(ctx, c.nodeID, s)
}
