package poolsyncer

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/nodesyncer"
)

type syncer struct {
	storage Storage
	client  Client
}

func (s *syncer) SyncPoolState(ctx context.Context) (*models.PoolSyncResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	nodes, err := s.listSyncingNodes(ctx)
	if err != nil {
		return nil, err
	}
	syncResult := s.syncNodes(ctx, nodes)
	return &syncResult, nil
}

type syncingNode struct {
	node    models.Node
	storage nodesyncer.Storage
	client  nodesyncer.Client
}

func (s *syncer) listSyncingNodes(ctx context.Context) ([]syncingNode, error) {
	var nodes []models.Node
	if err := s.storage.DoUoW(ctx, func(uow UoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, err
	}

	syncingNodes := make([]syncingNode, 0, len(nodes))
	for _, node := range nodes {
		nodeStorage := &nodeStorage{
			base:   s.storage,
			nodeID: node.ID,
		}
		nodeClient, err := s.client.GetNodeClient(ctx, node.Config.ConnectionInfo)
		if err != nil {
			return nil, err
		}
		syncingNodes = append(syncingNodes, syncingNode{
			node:    node,
			storage: nodeStorage,
			client:  nodeClient,
		})
	}

	return syncingNodes, nil
}

func (s *syncer) syncNodes(ctx context.Context, nodes []syncingNode) models.PoolSyncResult {
	nodeSyncResults := make([]models.NodeSyncResult, len(nodes))
	var wg sync.WaitGroup
	for idx, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeSyncResults[idx] = syncNode(ctx, node)
		}()
	}
	wg.Wait()

	return models.PoolSyncResult{
		Nodes: nodeSyncResults,
	}
}

func syncNode(ctx context.Context, node syncingNode) models.NodeSyncResult {
	syncErr := nodesyncer.SyncState(ctx, node.client, node.storage)
	if syncErr != nil {
		syncErr = errdefs.WrapWithf(syncErr, "nodeID: %v", node.node.ID)
	}
	return models.NodeSyncResult{
		ID:       node.node.ID,
		Endpoint: node.node.Config.ConnectionInfo.Endpoint,
		Err:      syncErr,
	}
}
