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

func (s *nodeStorage) UpdateStats(ctx context.Context, stats models.NodeStats) error {
	return s.base.UpdateStats(ctx, s.nodeID, stats)
}
