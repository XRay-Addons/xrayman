package nodestats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	GetStats(ctx context.Context) (*models.NodeStats, error)
}
