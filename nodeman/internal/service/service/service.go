package service

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	storage Storage
	keygen  Keygen
	poolmon PoolMonitor
}

var _ handler.Service = (*Service)(nil)

func New(storage Storage, keygen Keygen, poolmon PoolMonitor) (*Service, error) {
	if storage == nil {
		return nil, fmt.Errorf("service init: uow: %w", errdefs.ErrNilArgPassed)
	}
	if keygen == nil {
		return nil, fmt.Errorf("service init: keygen: %w", errdefs.ErrNilArgPassed)
	}
	if poolmon == nil {
		return nil, fmt.Errorf("service init: poolmon: %w", errdefs.ErrNilArgPassed)
	}
	return &Service{
		storage: storage,
		keygen:  keygen,
		poolmon: poolmon,
	}, nil
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: start: %w", errdefs.ErrNilObjectCall)
	}
	accessSecret, err := s.keygen.GenerateHS256Secret()
	if err != nil {
		return nil, fmt.Errorf("new node: %w", err)
	}
	var node models.Node
	node.Config.ConnectionInfo.AccessSecret = accessSecret
	node.Config.ConnectionInfo.Endpoint = p.Endpoint
	node.CurrentStatus = models.NodeStatusStopped
	node.TargetStatus = models.NodeStatusStopped
	if err := s.storage.NewUoW().Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewNode(ctx, &node)
		return
	}); err != nil {
		return nil, fmt.Errorf("service: new node: %w", err)
	}
	return &models.NewNodeResult{
		ID:           node.ID,
		Endpoint:     node.Config.ConnectionInfo.Endpoint,
		AccessSecret: node.Config.ConnectionInfo.AccessSecret,
	}, nil
}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, fmt.Errorf("service: start node: %w", err)
	}
	return &models.StartNodeResult{}, nil
}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeParams) (*models.StopNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, fmt.Errorf("service: stop node: %w", err)
	}
	return &models.StopNodeResult{}, nil
}

func (s *Service) ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: list nodes: %w", errdefs.ErrNilObjectCall)
	}
	var nodes []models.Node
	if err := s.storage.NewUoW().Do(ctx, func(uowctx UoWContext) (err error) {
		nodes, err = uowctx.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("set node status: %w", err)
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) setNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	if s == nil {
		return fmt.Errorf("set node status: %w", errdefs.ErrNilObjectCall)
	}
	// set target node state to storage
	if err := s.storage.NewUoW().Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, models.NodeStatusRunning)
		return
	}); err != nil {
		return fmt.Errorf("set node status: %w", err)
	}

	res, err := s.poolmon.Sync(ctx)
	if err != nil {
		return fmt.Errorf("set node status: %w", err)
	}
	if err := res.GetNodeErr(id); err != nil {
		return fmt.Errorf("set node status: %w", err)
	}
	return nil
}
