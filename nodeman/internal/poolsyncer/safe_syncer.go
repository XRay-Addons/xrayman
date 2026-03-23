package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/waveexec"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Syncer interface {
	SyncPoolState(ctx context.Context) (*models.PoolSyncResult, error)
	Close()
}

type safeSyncer struct {
	exec *waveexec.WaveExecutor[models.PoolSyncResult]
}

var _ Syncer = (*safeSyncer)(nil)

func New(client Client, storage Storage) (Syncer, error) {
	if client == nil {
		return nil, errdefs.NewNilArg("client")
	}
	if storage == nil {
		return nil, errdefs.NewNilArg("storage")
	}
	unsafeSyncer := syncer{
		storage: storage,
		client:  client,
	}
	syncFn := func(ctx context.Context) (*models.PoolSyncResult, error) {
		return unsafeSyncer.SyncPoolState(ctx)
	}
	return &safeSyncer{
		exec: waveexec.New(syncFn),
	}, nil
}

func (s *safeSyncer) Close() {
	if s == nil || s.exec == nil {
		return
	}
	s.exec.Close()
}

func (s *safeSyncer) SyncPoolState(ctx context.Context) (*models.PoolSyncResult, error) {
	return s.exec.Invoke(ctx)
}
