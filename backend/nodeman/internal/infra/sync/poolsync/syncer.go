package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/zap"
)

type Syncer struct {
	op *poolop.PoolOp
}

var _ users.Syncer = (*Syncer)(nil)
var _ nodes.Syncer = (*Syncer)(nil)
var _ syncman.PoolSyncer = (*Syncer)(nil)

func New(client Client, storage Storage, log *zap.Logger) (*Syncer, error) {
	if client == nil {
		return nil, errdefs.NilArg("client")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	if log == nil {
		return nil, errdefs.NilArg("log")
	}

	op, err := poolop.New(
		storage,
		&nodeOp{storage: storage, client: client},
		log,
	)
	if err != nil {
		return nil, err
	}
	return &Syncer{
		op: op,
	}, nil
}

func (s *Syncer) Close() {
	if s == nil || s.op == nil {
		return
	}
	s.op.Close()
}

func (s *Syncer) SyncPoolState(ctx context.Context) (
	*models.PoolOpResult, error,
) {
	return s.op.ExecAll(ctx)
}

func (s *Syncer) SyncNodeState(ctx context.Context,
	id models.NodeID,
) (*models.NodeOpResult, error) {
	return s.op.ExecNode(ctx, id)
}
