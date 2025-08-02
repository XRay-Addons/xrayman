package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type IService interface {
	NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error)
	StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error)
	StopNode(ctx context.Context, p models.StopNodeResult) (*models.StopNodeResult, error)
	ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error)
}

type Service struct {
	storage Storage
	keygen  Keygen
	poolmon PoolMonitor
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error) {

}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error) {

}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeResult) (*models.StopNodeResult, error) {

}

func (s *Service) ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error) {

}
