package dynconfig

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	EnsureDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error
	GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error)
	SetDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
