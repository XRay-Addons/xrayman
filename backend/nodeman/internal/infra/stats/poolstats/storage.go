package poolstats

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type StatsStorage interface {
	ListNodes(ctx context.Context) (
		[]models.Node, error)
	UpdateNodeStats(ctx context.Context, id models.NodeID,
		stats models.NodeStats) error
	UpdateDailyStats(ctx context.Context,
		day time.Time) error
}

type UoWContext interface {
	StatsStorage
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
