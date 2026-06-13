package statsman

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/pooljob"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type StatsMan struct {
	job *pooljob.PoolJob
}

type options struct {
	interval time.Duration
	log      *zap.Logger
}

type Option func(o *options)

func WithSyncInterval(interval time.Duration) Option {
	return func(o *options) {
		o.interval = interval
	}
}

func WithLogger(log *zap.Logger) Option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

const (
	defaultSyncInterval = 5 * time.Second
)

func New(updater StatsUpdater, opts ...Option) (*StatsMan, error) {
	if updater == nil {
		return nil, errdefs.NilArg("updater")
	}
	cfg := options{
		interval: defaultSyncInterval,
		log:      zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	jobFn := func(ctx context.Context) (*models.PoolOpResult, error) {
		return updater.UpdatePoolStats(ctx)
	}
	job, err := pooljob.New(jobFn, cfg.interval, "update stats", cfg.log)
	if err != nil {
		return nil, err
	}

	// init default options
	m := &StatsMan{
		job: job,
	}

	return m, nil
}

func (m *StatsMan) Run() error {
	if m == nil || m.job == nil {
		return errdefs.NilCall()
	}
	return m.job.Run()
}

func (m *StatsMan) Stop() error {
	if m == nil || m.job == nil {
		return nil
	}
	return m.job.Stop()
}
