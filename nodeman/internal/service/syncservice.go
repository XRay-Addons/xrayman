package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type SyncService interface {
	SyncNodesPool(ctx context.Context) ([]models.NodeSyncResult, error)
}
