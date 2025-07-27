package node

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeConfigStorage interface {
	UpdateClientTemplate(ctx context.Context,
		tmpl *models.ClientTemplate) error
}

type NodeStatusStorage interface {
	FetchNodeStatus(ctx context.Context) (
		target, current models.NodeStatus, err error)
	UpdateCurrentStatus(ctx context.Context,
		s models.NodeStatus) error
}

type UsersStorage interface {
	ListUsers(ctx context.Context) (
		[]models.UserTargetState, error)
}

type PendingSyncsStorage interface {
	FindPendingSyncs(ctx context.Context) (
		[]models.UserSyncStatus, error)
	PatchPendingSyncs(ctx context.Context,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	NodeConfigStorage() NodeConfigStorage
	NodeStatusStorage() NodeStatusStorage
	UsersStorage() UsersStorage
	PendingSyncsStorage() PendingSyncsStorage
}

type UoWFn func(UoWContext) error

type UoW interface {
	Do(ctx context.Context, fn UoWFn) error
}

type Storage interface {
	NewUoW() (UoW, error)
	DoUoW(ctx context.Context, fn UoWFn) error
}
