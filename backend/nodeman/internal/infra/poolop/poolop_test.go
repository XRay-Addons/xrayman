package poolop

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type nodePool struct {
	nodes []models.Node
}

func (p *nodePool) ListNodes(ctx context.Context) ([]models.Node, error) {
	return p.nodes, nil
}

func (p *nodePool) GetNode(ctx context.Context, id models.NodeID) (*models.Node, error) {
	for _, node := range p.nodes {
		if node.ID == id {
			return &node, nil
		}
	}
	return nil, xerr.New("Node not found")
}

type nodeOp struct {
	deferCallsCount int
	lock            sync.Mutex
}

func (o *nodeOp) Exec(ctx context.Context, node models.Node, log *zap.Logger) error {
	defer func() {
		log.Info(fmt.Sprintf("defer for node: %d\n", node.ID))
		o.lock.Lock()
		defer o.lock.Unlock()
		o.deferCallsCount++
	}()

	if node.ID == 1 {
		panic("node panic")
	}
	return nil
}

func TestPoolOp(t *testing.T) {
	np := nodePool{
		nodes: make([]models.Node, 4),
	}
	for i := range np.nodes {
		np.nodes[i].ID = i
	}

	op := nodeOp{}
	log := zaptest.NewLogger(t)
	poolOp, err := New(&np, &op, log)
	require.NoError(t, err)

	res, err := poolOp.ExecAll(t.Context())
	require.NoError(t, err)
	require.Equal(t, len(np.nodes), len(res.Nodes))
	require.Error(t, res.Nodes[1].Err)
	log.Info("panic test error", zap.Error(res.Nodes[1].Err))
	require.Equal(t, len(np.nodes), op.deferCallsCount)
}
