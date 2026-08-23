package nodesync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type TxFn = func(context.Context) error

type UsersStorage interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type StateStorage interface {
	GetNodeStatus(ctx context.Context) (
		target, current models.NodeStatus, err error)
	SetNodeSettings(ctx context.Context,
		cfg *models.NodeSettings) error
	SetCurrentNodeStatus(ctx context.Context,
		s models.NodeStatus) error
}

type SyncsStorage interface {
	FindPendingSyncs(ctx context.Context, rev models.Revision) (
		[]models.UserSyncStatus, error)
	GetNodeRev(ctx context.Context) (
		models.Revision, error)
	SetNodeRev(ctx context.Context,
		rev models.Revision) error
	UpdateNodeUsers(ctx context.Context,
		patch []models.UserStatusPatch) error
	SetNodeUsers(ctx context.Context,
		patch []models.UserStatusPatch) error
}

type Storage interface {
	UsersStorage
	StateStorage
	SyncsStorage
	// call multiple operations as tx
	DoTx(ctx context.Context, fn TxFn) error
}
