package bootstrap

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type DynConfig interface {
	EnsureDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error
}
