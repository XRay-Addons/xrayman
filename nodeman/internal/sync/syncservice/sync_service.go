package syncservice

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/waveexec"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"go.uber.org/zap"
)

type SyncService struct {
	executor *waveexec.WaveExecutor

	syncInterval time.Duration
	cancel       context.CancelFunc
	wg           sync.WaitGroup

	log *zap.Logger
}

var _ service.SyncService = (*SyncService)(nil)

type Option func(ss *SyncService)

func WithSyncInterval(interval time.Duration) Option {
	return func(ss *SyncService) {
		ss.syncInterval = interval
	}
}

func WithLog(log *zap.Logger) Option {
	return func(ss *SyncService) {
		if log != nil {
			ss.log = log
		}
	}
}

func New(syncer PoolSyncer, options ...Option) (*SyncService, error) {
	if syncer == nil {
		return nil, fmt.Errorf("pool monitor: init: %w", errdefs.ErrNilArgPassed)
	}
	// init default options
	ctx, cancel := context.WithCancel(context.Background())
	ss := &SyncService{
		syncInterval: 5 * time.Second,
		cancel:       cancel,
		log:          zap.NewNop(),
	}
	// apply options
	for _, o := range options {
		o(ss)
	}
	// add sync loop
	syncFn := ss.syncFn(syncer)
	ss.executor = waveexec.NewWaveExecutor(syncFn)

	// run sync loop
	ss.wg.Add(1)
	go func() {
		defer ss.wg.Done()
		ss.syncLoop(ctx)
	}()

	return ss, nil
}

func (ss *SyncService) Close() error {
	if ss == nil {
		return nil
	}
	if ss.cancel != nil {
		ss.cancel()
	}
	ss.wg.Wait()
	ss.executor.Close()
	return nil
}

func (ss *SyncService) SyncNodesPool(ctx context.Context) ([]models.NodeSyncResult, error) {
	if ss == nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", errdefs.ErrNilObjectCall)
	}
	syncResult, err := ss.executor.Invoke(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool monitor: sync: %w", err)
	}
	res, ok := syncResult.([]models.NodeSyncResult)
	if !ok {
		return nil, fmt.Errorf("pool monitor: sync: cast result: %w", errdefs.ErrIPE)
	}
	return res, nil
}

func (ss *SyncService) syncFn(ps PoolSyncer) waveexec.Fn {
	return func(ctx context.Context) (any, error) {
		syncResult, err := ps.SyncPoolState(ctx)
		if err != nil {
			ss.logSyncResult(syncResult)
		}

		return syncResult, err
	}
}

func (ss *SyncService) syncLoop(ctx context.Context) {
	for {
		// set sync time limit to ss.syncInterval
		syncCtx, cancel := context.WithTimeout(ctx, ss.syncInterval)
		syncResult, err := ss.SyncNodesPool(syncCtx)
		cancel()

		// log results
		if err != nil {
			ss.log.Error("pool sync", zap.Error(err))
		} else {
			ss.logSyncResult(syncResult)
		}

		// wait 'syncInterval' time and sync again
		select {
		case <-time.After(ss.syncInterval):
		case <-ctx.Done():
			return
		}
	}
}

func (ss *SyncService) logSyncResult(r []models.NodeSyncResult) {
	for _, n := range r {
		if n.Err == nil {
			ss.log.Info("background node sync OK",
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		} else {
			ss.log.Error("background node sync",
				zap.Error(n.Err),
				zap.String("nodeID", strconv.Itoa(int(n.ID))),
				zap.String("endpoint", n.Endpoint))
		}
	}
}
