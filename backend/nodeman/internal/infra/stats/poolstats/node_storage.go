package poolstats

import (
	"context"

	node "github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/nodestats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type nodeStorage struct {
	base   Storage
	nodeID models.NodeID
}

var _ node.Storage = (*nodeStorage)(nil)

func (s *nodeStorage) UpdateNodeStats(ctx context.Context, stats models.NodeStats) error {
	return s.base.UpdateNodeStats(ctx, s.nodeID, stats)
}
