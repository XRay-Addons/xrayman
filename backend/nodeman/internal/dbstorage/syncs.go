package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) SetNodeRev(ctx context.Context,
	id models.NodeID, rev models.Revision,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetNodeRev(ctx, queries.SetNodeRevParams{
			Revision: int64(rev),
			NodeID:   int64(id),
		})
	})
}

func (s *Storage) FindPendingSyncs(ctx context.Context,
	id models.NodeID,
) ([]models.UserSyncStatus, error) {
	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.FindPendingSyncsRow, error) {
		return q.FindPendingSyncs(ctx, int64(id))
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.FindPendingSyncsResp(resp), nil
}

func (s *Storage) GetUserNodes(ctx context.Context, id models.UserID) (
	[]models.Node, error,
) {
	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.GetUserNodesRow, error) {
		return q.GetUserNodes(ctx, int64(id))
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.GetUserNodesResp(resp)
}
