package syncer

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolSyncer struct {
	uow    PoolUoW
	client PoolClient
}

func New(uow PoolUoW, client PoolClient) (*PoolSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("pool syncer init: pool client: %w", errdefs.ErrNilArgPassed)
	}
	if uow == nil {
		return nil, fmt.Errorf("pool syncer init: uow %w", errdefs.ErrNilArgPassed)
	}

	return &PoolSyncer{
		uow:    uow,
		client: client,
	}, nil
}

func (s *PoolSyncer) SyncPoolState(ctx context.Context) ([]models.NodeSyncResult, error) {
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

func (s *PoolSyncer) listNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	if err := s.uow.Do(ctx, func(uow PoolUoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return nodes, nil
}

func (s *PoolSyncer) syncNodes(ctx context.Context, nodes []models.Node) []models.NodeSyncResult {
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

func (s *PoolSyncer) syncNode(ctx context.Context, node models.Node) (res models.NodeSyncResult) {
	res = models.NodeSyncResult{
		ID:       node.ID,
		Endpoint: node.Config.ConnectionInfo.Endpoint,
	}

	// get node storage
	nodeUoW := &PoolNodeUoW{
		base:   s.uow,
		nodeID: node.ID,
	}
	// get node client
	nodeClient, err := s.client.GetNodeClient(ctx, node.Config.ConnectionInfo)
	if err != nil {
		res.Err = fmt.Errorf("pool sync node: %w", err)
		return
	}
	// get node syncer
	syncer, err := NewNodeSyncer(nodeUoW, nodeClient)
	if err != nil {
		res.Err = fmt.Errorf("pool sync node: %w", err)
		return
	}
	// sync
	res.Err = syncer.SyncNodeState(ctx)
	return
}
