package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolMonitor interface {
	Sync(ctx context.Context) (*models.PoolSyncResult, error)
}
