package dbstorage

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) UpdateNodeStats(ctx context.Context,
	nodeID models.NodeID, stats models.NodeStats,
) error {
	// pre-convert
	args := convert.UpdateNodeStatsReq(nodeID, stats)

	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.UpdateTotalStats(ctx, args)
	})
}

func (s *Storage) UpdateDailyStats(ctx context.Context,
	day time.Time,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.UpdateDailyStats(ctx, day)
	})
}
