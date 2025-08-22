package node

import (
	"context"

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
		return errdefs.NewNilArg("uow")
	}
	if client == nil {
		return errdefs.NewNilArg("client")
	}
	nodeSyncer := syncer{
		uow:    uow,
		client: client,
	}
	if err := nodeSyncer.SyncNodeState(ctx); err != nil {
		return err
	}
	return nil
}
