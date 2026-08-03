package nodes

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/supervisor"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	storage    Storage
	poolSyncer Syncer

	syncTimeout time.Duration
	sv          *supervisor.Supervisor

	logger *zap.Logger
}

var _ handler.NodesService = (*Service)(nil)

func New(poolSyncer Syncer,
	storage Storage,
	syncTimeout time.Duration,
	logger *zap.Logger,
) (*Service, error) {
	if poolSyncer == nil {
		return nil, errdefs.NilArg("poolSyncer")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	if logger == nil {
		return nil, errdefs.NilArg("logger")
	}

	return &Service{
		storage:     storage,
		poolSyncer:  poolSyncer,
		syncTimeout: syncTimeout,
		sv:          supervisor.New(),
		logger:      logger,
	}, nil
}

func (s *Service) Close() {
	if s == nil || s.sv == nil {
		return
	}
	s.sv.Close()
}

func (s *Service) requestNodeSync(id models.NodeID) {
	s.sv.Go(func(ctx context.Context) {
		if err := s.poolSyncer.SyncNodeState(ctx, id); err != nil {
			s.logger.Warn("node sync request", zap.Error(err))
		}
	}, s.syncTimeout)
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (
	*models.NewNodeResult, error,
) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	var node models.Node
	node.Config.ConnectionInfo.Endpoint = p.Endpoint
	node.Config.ConnectionInfo.AccessKey = p.AccessKey

	node.CurrentStatus = models.NodeStatusStopped
	node.TargetStatus = models.NodeStatusRunning
	if err := s.storage.NewNode(ctx, &node); err != nil {
		return nil, err
	}

	s.requestNodeSync(node.ID)

	return &models.NewNodeResult{
		Node: node,
	}, nil
}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) error {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusRunning); err != nil {
		return err
	}
	return nil
}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeParams) error {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return err
	}
	return nil
}

func (s *Service) ListNodes(ctx context.Context) (
	*models.ListNodeResult, error,
) {
	nodes, err := s.storage.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) DeleteNode(ctx context.Context, p models.DeleteNodeParams) error {
	// mark node stopped and deleting
	if err := s.storage.DoTx(ctx, func(ctx context.Context) error {
		if err := s.storage.SetTargetNodeStatus(ctx,
			p.ID, models.NodeStatusStopped,
		); err != nil {
			return err
		}
		if err := s.storage.DeleteNode(ctx, p.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	s.requestNodeSync(p.ID)

	return nil
}

func (s *Service) setNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	if s == nil {
		return errdefs.NilCall()
	}
	// set target node state to storage
	if err := s.storage.SetTargetNodeStatus(ctx, id, status); err != nil {
		return err
	}

	s.requestNodeSync(id)

	return nil
}
