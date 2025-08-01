package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

////////////////////////////////////////////////////////////////////////////////
// storage to use in PoolSyncer

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type NodeStatesStorage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	UpdateClientConfig(ctx context.Context, id models.NodeID,
		cfg models.ClientConfig) error
	FetchNodeStatus(ctx context.Context, id models.NodeID) (
		target, current models.NodeStatus, err error)
	UpdateCurrentStatus(ctx context.Context, id models.NodeID,
		s models.NodeStatus) error
}

type UserSyncsStorage interface {
	FindPendingSyncs(ctx context.Context, id models.NodeID) (
		[]models.UserSyncStatus, error)
	PatchPendingSyncs(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	UsersStorage
	NodeStatesStorage
	UserSyncsStorage
}

type UoWFn func(UoWContext) error

type UoW interface {
	Do(ctx context.Context, fn UoWFn) error
}

type PoolStorage interface {
	NewUoW() (UoW, error)
	DoUoW(ctx context.Context, fn UoWFn) error
}
