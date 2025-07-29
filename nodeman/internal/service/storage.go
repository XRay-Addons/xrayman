package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeConfigStorage interface {
	AddNode(ctx context.Context,
		node *models.NodeConfig) error
	ListNodes(ctx context.Context) (
		[]models.NodeConfig, error)

	UpdateConnectionInfo(ctx context.Context, id models.NodeID,
		connInfo *models.NodeConnectionInfo) error
	GetConnectionInfo(ctx context.Context, id models.NodeID) (
		*models.NodeConnectionInfo, error)

	UpdateClientConfig(ctx context.Context, id models.NodeID,
		cfg *models.ClientConfig) error
	GetClientConfig(ctx context.Context, id models.NodeID) (
		*models.ClientConfig, error)
}

type NodeStatusStorage interface {
	FetchNodeStatus(ctx context.Context,
		id models.NodeID) (
		target, current models.NodeStatus, err error)
	UpdateTargetStatus(ctx context.Context,
		id models.NodeID, s models.NodeStatus) error
	UpdateCurrentStatus(ctx context.Context,
		id models.NodeID, s models.NodeStatus) error
}

type UsersStorage interface {
	AddUser(ctx context.Context,
		user *models.UserProfile) error
	ListUsers(ctx context.Context) (
		[]models.UserTargetState, error)
}

type UserStatusStorage interface {
	GetUserStatus(ctx context.Context,
		id models.UserID) (models.UserStatus, error)
	SetUserStatus(ctx context.Context,
		int models.UserID, s models.UserStatus) error
}

type PendingSyncsStorage interface {
	FindPendingSyncs(ctx context.Context, id models.NodeID) (
		[]models.UserSyncStatus, error)
	PatchPendingSyncs(ctx context.Context, id models.NodeID,
		patch []models.UserStatusPatch) error
}

type UoWContext interface {
	NodeConfigStorage() NodeConfigStorage
	NodeStatusStorage() NodeStatusStorage
	UsersStorage() UsersStorage
	UserStatusStorage() UserStatusStorage
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
