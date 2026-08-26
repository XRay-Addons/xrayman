package nodestats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	UpdateStats(ctx context.Context, s models.NodeStats) error
}
