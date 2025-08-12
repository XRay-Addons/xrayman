package pool

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type NodeStateStorage interface {
	UpdateClientConfig(ctx context.Context,
		cfg models.ClientConfig) error
	FetchNodeStatus(ctx context.Context) (
		target, current models.NodeStatus, err error)
	UpdateCurrentStatus(ctx context.Context,
		s models.NodeStatus) error
}

type NodeSyncsStorage interface {
	FindPendingSyncs(ctx context.Context) (
		[]models.UserSyncStatus, error)
	PatchPendingSyncs(ctx context.Context,
		patch []models.UserStatusPatch) error
}

type NodeUoWContext interface {
	UsersStorage
	NodeStateStorage
	NodeSyncsStorage
}

type NodeUoWFn = uow.Fn[NodeUoWContext]
type NodeUoW = uow.UoW[NodeUoWContext]
