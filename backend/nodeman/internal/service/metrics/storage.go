package metrics

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	GetNodeMetrics(ctx context.Context) ([]models.NodeMetrics, error)
}
