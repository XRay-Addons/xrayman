package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) GetNodeMetrics(ctx context.Context) ([]models.NodeMetrics, error) {
	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.GetMetricsRow, error) {
		return q.GetMetrics(ctx, int16(models.NodeStatusRunning))
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.GetNodeMetricsResp(resp)
}
