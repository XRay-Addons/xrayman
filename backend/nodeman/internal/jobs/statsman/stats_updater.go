package statsman

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type StatsUpdater interface {
	UpdatePoolStats(ctx context.Context) (*models.PoolOpResult, error)
	UpdateDailyStats(ctx context.Context) error
}
