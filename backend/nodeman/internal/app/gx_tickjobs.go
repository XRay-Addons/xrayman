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

// job result adapters for logging
type voidResultAdapter struct {
	name string
}

var _ zapcore.ObjectMarshaler = (*voidResultAdapter)(nil)

func (r voidResultAdapter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(r.name, "OK")
	return nil
}

type poolResultAdapter struct {
	models.PoolOpResult
	name string
}

var _ zapcore.ObjectMarshaler = (*poolResultAdapter)(nil)

func (r poolResultAdapter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, node := range r.Nodes {
		key := fmt.Sprintf("%s for node %s", r.name, node.Endpoint)
		if node.Err != nil {
			enc.AddString(key, fmt.Sprintf("%+v", node.Err))
		} else {
			enc.AddString(key, "OK")
		}
	}
	return nil
}

// sync state job
func syncState(s *poolsync.Service) func(context.Context) (*poolResultAdapter, error) {
	return func(ctx context.Context) (*poolResultAdapter, error) {
		r, err := s.SyncPoolState(ctx)
		if err != nil {
			return nil, err
		}
		return &poolResultAdapter{
			PoolOpResult: *r,
			name:         "sync state",
		}, nil
	}
}

var syncStateJob = gx.ProvideAnnotated(
	func(s *poolsync.Service, cfg *config.Config, l *zap.Logger) gx.Job {
		return gx.MakeTickJob("sync state", syncState(s), cfg.StateSyncInterval, l)
	},
	gx.ResultTags(`group:"tick-jobs"`),
)

// update stats job
func updateStats(s *poolstats.Service) func(context.Context) (*poolResultAdapter, error) {
	return func(ctx context.Context) (*poolResultAdapter, error) {
		r, err := s.UpdatePoolStats(ctx)
		if err != nil {
			return nil, err
		}
		return &poolResultAdapter{
			PoolOpResult: *r,
			name:         "sync state",
		}, nil
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
func refreshDailyStats(s *poolstats.Service) func(context.Context) (*voidResultAdapter, error) {
	return func(ctx context.Context) (*voidResultAdapter, error) {
		if err := s.RefreshDailyStats(ctx); err != nil {
			return nil, err
		}
		return &voidResultAdapter{
			name: "refresh daily stats",
		}, nil
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
