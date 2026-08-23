package poolsync

import (
	"context"

	node "github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type nodeStorage struct {
	base   Storage
	nodeID models.NodeID
}

var _ node.Storage = (*nodeStorage)(nil)

func (n *nodeStorage) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	return n.base.FindPendingSyncs(ctx, n.nodeID)
}

func (n *nodeStorage) GetNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	node, err := n.base.GetNode(ctx, n.nodeID)
	if err != nil {
		return
	}
	target = node.TargetStatus
	current = node.CurrentStatus
	return
}

func (n *nodeStorage) ListUsers(ctx context.Context) (
	[]models.User, error,
) {
	return n.base.ListUsers(ctx)
}

func (n *nodeStorage) SetNodeSettings(ctx context.Context,
	cfg *models.NodeSettings,
) error {
	return n.base.SetNodeSettings(ctx, n.nodeID, cfg)
}

func (n *nodeStorage) SetCurrentNodeStatus(ctx context.Context,
	s models.NodeStatus,
) error {
	return n.base.SetCurrentNodeStatus(ctx, n.nodeID, s)
}

func (n *nodeStorage) SetNodeRev(ctx context.Context, rev models.Revision) error {
	return n.base.SetNodeRev(ctx, n.nodeID, rev)
}

func (n *nodeStorage) DoTx(ctx context.Context, fn node.TxFn) error {
	return n.base.DoTx(ctx, fn)
}
