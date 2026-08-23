package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) NewNode(ctx context.Context, node *models.Node) error {
	// pre-convert
	arg, err := convert.NewNodeReq(node)
	if err != nil {
		return err
	}

	// request
	nodeID, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (int64, error) {
		return q.NewNode(ctx, *arg)
	})
	if err != nil {
		return err
	}

	// post-convert
	node.ID = models.NodeID(nodeID)

	return nil
}

func (s *Storage) GetNode(ctx context.Context,
	id models.NodeID,
) (*models.Node, error) {
	// request
	node, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (queries.GetNodeRow, error) {
		return q.GetNode(ctx, int64(id))
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.GetNodeResp(&node)
}

func (s *Storage) ListNodes(ctx context.Context) (
	[]models.Node, error,
) {
	// request
	nodes, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries) ([]queries.ListNodesRow, error,
	) {
		nodes, err := q.ListNodes(ctx)
		return nodes, err
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.ListNodesResp(nodes)
}

func (s *Storage) SetTargetNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetTargetNodeStatus(ctx, queries.SetTargetNodeStatusParams{
			NodeID:           int64(id),
			NodeTargetStatus: int16(status),
		})
	})
}

func (s *Storage) SetCurrentNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetCurrentNodeStatus(ctx, queries.SetCurrentNodeStatusParams{
			NodeID:            int64(id),
			NodeCurrentStatus: int16(status),
		})
	})
}

func (s *Storage) SetNodeSettings(ctx context.Context,
	id models.NodeID, settings *models.NodeSettings,
) error {
	// pre-convert
	arg, err := convert.SetNodeSettingsReq(id, settings)
	if err != nil {
		return err
	}
	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetNodeSettings(ctx, *arg)
	})
}

func (s *Storage) DeleteNode(ctx context.Context,
	id models.NodeID,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.DeleteNode(ctx, int64(id))
	})
}

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
