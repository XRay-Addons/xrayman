package poolmonitor

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=pool_syncer.go -destination=./mocks/mock_pool_syncer.go -package=mocks PoolSyncer
type PoolSyncer interface {
	SyncNodesPool(ctx context.Context) (*models.PoolSyncResult, error)
}
