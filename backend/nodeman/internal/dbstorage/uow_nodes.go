package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) NewNode(ctx context.Context, node *models.Node) error {
	// pre-convert
	arg, err := Convert[models.Node, queries.NewNodeParams](node,
		With(func(from *models.Node, to *queries.NewNodeParams) {
			to.NodeEndpoint = from.Config.ConnectionInfo.Endpoint
			to.NodeCurrentStatus = int16(from.CurrentStatus)
			to.NodeTargetStatus = int16(from.TargetStatus)
		}),
		WithE(func(from *models.Node, to *queries.NewNodeParams) (err error) {
			to.ClientCfgTemplate, err = from.Config.ClientConfigTemplate.Value()
			return
		}),
		WithE(func(from *models.Node, to *queries.NewNodeParams) (err error) {
			to.NodeAccessKey, err = from.Config.ConnectionInfo.AccessKey.Value()
			return
		}),
	)
	if err != nil {
		return err
	}

	// request
	nodeID, err := uow.q.NewNode(ctx, *arg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// post-convert
	node.ID = models.NodeID(nodeID)

	return nil
}

func (uow *uowctx) GetNode(ctx context.Context,
	id models.NodeID,
) (*models.Node, error) {
	// request
	resp, err := uow.q.GetNode(ctx, int64(id))
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	// post-convert
	node, err := Convert[queries.GetNodeRow, models.Node](&resp,
		With(func(from *queries.GetNodeRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
		}),
		WithE(func(from *queries.GetNodeRow, to *models.Node) error {
			return to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
		}),
		WithE(func(from *queries.GetNodeRow, to *models.Node) error {
			return to.Config.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
		}),
	)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (uow *uowctx) ListNodes(ctx context.Context) (
	[]models.Node, error,
) {
	// request
	resp, err := uow.q.ListNodes(ctx)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	// post-convert
	nodes, err := ConvertArray[queries.ListNodesRow, models.Node](resp,
		With(func(from *queries.ListNodesRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
		}),
		WithE(func(from *queries.ListNodesRow, to *models.Node) error {
			return to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
		}),
		WithE(func(from *queries.ListNodesRow, to *models.Node) error {
			return to.Config.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
		}),
	)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (uow *uowctx) SetTargetNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	if err := uow.q.SetTargetNodeStatus(ctx, queries.SetTargetNodeStatusParams{
		NodeID:           int64(id),
		NodeTargetStatus: int16(status),
	}); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (uow *uowctx) SetCurrentNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	err := uow.q.SetCurrentNodeStatus(ctx, queries.SetCurrentNodeStatusParams{
		NodeID:            int64(id),
		NodeCurrentStatus: int16(status),
	})
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (uow *uowctx) SetClientConfig(ctx context.Context,
	id models.NodeID, cfg models.ClientConfigTemplate,
) (err error) {
	// pre-convert
	var arg queries.SetClientConfigParams
	arg.NodeID = int64(id)
	if arg.ClientCfgTemplate, err = cfg.Value(); err != nil {
		return
	}

	// request
	if err = uow.q.SetClientConfig(ctx, arg); err != nil {
		return xerr.WrapWithStack(err)
	}
	return
}

func (uow *uowctx) DeleteNode(ctx context.Context,
	id models.NodeID,
) error {
	if err := uow.q.DeleteNode(ctx, int64(id)); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
