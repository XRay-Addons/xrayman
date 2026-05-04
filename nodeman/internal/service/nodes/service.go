package nodes

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	storage    Storage
	poolSyncer poolsync.Syncer
}

var _ handler.NodesService = (*Service)(nil)

func New(poolSyncer poolsync.Syncer,
	storage Storage,
) (*Service, error) {
	if poolSyncer == nil {
		return nil, errdefs.NewNilArg("poolSyncer")
	}
	if storage == nil {
		return nil, errdefs.NewNilArg("storage")
	}

	return &Service{
		storage:    storage,
		poolSyncer: poolSyncer,
	}, nil
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (
	*models.NewNodeResult, error,
) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var node models.Node
	node.Config.ConnectionInfo.Endpoint = p.Endpoint
	node.Config.ConnectionInfo.AccessKey = p.AccessKey

	node.CurrentStatus = models.NodeStatusStopped
	node.TargetStatus = models.NodeStatusRunning
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewNode(ctx, &node)
		return
	}); err != nil {
		return nil, err
	}

	_ = s.syncNode(ctx, node.ID)

	return &models.NewNodeResult{
		Node: node,
	}, nil
}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) (
	*models.StartNodeResult, error,
) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusRunning); err != nil {
		return nil, err
	}
	return &models.StartNodeResult{}, nil
}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeParams) (
	*models.StopNodeResult, error,
) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, err
	}
	return &models.StopNodeResult{}, nil
}

func (s *Service) ListNodes(ctx context.Context, p models.ListNodeParams) (
	*models.ListNodeResult, error,
) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var nodes []models.Node
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		nodes, err = uowctx.ListNodes(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) DeleteNode(ctx context.Context, p models.DeleteNodeParams) (
	*models.DeleteNodeResult, error,
) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}

	// stop node before deleting
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, err
	}

	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.DeleteNode(ctx, p.ID)
		return
	}); err != nil {
		return nil, err
	}

	return &models.DeleteNodeResult{}, nil
}

func (s *Service) setNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	// set target node state to storage
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, status)
		return
	}); err != nil {
		return err
	}

	_ = s.syncNode(ctx, id)
	return nil
}

func (s *Service) syncNode(ctx context.Context, id models.NodeID) error {
	syncResults, err := s.poolSyncer.SyncPoolState(ctx)
	if err != nil {
		return err
	}
	if err = syncResults.GetNodeErr(id); err != nil {
		return err
	}
	return nil
}

// sync all nodes, return nil if at least one node synced ok
func (s *Service) syncAllNodes(ctx context.Context) error {
	syncResults, err := s.poolSyncer.SyncPoolState(ctx)
	if err != nil {
		return err
	}
	if len(syncResults.Nodes) == 0 {
		return nil
	}
	var errs []error
	for _, syncRes := range syncResults.Nodes {
		if syncRes.Err == nil {
			return nil
		}
		errs = append(errs, syncRes.Err)
	}
	return errors.Join(errs...)
}
