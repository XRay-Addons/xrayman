package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) FindPendingSyncs(ctx context.Context,
	id models.NodeID,
) ([]models.UserSyncStatus, error) {
	resp, err := uow.q.FindPendingSyncs(ctx, queries.FindPendingSyncsParams{
		NodeID:            int64(id),
		DefaultUserStatus: int16(models.UserStatusDisabled),
	})
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	syncs, err := ConvertArray[queries.FindPendingSyncsRow, models.UserSyncStatus](resp,
		With(func(from *queries.FindPendingSyncsRow, to *models.UserSyncStatus) {
			to.CurrentStatus = models.UserStatus(from.UserCurrentStatus)
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.VlessUUID = from.VlessUuid
		}),
	)
	if err != nil {
		return nil, err
	}

	return syncs, nil
}

const UserNodeLocksMask = 1 << 33

func (uow *uowctx) SetNodeUsers(ctx context.Context, id models.NodeID,
	patch []models.UserStatusPatch,
) error {
	// remove old users
	// TODO: 3-steps via temp table
	if err := uow.q.Lock(ctx, UserNodeLocksMask|int64(id)); err != nil {
		return err
	}

	if err := uow.q.DeleteNodeUsers(ctx, int64(id)); err != nil {
		return xerr.WrapWithStack(err)
	}

	return uow.UpdateNodeUsers(ctx, id, patch)
}

func (uow *uowctx) UpdateNodeUsers(ctx context.Context, id models.NodeID,
	patch []models.UserStatusPatch,
) error {
	if len(patch) == 0 {
		return nil
	}

	// bulk insert new users
	arg := queries.InsertNodeUsersParams{
		NodeID:            int64(id),
		UserID:            make([]int64, 0, len(patch)),
		UserCurrentStatus: make([]int16, 0, len(patch)),
	}
	for _, p := range patch {
		arg.UserID = append(arg.UserID, int64(p.UserID))
		arg.UserCurrentStatus = append(arg.UserCurrentStatus, int16(p.Status))
	}

	if err := uow.q.Lock(ctx, UserNodeLocksMask|int64(id)); err != nil {
		return err
	}
	if err := uow.q.InsertNodeUsers(ctx, arg); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) GetUserNodes(ctx context.Context, id models.UserID) (
	[]models.Node, error,
) {
	resp, err := uow.q.GetUserNodes(ctx, queries.GetUserNodesParams{
		UserID:            int64(id),
		UserStatusEnabled: int16(models.UserStatusEnabled),
		NodeStatusRunning: int16(models.NodeStatusRunning),
	})
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	nodes, err := ConvertArray[queries.GetUserNodesRow, models.Node](resp,
		With(func(from *queries.GetUserNodesRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
		}),
		WithE(func(from *queries.GetUserNodesRow, to *models.Node) (err error) {
			err = to.Config.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
			return
		}),
		WithE(func(from *queries.GetUserNodesRow, to *models.Node) (err error) {
			err = to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
			return
		}),
	)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
