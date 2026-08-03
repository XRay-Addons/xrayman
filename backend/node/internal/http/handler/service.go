package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

//go:generate mockgen -source=service.go -destination=./mocks/mock_service.go -package=mocks
type Service interface {
	Start(ctx context.Context, params models.StartParams) (*models.StartResult, error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (*models.StatusResult, error)
	EditUsers(ctx context.Context, params models.EditUsersParams) error
	GetStats(ctx context.Context) (*models.StatsResult, error)
}
