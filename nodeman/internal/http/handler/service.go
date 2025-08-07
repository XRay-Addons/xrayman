package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=service.go -destination=./mocks/mock_service.go -package=mocks
type Service interface {
	NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error)
	StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error)
	StopNode(ctx context.Context, p models.StopNodeParams) (*models.StopNodeResult, error)
	ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error)
}
