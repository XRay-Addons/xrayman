package dbstorage

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// poolstats.UoWContext impl
func (uow *uowctx) UpdateNodeStats(ctx context.Context,
	nodeID models.NodeID, stats models.NodeStats,
) error {
	args := queries.UpdateTotalStatsParams{
		UserID:   make([]int64, 0, len(stats.Users)),
		Upload:   make([]int64, 0, len(stats.Users)),
		Download: make([]int64, 0, len(stats.Users)),
	}
	for _, u := range stats.Users {
		args.UserID = append(args.UserID, int64(u.ID))
		args.Upload = append(args.Upload, u.Upload)
		args.Download = append(args.Download, u.Download)
	}
	if err := uow.q.UpdateTotalStats(ctx, args); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (uow *uowctx) UpdateDailyStats(ctx context.Context,
	day time.Time,
) error {
	if err := uow.q.UpdateDailyStats(ctx, day); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
