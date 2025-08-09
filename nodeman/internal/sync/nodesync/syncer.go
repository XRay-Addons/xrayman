package nodesync

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/sync/poolsync"
)

type Syncer struct {
}

func New() *Syncer {
	return &Syncer{}
}

var _ poolsync.NodeSyncer = (*Syncer)(nil)

func (s *Syncer) SyncNodeState(ctx context.Context,
	client poolsync.NodeClient, uow poolsync.NodeUoW,
) error {
	if uow == nil {
		return fmt.Errorf("syncer: sync node state: uow: %w", errdefs.ErrNilArgPassed)
	}
	if client == nil {
		return fmt.Errorf("syncer: sync node state: client: %w", errdefs.ErrNilArgPassed)
	}
	nodeSyncer := syncer{
		uow:    uow,
		client: client,
	}
	if err := nodeSyncer.SyncNodeState(ctx); err != nil {
		return fmt.Errorf("syncer: sync node state: %w", err)
	}
	return nil
}
