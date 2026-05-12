package syncman

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type SyncMan struct {
	syncer PoolSyncer

	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	log *zap.Logger
}

type Option func(m *SyncMan)

func WithSyncInterval(interval time.Duration) Option {
	return func(m *SyncMan) {
		m.interval = interval
	}
}

func WithLogger(log *zap.Logger) Option {
	return func(m *SyncMan) {
		if log != nil {
			m.log = log
		}
	}
}

const (
	defaultSyncInterval = 5 * time.Second
)

func New(syncer PoolSyncer, options ...Option) (*SyncMan, error) {
	if syncer == nil {
		return nil, errdefs.NilArg("syncer")
	}
	// init default options
	m := &SyncMan{
		syncer:   syncer,
		interval: defaultSyncInterval,
		log:      zap.NewNop(),
	}
	// apply options
	for _, o := range options {
		o(m)
	}

	return m, nil
}

func (m *SyncMan) Run() error {
	if m == nil {
		return errdefs.NilCall()
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	// run sync loop
	m.wg.Add(1)
	defer m.wg.Done()
	m.syncLoop(ctx)

	return ctx.Err() //nolint:wrapcheck
}

func (m *SyncMan) Stop() error {
	if m == nil {
		return nil
	}
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.wg.Wait()
	return nil
}

func (m *SyncMan) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		// set sync time limit to m.interval
		syncCtx, cancel := context.WithTimeout(ctx, m.interval)
		syncRes, err := m.syncer.SyncPoolState(syncCtx)
		cancel()
		m.logSyncResult(syncRes, err)

		select {
		case <-time.After(m.interval):
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (m *SyncMan) logSyncResult(r *models.PoolSyncResult, err error) {
	if err != nil {
		m.log.Error("background node sync", zap.Error(err))
		return
	}
	for _, n := range r.Nodes {
		if n.Err == nil {
			m.log.Info("background node sync OK",
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		} else {
			m.log.Error("background node sync",
				zap.Error(n.Err),
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		}
	}
}
