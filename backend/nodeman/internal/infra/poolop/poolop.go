package poolop

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/common/safego"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/waveexec"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type empty = struct{}
type nodeExec = *waveexec.WaveExecutor[empty]

type execItem struct {
	node models.Node
	err  error
}

type PoolOp struct {
	storage Storage
	nodeOp  NodeOp
	log     *zap.Logger

	nodeExecs map[models.NodeID]nodeExec
	mu        sync.RWMutex
}

func New(s Storage, op NodeOp, log *zap.Logger) (*PoolOp, error) {
	if s == nil {
		return nil, xerr.NilArg("s")
	}
	if op == nil {
		return nil, xerr.NilArg("op")
	}
	if log == nil {
		return nil, xerr.NilArg("log")
	}
	return &PoolOp{
		storage: s,
		nodeOp:  op,
		log:     log,

		nodeExecs: make(map[models.NodeID]nodeExec),
	}, nil
}

func (o *PoolOp) Close() {
	if o == nil {
		return
	}
	for _, exec := range o.nodeExecs {
		if exec != nil {
			exec.Close()
		}
	}
}

func (o *PoolOp) ExecAll(ctx context.Context) (*models.PoolOpResult, error) {
	nodes, err := o.storage.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	execItems := make([]execItem, 0, len(nodes))
	for _, node := range nodes {
		execItems = append(execItems, execItem{node: node})
	}

	o.exec(ctx, execItems)

	nodeResults := make([]models.NodeOpResult, len(execItems))
	for i, exec := range execItems {
		nodeResults[i].ID = exec.node.ID
		nodeResults[i].Endpoint = exec.node.Config.ConnectionInfo.Endpoint
		nodeResults[i].Err = exec.err
	}

	return &models.PoolOpResult{
		Nodes: nodeResults,
	}, nil
}

func (o *PoolOp) ExecNode(ctx context.Context, id models.NodeID) error {
	node, err := o.storage.GetNode(ctx, id)
	if err != nil {
		return err
	}
	execItems := []execItem{
		execItem{node: *node},
	}

	o.exec(ctx, execItems)

	return execItems[0].err
}

func (o *PoolOp) exec(ctx context.Context, items []execItem) {
	// create execs
	execs := make([]nodeExec, 0, len(items))
	o.mu.Lock()
	for _, item := range items {
		var nodeExec nodeExec
		var exists bool
		if nodeExec, exists = o.nodeExecs[item.node.ID]; !exists {
			nodeOp := func(ctx context.Context) (*empty, error) {
				err := o.nodeOp.Exec(ctx, item.node, o.log)
				return nil, err
			}
			nodeExec = waveexec.New(nodeOp)
			o.nodeExecs[item.node.ID] = nodeExec
		}
		execs = append(execs, nodeExec)
	}
	o.mu.Unlock()

	// run execs
	var wg sync.WaitGroup
	for idx, exec := range execs {
		wg.Add(1)
		items[idx].err = safego.Invoke(func() error {
			defer func() {
				wg.Done()
			}()
			_, err := exec.Invoke(ctx)
			return err
		})
	}
	wg.Wait()
}
