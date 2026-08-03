package settings

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	EnsureSettings(ctx context.Context, cfg models.Settings) error
	GetSettings(ctx context.Context) (*models.Settings, error)
	SetSettings(ctx context.Context, cfg models.Settings) error
}
