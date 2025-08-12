package node

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
)

type Syncer struct {
}

func NewSyncer() *Syncer {
	return &Syncer{}
}

var _ pool.NodeSyncer = (*Syncer)(nil)

func (s *Syncer) SyncNodeState(ctx context.Context,
	client pool.NodeClient, uow pool.NodeUoW,
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
