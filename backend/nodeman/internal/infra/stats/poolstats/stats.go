package poolstats

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/nodestats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type Stats struct {
	op *poolop.PoolOp
}

var _ statsman.StatsUpdater = (*Stats)(nil)

func New(client Client, storage Storage, log *zap.Logger) (*Stats, error) {
	if client == nil {
		return nil, errdefs.NilArg("client")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	op, err := poolop.New(
		&nodePool{storage: storage},
		&nodeOp{storage: storage, client: client},
		log,
	)
	if err != nil {
		return nil, err
	}
	return &Stats{
		op: op,
	}, nil
}

func (s *Stats) Close() {
	if s == nil || s.op == nil {
		return
	}
	s.op.Close()
}

func (s *Stats) UpdatePoolStats(ctx context.Context) (*models.PoolOpResult, error) {
	return s.op.Exec(ctx)
}

// node pool impl
type nodePool struct {
	storage Storage
}

var _ poolop.NodePool = (*nodePool)(nil)

func (p *nodePool) ListNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	if err := p.storage.DoUoW(ctx, func(uow UoWContext) (err error) {
		nodes, err = uow.ListNodes(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return nodes, nil
}

// node op impl
type nodeOp struct {
	storage Storage
	client  Client
}

var _ poolop.NodeOp = (*nodeOp)(nil)

func (op *nodeOp) Exec(ctx context.Context, node models.Node, log *zap.Logger) error {
	// get stats only from running nodes
	if node.TargetStatus != models.NodeStatusRunning {
		return nil
	}

	nodeStorage := &nodeStorage{
		base:   op.storage,
		nodeID: node.ID,
	}
	nodeClient, err := op.client.GetNodeClient(node.Config.ConnectionInfo)
	if err != nil {
		return err
	}
	if err := nodestats.UpdateNodeStats(ctx, nodeClient, nodeStorage, log); err != nil {
		return err
	}
	return nil
}
