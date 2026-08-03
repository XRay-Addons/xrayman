package nodestats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	UpdateNodeStats(ctx context.Context, s models.NodeStats) error
}
