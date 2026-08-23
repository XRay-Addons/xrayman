package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type TxFn = func(context.Context) error

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type StatesStorage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	GetNode(ctx context.Context, id models.NodeID) (
		*models.Node, error)
	SetNodeSettings(ctx context.Context, id models.NodeID,
		s *models.NodeSettings) error
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
	SetNodeRev(ctx context.Context,
		id models.NodeID, rev models.Revision) error
}

type Storage interface {
	UsersStorage
	StatesStorage
	SyncsStorage
	// call multiple operations as tx
	DoTx(ctx context.Context, fn TxFn) error
}

var _ poolop.Storage = (Storage)(nil)
