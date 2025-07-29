package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/node"
)

// create Storage proxy to implement node.Storage.
// (pass nodeID in ctor and fix it, implement proxy UoW, proxy UoWCtx + inheritance tricks)

type NodeStorage struct {
	nodeID models.NodeID
	base   Storage
}

var _ node.Storage = (*NodeStorage)(nil)

func (n *NodeStorage) DoUoW(ctx context.Context, fn node.UoWFn) error {
	uow, err := n.NewUoW()
	if err != nil {
		return err
	}
	return uow.Do(ctx, fn)
}

func (n *NodeStorage) NewUoW() (node.UoW, error) {
	baseUoW, err := n.base.NewUoW()
	if err != nil {
		return nil, err
	}
	return &NodeUoW{nodeID: n.nodeID, base: baseUoW}, nil
}

type NodeUoW struct {
	nodeID models.NodeID
	base   UoW
}

var _ node.UoW = (*NodeUoW)(nil)

func (uow *NodeUoW) Do(ctx context.Context, fn node.UoWFn) error {
	return uow.base.Do(ctx, func(uowctx UoWContext) error {
		return fn(&NodeUoWContext{
			nodeID: uow.nodeID,
			base:   uowctx,
		})
	})
}

// node uow context implements node.UOWContext by itself:
// it implements all of components of node.UowContext
type NodeUoWContext struct {
	nodeID models.NodeID
	base   UoWContext
}

// UpdateClientTemplate implements node.NodeConfigStorage.

var _ node.UoWContext = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) NodeConfigStorage() node.NodeConfigStorage {
	return c
}

func (c *NodeUoWContext) NodeStatusStorage() node.NodeStatusStorage {
	return c
}

func (c *NodeUoWContext) PendingSyncsStorage() node.PendingSyncsStorage {
	return c
}

func (c *NodeUoWContext) UsersStorage() node.UsersStorage {
	return c
}

var _ node.NodeConfigStorage = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) UpdateClientConfig(ctx context.Context,
	tmpl *models.ClientConfig,
) error {
	return c.base.NodeConfigStorage().UpdateClientConfig(ctx, c.nodeID, tmpl)
}

var _ node.UsersStorage = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) ListUsers(ctx context.Context) (
	[]models.UserTargetState, error,
) {
	return c.base.UsersStorage().ListUsers(ctx)
}

var _ node.NodeStatusStorage = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) UpdateCurrentStatus(ctx context.Context,
	s models.NodeStatus,
) error {
	return c.base.NodeStatusStorage().UpdateCurrentStatus(ctx, c.nodeID, s)
}

func (c *NodeUoWContext) FetchNodeStatus(ctx context.Context) (
	target models.NodeStatus, current models.NodeStatus, err error,
) {
	return c.base.NodeStatusStorage().FetchNodeStatus(ctx, c.nodeID)
}

var _ node.PendingSyncsStorage = (*NodeUoWContext)(nil)

func (c *NodeUoWContext) FindPendingSyncs(ctx context.Context) (
	[]models.UserSyncStatus, error,
) {
	return c.base.PendingSyncsStorage().FindPendingSyncs(ctx, c.nodeID)
}

func (c *NodeUoWContext) PatchPendingSyncs(ctx context.Context,
	patch []models.UserStatusPatch,
) error {
	return c.base.PendingSyncsStorage().PatchPendingSyncs(ctx, c.nodeID, patch)
}
