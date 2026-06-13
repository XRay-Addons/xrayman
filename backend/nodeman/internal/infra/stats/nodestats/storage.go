package nodestats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type StatsStorage interface {
	UpdateNodeStats(ctx context.Context, s models.NodeStats) error
}

type UoWContext interface {
	StatsStorage
}

type UoWFn = uow.Fn[UoWContext]

type Storage = uow.Storage[UoWContext]
