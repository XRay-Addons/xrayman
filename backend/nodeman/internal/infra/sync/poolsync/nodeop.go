package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

// node op impl
type nodeOp struct {
	storage Storage
	client  Client
}

var _ poolop.NodeOp = (*nodeOp)(nil)

func (op *nodeOp) Exec(ctx context.Context, node models.Node, log *zap.Logger) error {
	nodeStorage := &nodeStorage{
		base:   op.storage,
		nodeID: node.ID,
	}
	nodeClient, err := op.client.GetNodeClient(node.Config.ConnectionInfo)
	if err != nil {
		return err
	}
	if err := nodesync.SyncState(ctx, nodeClient, nodeStorage); err != nil {
		return err
	}
	return nil
}
