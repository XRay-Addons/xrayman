package syncman

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/job"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type SyncMan struct {
	job *job.PoolJob
}

type options struct {
	interval time.Duration
	log      *zap.Logger
}

type Option func(o *options)

func WithLogger(log *zap.Logger) Option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

func New(syncer PoolSyncer, interval time.Duration, opts ...Option) (*SyncMan, error) {
	if syncer == nil {
		return nil, errdefs.NilArg("syncer")
	}
	if interval == 0 {
		return nil, errdefs.NilArg("interval")
	}
	cfg := options{
		interval: interval,
		log:      zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	jobFn := func(ctx context.Context) (*models.PoolOpResult, error) {
		return syncer.SyncPoolState(ctx)
	}
	job, err := job.NewPoolJob(jobFn, cfg.interval, "sync state", cfg.log)
	if err != nil {
		return nil, err
	}

	// init default options
	m := &SyncMan{
		job: job,
	}

	return m, nil
}

func (m *SyncMan) Run() error {
	if m == nil || m.job == nil {
		return errdefs.NilCall()
	}
	return m.job.Run()
}

func (m *SyncMan) Stop() {
	if m == nil || m.job == nil {
		return
	}
	m.job.Stop()
}
