package dynconfig

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	EnsureDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error
	GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error)
	SetDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error
}
