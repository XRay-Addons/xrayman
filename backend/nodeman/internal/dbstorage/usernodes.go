package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) FindPendingSyncs(ctx context.Context,
	rev models.Revision,
) ([]models.UserSyncStatus, error) {
	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.FindPendingSyncsRow, error) {
		return q.FindPendingSyncs(ctx, int64(rev))
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.FindPendingSyncsResp(resp), nil
}

const UserNodeLocksMask = 1 << 33

func (s *Storage) SetNodeUsers(ctx context.Context, id models.NodeID,
	patch []models.UserStatusPatch,
) error {
	return s.DoTx(ctx, func(ctx context.Context) error {
		// lock user-id for this tx and operations like this
		if err := doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
			return q.Lock(ctx, UserNodeLocksMask|int64(id))
		}); err != nil {
			return err
		}
		// delete all related to user id
		if err := doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
			return q.DeleteNodeUsers(ctx, int64(id))
		}); err != nil {
			return err
		}
		// update node users
		return s.UpdateNodeUsers(ctx, id, patch)
	})
}

func (s *Storage) UpdateNodeUsers(ctx context.Context, id models.NodeID,
	patch []models.UserStatusPatch,
) error {
	if len(patch) == 0 {
		return nil
	}

	// pre-convert
	arg := convert.UpdateNodeUsersReq(id, patch)

	// request
	if err := doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.InsertNodeUsers(ctx, arg)
	}); err != nil {
		return err
	}

	return nil
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
