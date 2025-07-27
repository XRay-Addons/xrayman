package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeConfigStorage interface {
	AddNode(ctx context.Context,
		node *models.NodeConfig) error
	RemoveNode(ctx context.Context,
		id models.NodeID) error
	ListNodes(ctx context.Context) (
		[]models.NodeConfig, error)

	UpdateConnectionInfo(ctx context.Context, id models.NodeID,
		connInfo *models.NodeConnectionInfo)
	GetConnectionInfo(ctx context.Context, id models.NodeID) (
		*models.NodeConnectionInfo, error)

	UpdateClientTemplate(ctx context.Context, id models.NodeID,
		tmpl *models.ClientTemplate) error
	GetClientTemplate(ctx context.Context,
		id models.NodeID) error
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
		user *models.UserTargetState) error
	RemoveUser(ctx context.Context,
		id models.UserID) error
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
