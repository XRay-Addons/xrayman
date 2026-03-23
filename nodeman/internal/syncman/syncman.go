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

func WithLog(log *zap.Logger) Option {
	return func(m *SyncMan) {
		if log != nil {
			m.log = log
		}
	}
}

const (
	defaultSyncInterval = 60 * time.Second
)

func New(syncer PoolSyncer, options ...Option) (*SyncMan, error) {
	if syncer == nil {
		return nil, errdefs.NewNilArg("syncer")
	}
	// init default options
	ctx, cancel := context.WithCancel(context.Background())
	m := &SyncMan{
		syncer:   syncer,
		interval: defaultSyncInterval,
		cancel:   cancel,
		log:      zap.NewNop(),
	}
	// apply options
	for _, o := range options {
		o(m)
	}

	// run sync loop
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.syncLoop(ctx)
	}()

	return m, nil
}

func (m *SyncMan) Close() error {
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
		select {
		case <-ticker.C:
			// set sync time limit to m.interval
			syncCtx, cancel := context.WithTimeout(ctx, m.interval)
			syncRes, err := m.syncer.SyncPoolState(syncCtx)
			cancel()

			m.logSyncResult(syncRes, err)
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
