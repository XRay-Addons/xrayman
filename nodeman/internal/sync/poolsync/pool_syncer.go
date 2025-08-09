package poolsync

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Syncer struct {
	uow    UoW
	client Client
	syncer NodeSyncer
}

func New(uow UoW, client Client, syncer NodeSyncer) (*Syncer, error) {
	if client == nil {
		return nil, fmt.Errorf("pool syncer init: pool client: %w", errdefs.ErrNilArgPassed)
	}
	if uow == nil {
		return nil, fmt.Errorf("pool syncer init: uow %w", errdefs.ErrNilArgPassed)
	}
	if syncer == nil {
		return nil, fmt.Errorf("pool syncer init: node syncer %w", errdefs.ErrNilArgPassed)
	}

	return &Syncer{
		uow:    uow,
		client: client,
		syncer: syncer,
	}, nil
}

func (s *Syncer) SyncPoolState(ctx context.Context) ([]models.NodeSyncResult, error) {
	if s == nil || s.uow == nil || s.client == nil {
		return nil, fmt.Errorf("pool syncer: sync pool state: %w", errdefs.ErrNilObjectCall)
	}

	nodes, err := s.listNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool syncer: sync pool state: %w", err)
	}

	nodesSyncResult := s.syncNodes(ctx, nodes)

	return nodesSyncResult, nil
}

func (s *Syncer) listNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	if err := s.uow.Do(ctx, func(uow UoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return nodes, nil
}

func (s *Syncer) syncNodes(ctx context.Context, nodes []models.Node) []models.NodeSyncResult {
	nodeSyncResults := make([]models.NodeSyncResult, len(nodes))
	var wg sync.WaitGroup
	for idx, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeSyncResults[idx] = s.syncNode(ctx, node)
		}()
	}
	wg.Wait()

	return nodeSyncResults
}

func (s *Syncer) syncNode(ctx context.Context, node models.Node) (res models.NodeSyncResult) {
	res = models.NodeSyncResult{
		ID:       node.ID,
		Endpoint: node.Config.ConnectionInfo.Endpoint,
	}

	// get node storage
	nodeUoW := &poolNodeUoW{
		base:   s.uow,
		nodeID: node.ID,
	}
	// get node client
	nodeClient, err := s.client.GetNodeClient(ctx, node.Config.ConnectionInfo)
	if err != nil {
		res.Err = fmt.Errorf("pool sync node: %w", err)
		return
	}
	// sync node
	if err = s.syncer.SyncNodeState(ctx, nodeClient, nodeUoW); err != nil {
		res.Err = fmt.Errorf("pool sync node: %w", err)
		return
	}

	return
}
