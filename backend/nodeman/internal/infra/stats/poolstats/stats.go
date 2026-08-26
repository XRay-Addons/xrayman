package poolstats

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type Stats struct {
	storage Storage
	op      *poolop.PoolOp
}

var _ statsman.StatsUpdater = (*Stats)(nil)

func New(client Client, storage Storage, log *zap.Logger) (*Stats, error) {
	if client == nil {
		return nil, errdefs.NilArg("client")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	op, err := poolop.New(
		storage,
		&nodeOp{storage: storage, client: client},
		log,
	)
	if err != nil {
		return nil, err
	}
	return &Stats{
		op:      op,
		storage: storage,
	}, nil
}

func (s *Stats) Close() {
	if s == nil || s.op == nil {
		return
	}
	s.op.Close()
}

func (s *Stats) UpdatePoolStats(ctx context.Context) (*models.PoolOpResult, error) {
	return s.op.ExecAll(ctx)
}

func (s *Stats) RefreshDailyStats(ctx context.Context) error {
	return s.storage.RefreshDailyStats(ctx, time.Now())
}
