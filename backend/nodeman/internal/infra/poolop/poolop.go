package poolop

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/waveexec"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type NodePool interface {
	ListNodes(ctx context.Context) ([]models.Node, error)
}

type NodeOp interface {
	Exec(ctx context.Context, node models.Node, log *zap.Logger) error
}

type PoolOpFn = func(ctx context.Context) (*models.PoolOpResult, error)

type PoolOp struct {
	exec *waveexec.WaveExecutor[models.PoolOpResult]
}

func New(pool NodePool, op NodeOp, log *zap.Logger) (*PoolOp, error) {
	if pool == nil {
		return nil, xerr.NilArg("pool")
	}
	if op == nil {
		return nil, xerr.NilArg("op")
	}
	if log == nil {
		return nil, xerr.NilArg("log")
	}
	poolOpFn := getPoolOp(pool, op, log)
	return &PoolOp{
		exec: waveexec.New(poolOpFn),
	}, nil
}

func (op *PoolOp) Close() {
	if op == nil || op.exec == nil {
		return
	}
	op.exec.Close()
}

func (op *PoolOp) Exec(ctx context.Context) (*models.PoolOpResult, error) {
	return op.exec.Invoke(ctx)
}

func getPoolOp(pool NodePool, op NodeOp, log *zap.Logger) PoolOpFn {
	return func(ctx context.Context) (*models.PoolOpResult, error) {
		return poolOp(ctx, pool, op, log)
	}
}

func poolOp(ctx context.Context,
	pool NodePool, op NodeOp, log *zap.Logger,
) (*models.PoolOpResult, error) {
	nodes, err := pool.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodeResults := make([]models.NodeOpResult, len(nodes))
	var wg sync.WaitGroup
	for idx, node := range nodes {
		nodeResults[idx].ID = node.ID
		nodeResults[idx].Endpoint = node.Config.ConnectionInfo.Endpoint

		wg.Add(1)
		go func() {
			defer func() {
				// panic to error
				if p := recover(); p != nil {
					nodeResults[idx].Err = xerr.Panic(p)
				}
				defer wg.Done()
			}()

			// PANIC
			nodeResults[idx].Err = op.Exec(ctx, node, log)
		}()
	}
	wg.Wait()

	return &models.PoolOpResult{
		Nodes: nodeResults,
	}, nil
}
