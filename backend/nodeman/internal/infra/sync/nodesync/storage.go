package nodesync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type StateStorage interface {
	GetNodeStatus(ctx context.Context) (
		target, current models.NodeStatus, err error)
	SetClientConfig(ctx context.Context,
		cfg models.ClientConfigTemplate) error
	SetCurrentNodeStatus(ctx context.Context,
		s models.NodeStatus) error
}

type SyncsStorage interface {
	FindPendingSyncs(ctx context.Context) (
		[]models.UserSyncStatus, error)
	UpdateNodeUsers(ctx context.Context,
		patch []models.UserStatusPatch) error
	SetNodeUsers(ctx context.Context,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	UsersStorage
	StateStorage
	SyncsStorage
}

type UoWFn = uow.Fn[UoWContext]

type Storage = uow.Storage[UoWContext]
