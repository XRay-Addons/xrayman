package poolstats

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	GetNode(ctx context.Context, id models.NodeID) (
		*models.Node, error)
	UpdateStats(ctx context.Context, id models.NodeID,
		stats models.NodeStats) error
	RefreshDailyStats(ctx context.Context,
		day time.Time) error
}

var _ poolop.Storage = (Storage)(nil)
