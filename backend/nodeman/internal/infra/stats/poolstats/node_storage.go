package poolstats

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	node "github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/nodestats"
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
			return xerr.WrapWithInfof(err, "node: %v", s.nodeID)
		}
		return nil
	})
}

type PoolNodeUoWContext struct {
	base   UoWContext
	nodeID models.NodeID
}

var _ node.UoWContext = (*PoolNodeUoWContext)(nil)

// UpdateNodeStats implements nodestats.UoWContext.
func (c *PoolNodeUoWContext) UpdateNodeStats(ctx context.Context, s models.NodeStats) error {
	return c.base.UpdateNodeStats(ctx, c.nodeID, s)
}
