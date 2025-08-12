package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/syncman"
)

type syncFn = func(ctx context.Context) ([]models.NodeSyncResult, error)

type Syncer struct {
	fn syncFn
}

var _ syncman.PoolSyncer = (*Syncer)(nil)

func NewSyncer(uow UoW, client Client, syncer NodeSyncer) (*Syncer, error) {
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
		fn: getSyncFn(uow, client, syncer),
	}, nil
}

func (s *Syncer) SyncPoolState(ctx context.Context) ([]models.NodeSyncResult, error) {
	if s == nil {
		return nil, fmt.Errorf("syncer: sync pool state: %w", errdefs.ErrNilObjectCall)
	}
	return s.fn(ctx)
}

func getSyncFn(uow UoW, client Client, syncer NodeSyncer) syncFn {
	return func(ctx context.Context) ([]models.NodeSyncResult, error) {
		nodes, err := listSyncingNodes(ctx, uow, client)
		if err != nil {
			return nil, fmt.Errorf("sync nodes: %w", err)
		}
		return syncNodes(ctx, syncer, nodes), nil
	}
}

type syncingNode struct {
	node   models.Node
	uow    NodeUoW
	client NodeClient
}

func listSyncingNodes(ctx context.Context, uow UoW, client Client) ([]syncingNode, error) {
	var nodes []models.Node
	if err := uow.Do(ctx, func(uow UoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	syncingNodes := make([]syncingNode, 0, len(nodes))
	for _, node := range nodes {
		nodeUoW := &poolNodeUoW{
			base:   uow,
			nodeID: node.ID,
		}
		nodeClient, err := client.GetNodeClient(ctx, node.Config.ConnectionInfo)
		if err != nil {
			return nil, fmt.Errorf("list syncing nodes: %w", err)
		}
		syncingNodes = append(syncingNodes, syncingNode{
			node:   node,
			uow:    nodeUoW,
			client: nodeClient,
		})
	}

	return syncingNodes, nil
}

func syncNodes(ctx context.Context, syncer NodeSyncer, nodes []syncingNode) []models.NodeSyncResult {
	nodeSyncResults := make([]models.NodeSyncResult, len(nodes))
	var wg sync.WaitGroup
	for idx, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeSyncResults[idx] = syncNode(ctx, syncer, node)
		}()
	}
	wg.Wait()

	return nodeSyncResults
}

func syncNode(ctx context.Context, syncer NodeSyncer, node syncingNode) (res models.NodeSyncResult) {
	res = models.NodeSyncResult{
		ID:       node.node.ID,
		Endpoint: node.node.Config.ConnectionInfo.Endpoint,
	}

	if err := syncer.SyncNodeState(ctx, node.client, node.uow); err != nil {
		res.Err = fmt.Errorf("pool sync node: %w", err)
		return
	}

	return
}
