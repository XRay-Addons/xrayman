package syncman

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/waveexec"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type Manager struct {
	executor *waveexec.WaveExecutor[[]models.NodeSyncResult]

	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	log *zap.Logger
}

type Option func(m *Manager)

func WithSyncInterval(interval time.Duration) Option {
	return func(m *Manager) {
		m.interval = interval
	}
}

func WithLog(log *zap.Logger) Option {
	return func(m *Manager) {
		if log != nil {
			m.log = log
		}
	}
}

func New(syncer PoolSyncer, options ...Option) (*Manager, error) {
	if syncer == nil {
		return nil, fmt.Errorf("pool monitor: init: %w", errdefs.ErrNilArgPassed)
	}
	// init default options
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		interval: 5 * time.Second,
		cancel:   cancel,
		log:      zap.NewNop(),
	}
	// apply options
	for _, o := range options {
		o(m)
	}
	// add sync loop
	syncFn := m.syncFn(syncer)
	m.executor = waveexec.NewWaveExecutor(syncFn)

	// run sync loop
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.syncLoop(ctx)
	}()

	return m, nil
}

func (m *Manager) Close() error {
	if m == nil {
		return nil
	}
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.wg.Wait()
	m.executor.Close()
	return nil
}

func (m *Manager) SyncNodesPool(ctx context.Context) ([]models.NodeSyncResult, error) {
	if m == nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", errdefs.ErrNilObjectCall)
	}
	syncResult, err := m.executor.Invoke(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", err)
	}
	return *syncResult, nil
}

func (m *Manager) syncFn(ps PoolSyncer) waveexec.Fn[[]models.NodeSyncResult] {
	return func(ctx context.Context) ([]models.NodeSyncResult, error) {
		syncResult, err := ps.SyncPoolState(ctx)
		if err == nil {
			m.logSyncResult(syncResult)
		} else {
			m.log.Error("pool sync", zap.Error(err))
		}

		return syncResult, err
	}
}

func (m *Manager) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// set sync time limit to m.interval
			syncCtx, cancel := context.WithTimeout(ctx, m.interval)
			m.SyncNodesPool(syncCtx)
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) logSyncResult(r []models.NodeSyncResult) {
	for _, n := range r {
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
