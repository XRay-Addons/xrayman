package app

import (
	"context"
	"fmt"
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// sync state job
type SyncStateResult models.PoolOpResult

var _ zapcore.ObjectMarshaler = (*SyncStateResult)(nil)

func (r SyncStateResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, node := range r.Nodes {
		key := fmt.Sprintf("sync %s state", node.Endpoint)
		if node.Err != nil {
			enc.AddString(key, fmt.Sprintf("%+v", node.Err))
		} else {
			enc.AddString(key, "OK")
		}
	}
	return nil
}

func syncState(s *poolsync.Service) func(context.Context) (*SyncStateResult, error) {
	return func(ctx context.Context) (*SyncStateResult, error) {
		r, err := s.SyncPoolState(ctx)
		if err != nil {
			return nil, err
		}
		rw := SyncStateResult(*r)
		return &rw, nil
	}
}

var syncStateJob = gx.ProvideAnnotated(
	func(s *poolsync.Service, cfg *config.Config, l *zap.Logger) gx.Job {
		return gx.MakeTickJob("sync state", syncState(s), cfg.StateSyncInterval, l)
	},
	gx.ResultTags(`group:"tick-jobs"`),
)

// update stats job
type UpdateStatsResult models.PoolOpResult

var _ zapcore.ObjectMarshaler = (*UpdateStatsResult)(nil)

func (r UpdateStatsResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, node := range r.Nodes {
		key := fmt.Sprintf("update %s stats", node.Endpoint)
		if node.Err != nil {
			enc.AddString(key, fmt.Sprintf("%+v", node.Err))
		} else {
			enc.AddString(key, "OK")
		}
	}
	return nil
}

func updateStats(s *poolstats.Service) func(context.Context) (*UpdateStatsResult, error) {
	return func(ctx context.Context) (*UpdateStatsResult, error) {
		r, err := s.UpdatePoolStats(ctx)
		if err != nil {
			return nil, err
		}
		rw := UpdateStatsResult(*r)
		return &rw, nil
	}
}

var updateStatsJob = gx.ProvideAnnotated(
	func(s *poolstats.Service, cfg *config.Config, l *zap.Logger) gx.Job {
		return gx.MakeTickJob("update stats", updateStats(s), cfg.StatsSyncInterval, l)
	},
	gx.ResultTags(`group:"tick-jobs"`),
)

// refresh daily stats job. run every hour for:
// - to be sure it runs at least once every day
// - to not lost data for more than one hour in case of fail
type RefreshDailyStatsResult struct{}

var _ zapcore.ObjectMarshaler = (*RefreshDailyStatsResult)(nil)

func (r RefreshDailyStatsResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("refresh daily stats", "OK")
	return nil
}

func refreshDailyStats(s *poolstats.Service) func(context.Context) (*RefreshDailyStatsResult, error) {
	return func(ctx context.Context) (*RefreshDailyStatsResult, error) {
		if err := s.RefreshDailyStats(ctx); err != nil {
			return nil, err
		}
		return &RefreshDailyStatsResult{}, nil
	}
}

var refreshDailyStatsJob = gx.ProvideAnnotated(
	func(s *poolstats.Service, cfg *config.Config, l *zap.Logger) gx.Job {
		return gx.MakeTickJob("refresh daily stats", refreshDailyStats(s), time.Hour, l)
	},
	gx.ResultTags(`group:"tick-jobs"`),
)

// all tick jobs together
type TickJobsParams struct {
	gx.In
	Lc   gx.Lifecycle
	Jobs []gx.Job `group:"tick-jobs"`
}

var tickJobs = gx.Invoke(
	func(p TickJobsParams) {
		for _, j := range p.Jobs {
			p.Lc.AppendJob(j)
		}
	},
)

var TickJobs = gx.Module("tick-jobs",
	syncStateJob,
	updateStatsJob,
	refreshDailyStatsJob,
	tickJobs,
)
