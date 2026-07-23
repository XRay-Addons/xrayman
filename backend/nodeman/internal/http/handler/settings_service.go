package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=settings_service.go -destination=./mocks/mock_settings_service.go -package=mocks
type SettingsService interface {
	SetSettings(ctx context.Context, s models.Settings) error
	GetSettings(ctx context.Context) (*models.Settings, error)
}
