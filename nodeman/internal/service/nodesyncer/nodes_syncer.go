package nodesyncer

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

type NodesSyncer struct {
}

func New() *NodesSyncer {
	return &NodesSyncer{}
}

// NodeSyncer implements poolsyncer.NodeSyncer without import it to
// avoid cyclic dependency
func (s *NodesSyncer) SyncNode(ctx context.Context,
	client NodeClient, uow UoW,
) error {
	if s == nil {
		return fmt.Errorf("node syncer: sync node: %w", errdefs.ErrNilObjectCall)
	}
	if client == nil {
		return fmt.Errorf("node syncer: sync node: client: %w", errdefs.ErrNilArgPassed)
	}
	if uow == nil {
		return fmt.Errorf("node syncer: sync node: uow: %w", errdefs.ErrNilArgPassed)
	}

	nodeSyncer, err := NewNodeSyncer(uow, client)
	if err != nil {
		return fmt.Errorf("nodes syncer: sync node: %w", err)
	}
	if err := nodeSyncer.SyncState(ctx); err != nil {
		return fmt.Errorf("node syncer: sync node: %w", err)
	}
	return nil
}
