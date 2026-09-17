package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/poolop"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/zap"
)

type Service struct {
	op *poolop.PoolOp
}

var _ users.SyncService = (*Service)(nil)
var _ nodes.SyncService = (*Service)(nil)

//var _ syncman.PoolSyncer = (*Syncer)(nil)

func New(client Client, storage Storage, log *zap.Logger) (*Service, error) {
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
	return &Service{
		op: op,
	}, nil
}

func (s *Service) Close() {
	if s == nil || s.op == nil {
		return
	}
	s.op.Close()
}

func (s *Service) SyncPoolState(ctx context.Context) (
	*models.PoolOpResult, error,
) {
	return s.op.ExecAll(ctx)
}

func (s *Service) SyncNodeState(ctx context.Context,
	id models.NodeID,
) error {
	return s.op.ExecNode(ctx, id)
}
