package node

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// storage for node interface
type NodeStorage interface {
	FetchNodeStatus(ctx context.Context) (target, current models.NodeStatus, err error)
	ListManagedUsers(ctx context.Context) ([]models.UserTargetState, error)
	FindPendingSyncs(ctx context.Context) ([]models.UserSyncStatus, error)
	BeginTx() StorageTx
}

type StorageTx interface {
	UpdateNodeStatus(s models.NodeStatus)
	UpdateNodeConfig(c models.NodeConfig)
	UpdateNodeUsers(u []models.UserStatusPatch)
	Commit(ctx context.Context) error
}

// node api interface
type NodeAPI interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.NodeConfig, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUserStates(ctx context.Context, transitions []models.UserTargetState) error
}
