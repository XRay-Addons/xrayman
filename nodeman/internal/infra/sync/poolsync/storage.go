package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/uow"
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
	DeleteNode(ctx context.Context,
		id models.NodeID) error
}

type SyncsStorage interface {
	FindPendingSyncs(ctx context.Context, id models.NodeID) (
		[]models.UserSyncStatus, error)
	UpdateNodeUsers(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
	SetNodeUsers(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
	DeleteUser(ctx context.Context,
		id models.UserID) error
}

type UoWContext interface {
	UsersStorage
	StatesStorage
	SyncsStorage
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
