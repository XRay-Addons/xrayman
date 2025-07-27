package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/node"
)

type NodeStorage struct {
}

var _ node.NodeStorage = (*NodeStorage)(nil)

func (n *NodeStorage) BeginTx() node.StorageTx {
	panic("unimplemented")
}

func (n *NodeStorage) FetchNodeStatus(ctx context.Context) (target models.NodeStatus, current models.NodeStatus, err error) {
	panic("unimplemented")
}

// FindPendingSyncs implements node.NodeStorage.
func (n *NodeStorage) FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error) {
	panic("unimplemented")
}

// ListManagedUsers implements node.NodeStorage.
func (n *NodeStorage) ListManagedUsers(ctx context.Context) ([]models.UserTargetState, error) {
	panic("unimplemented")
}
