package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=dyn_config_service.go -destination=./mocks/mock_dynconfig_service.go -package=mocks
type DynConfigService interface {
	SetDynamicConfig(ctx context.Context, p models.DynamicConfig) error
	GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error)
}
