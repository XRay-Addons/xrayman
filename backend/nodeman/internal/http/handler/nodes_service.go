package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=nodes_service.go -destination=./mocks/mock_nodes_service.go -package=mocks
type NodesService interface {
	NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error)
	StartNode(ctx context.Context, p models.StartNodeParams) error
	StopNode(ctx context.Context, p models.StopNodeParams) error
	ListNodes(ctx context.Context) (*models.ListNodeResult, error)
	DeleteNode(ctx context.Context, p models.DeleteNodeParams) error
}
