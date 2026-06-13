package syncman

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/pooljob"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type SyncMan struct {
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
	defaultSyncInterval = 60 * time.Second
)

func New(syncer PoolSyncer, opts ...Option) (*SyncMan, error) {
	if syncer == nil {
		return nil, errdefs.NilArg("syncer")
	}
	cfg := options{
		interval: defaultSyncInterval,
		log:      zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	jobFn := func(ctx context.Context) (*models.PoolOpResult, error) {
		return syncer.SyncPoolState(ctx)
	}
	job, err := pooljob.New(jobFn, cfg.interval, "sync state", cfg.log)
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

func (m *SyncMan) Stop() error {
	if m == nil || m.job == nil {
		return nil
	}
	return m.job.Stop()
}
