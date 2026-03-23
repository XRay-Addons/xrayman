package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type StatesStorage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	GetNode(ctx context.Context, id models.NodeID) (
		*models.Node, bool, error)
	SetClientConfig(ctx context.Context, id models.NodeID,
		cfg models.ClientConfigTemplate) error
	SetCurrentNodeStatus(ctx context.Context, id models.NodeID,
		s models.NodeStatus) error
}

type SyncsStorage interface {
	FindPendingSyncs(ctx context.Context, id models.NodeID) (
		[]models.UserSyncStatus, error)
	UpdateNodeUsers(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
	SetNodeUsers(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	UsersStorage
	StatesStorage
	SyncsStorage
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
