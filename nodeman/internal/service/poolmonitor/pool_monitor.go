package poolmonitor

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

type PoolMonitor struct {
	executor *waveexec.WaveExecutor

	syncInterval time.Duration
	cancel       context.CancelFunc
	wg           sync.WaitGroup

	log *zap.Logger
}

type Option func(pm *PoolMonitor)

func WithSyncInterval(interval time.Duration) Option {
	return func(pm *PoolMonitor) {
		pm.syncInterval = interval
	}
}

func WithLog(log *zap.Logger) Option {
	return func(pm *PoolMonitor) {
		if log != nil {
			pm.log = log
		}
	}
}

func New(syncer PoolSyncer, options ...Option) (*PoolMonitor, error) {
	if syncer == nil {
		return nil, fmt.Errorf("pool monitor: init: %w", errdefs.ErrNilArgPassed)
	}
	// init
	ctx, cancel := context.WithCancel(context.Background())
	pm := &PoolMonitor{
		executor:     waveexec.NewWaveExecutor(syncFn(syncer)),
		syncInterval: 5 * time.Second,
		cancel:       cancel,
		log:          zap.NewNop(),
	}
	// apply options
	for _, o := range options {
		o(pm)
	}

	// run sync loop
	pm.wg.Add(1)
	go func() {
		defer pm.wg.Done()
		pm.syncLoop(ctx)
	}()

	return pm, nil
}

func (pm *PoolMonitor) Close() error {
	if pm == nil {
		return nil
	}
	if pm.cancel != nil {
		pm.cancel()
	}
	pm.wg.Wait()
	pm.executor.Close()
	return nil
}

func (pm *PoolMonitor) Sync(ctx context.Context) (*models.PoolSyncResult, error) {
	if pm == nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", errdefs.ErrNilObjectCall)
	}
	syncResult, err := pm.executor.Invoke(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", err)
	}
	res, ok := syncResult.(*models.PoolSyncResult)
	if !ok {
		return nil, fmt.Errorf("pool monitor: sync: cast result: %w", errdefs.ErrIPE)
	}
	return res, nil
}

func syncFn(ps PoolSyncer) waveexec.Fn {
	return func(ctx context.Context) (any, error) {
		return ps.SyncNodesPool(ctx)
	}
}

func (pm *PoolMonitor) syncLoop(ctx context.Context) {
	for {
		// set sync time limit to pm.syncInterval
		syncCtx, cancel := context.WithTimeout(ctx, pm.syncInterval)
		syncResult, err := pm.Sync(syncCtx)
		cancel()

		// log results
		if err != nil {
			pm.log.Error("pool sync", zap.Error(err))
		} else {
			pm.logSyncResult(*syncResult)
		}

		// wait 'syncInterval' time and sync again
		select {
		case <-time.After(pm.syncInterval):
		case <-ctx.Done():
			return
		}
	}
}

func (pm *PoolMonitor) logSyncResult(r models.PoolSyncResult) {
	for _, n := range r.Nodes {
		if n.Err == nil {
			pm.log.Info("background node sync OK",
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		} else {
			pm.log.Error("background node sync",
				zap.Error(n.Err),
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		}
	}
}
