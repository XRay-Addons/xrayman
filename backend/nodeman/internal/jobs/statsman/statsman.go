package statsman

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/job"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type StatsMan struct {
	updateStatsJob *job.PoolJob
	updateDailyJob *job.Job
}

type options struct {
	log *zap.Logger
}

type Option func(o *options)

func WithLogger(log *zap.Logger) Option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

func New(updater StatsUpdater, interval time.Duration, opts ...Option) (*StatsMan, error) {
	if updater == nil {
		return nil, errdefs.NilArg("updater")
	}
	cfg := options{
		log: zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	updateJobFn := func(ctx context.Context) (*models.PoolOpResult, error) {
		return updater.UpdatePoolStats(ctx)
	}
	updateStatsJob, err := job.NewPoolJob(updateJobFn, interval, "update stats", cfg.log)
	if err != nil {
		return nil, err
	}

	// update daily stats every hour for:
	// - to be sure it runs at least once every day
	// - to not lost data for more than one hour in case of fail
	updateDailyFn := func(ctx context.Context) error {
		return updater.UpdateDailyStats(ctx)
	}
	updateDailyJob, err := job.NewJob(updateDailyFn, time.Hour, "update daily stats", cfg.log)
	if err != nil {
		return nil, err
	}

	// init default options
	m := &StatsMan{
		updateStatsJob: updateStatsJob,
		updateDailyJob: updateDailyJob,
	}

	return m, nil
}

func (m *StatsMan) Run() error {
	if m == nil || m.updateStatsJob == nil || m.updateDailyJob == nil {
		return errdefs.NilCall()
	}
	return xerr.Join(
		m.updateStatsJob.Run(),
		m.updateDailyJob.Run(),
	)
}

func (m *StatsMan) Stop() {
	if m == nil {
		return
	}
	if m.updateStatsJob != nil {
		m.updateStatsJob.Stop()
	}
	if m.updateDailyJob != nil {
		m.updateDailyJob.Stop()
	}
}
