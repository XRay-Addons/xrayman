package poolsyncer

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolSyncer struct {
	client PoolClient
	uow    UoW
	syncer NodesSyncer
}

func New(client PoolClient, uow UoW, syncer NodesSyncer) (*PoolSyncer, error) {
	if client == nil {
		return nil, fmt.Errorf("pool syncer init: pool client: %w", errdefs.ErrNilArgPassed)
	}
	if uow == nil {
		return nil, fmt.Errorf("pool syncer init: uow %w", errdefs.ErrNilArgPassed)
	}
	if syncer == nil {
		return nil, fmt.Errorf("pool syncer init: node syncer %w", errdefs.ErrNilArgPassed)
	}

	return &PoolSyncer{
		client: client,
		uow:    uow,
		syncer: syncer,
	}, nil
}

func (s *PoolSyncer) SyncPool(ctx context.Context) (*models.PoolSyncResult, error) {
	// get nodes list
	var nodes []models.Node
	if err := s.uow.Do(ctx, func(uow UoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("pool sync: %w", err)
	}

	// run nodes syncs
	nodesSyncResult := s.syncNodes(ctx, nodes)

	// determine sync status
	poolSyncStatus := s.getPoolSyncStatus(nodesSyncResult)

	return &models.PoolSyncResult{
		Status: poolSyncStatus,
		Nodes:  nodesSyncResult,
	}, nil
}

func (s *PoolSyncer) syncNodes(ctx context.Context,
	nodes []models.Node,
) []models.NodeSyncResult {
	nodeSyncResults := make([]models.NodeSyncResult, len(nodes))

	var wg sync.WaitGroup
	for idx, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeSyncResults[idx] = models.NodeSyncResult{
				ID:       node.ID,
				Endpoint: node.Config.ConnectionInfo.Endpoint,
				Err:      s.syncNode(ctx, node),
			}
		}()
	}
	wg.Wait()

	return nodeSyncResults
}

func (s *PoolSyncer) syncNode(ctx context.Context, node models.Node) error {
	nodeUoW := &NodeUoW{
		base:   s.uow,
		nodeID: node.ID,
	}

	nodeClient, err := s.client.GetNodeClient(ctx, node.Config.ConnectionInfo)
	if err != nil {
		return fmt.Errorf("pool sync node: %w", err)
	}

	if err := s.syncer.SyncNode(ctx, nodeUoW, nodeClient); err != nil {
		return fmt.Errorf("pool sync node: %w", err)
	}
	return nil
}

func (s *PoolSyncer) getPoolSyncStatus(
	nodes []models.NodeSyncResult,
) models.PoolSyncStatus {
	// determine pool sync status
	syncedCount := 0
	for _, node := range nodes {
		if node.Err == nil {
			syncedCount++
		}
	}

	switch {
	case syncedCount == len(nodes):
		return models.PoolSyncOk
	case syncedCount == 0:
		return models.PoolSyncFailed
	default:
		return models.PoolSyncPartially
	}
}
