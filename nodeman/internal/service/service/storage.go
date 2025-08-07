package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/poolsyncer"
)

type NodesStorage interface {
	// add new node to storage, assign NodeID to node
	NewNode(ctx context.Context, node *models.Node) error
	// change node target state
	SetTargetNodeStatus(ctx context.Context, id models.NodeID, state models.NodeStatus) error
}

type UoWContext interface {
	poolsyncer.UoWContext
	NodesStorage
}

type UoWFn = uow.Fn[UoWContext]

type UoW = uow.UoW[UoWContext]

type Storage = uow.Storage[UoWContext]
