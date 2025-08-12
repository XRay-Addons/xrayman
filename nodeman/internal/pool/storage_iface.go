package pool

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type StatesStorage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	UpdateClientConfig(ctx context.Context, id models.NodeID,
		cfg models.ClientConfig) error
	FetchNodeStatus(ctx context.Context, id models.NodeID) (
		target, current models.NodeStatus, err error)
	UpdateCurrentStatus(ctx context.Context, id models.NodeID,
		s models.NodeStatus) error
}

type SyncsStorage interface {
	FindPendingSyncs(ctx context.Context, id models.NodeID) (
		[]models.UserSyncStatus, error)
	PatchPendingSyncs(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	UsersStorage
	StatesStorage
	SyncsStorage
}

type UoWFn = uow.Fn[UoWContext]

type UoW = uow.UoW[UoWContext]
