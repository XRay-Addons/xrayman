package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	uow    UoW
	sync   SyncService
	keygen Keygen
}

var _ handler.Service = (*Service)(nil)

func New(uow UoW, sync SyncService, keygen Keygen) (*Service, error) {
	if uow == nil {
		return nil, fmt.Errorf("service init: uow: %w", errdefs.ErrNilArgPassed)
	}
	if sync == nil {
		return nil, fmt.Errorf("service init: sync: %w", errdefs.ErrNilArgPassed)
	}
	if keygen == nil {
		return nil, fmt.Errorf("service init: keygen: %w", errdefs.ErrNilArgPassed)
	}
	return &Service{
		uow:    uow,
		sync:   sync,
		keygen: keygen,
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
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
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
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusRunning); err != nil {
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
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
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
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, status)
		return
	}); err != nil {
		return fmt.Errorf("set node status: %w", err)
	}

	err := s.syncNode(ctx, id)
	if err != nil {
		return fmt.Errorf("set node status: %w", err)
	}
	return nil
}

func (s *Service) syncNode(ctx context.Context, id models.NodeID) error {
	syncResults, err := s.sync.SyncNodesPool(ctx)
	if err != nil {
		return fmt.Errorf("service: sync node: %w", err)
	}
	for _, syncRes := range syncResults {
		if syncRes.ID != id {
			continue
		}
		if syncRes.Err == nil {
			return nil
		}
		return fmt.Errorf("service: sync node: %w", syncRes.Err)
	}
	return fmt.Errorf("servuce: sync node: node not found: %w", errdefs.ErrIPE)
}

// sync all nodes, return nil if at least one node synced ok
func (s *Service) syncAllNodes(ctx context.Context) error {
	syncResults, err := s.sync.SyncNodesPool(ctx)
	if err != nil {
		return fmt.Errorf("service: sync all nodes: %w", err)
	}
	if len(syncResults) == 0 {
		return nil
	}
	var errs []error
	for _, syncRes := range syncResults {
		if syncRes.Err == nil {
			return nil
		}
		errs = append(errs, syncRes.Err)
	}
	return fmt.Errorf("servuce: sync all nodes: %w", errors.Join(errs...))
}
