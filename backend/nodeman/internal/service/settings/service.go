package settings

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	storage Storage
}

var _ handler.SettingsService = (*Service)(nil)

func New(storage Storage) (*Service, error) {
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}

	return &Service{
		storage: storage,
	}, nil
}

func (s *Service) GetSettings(ctx context.Context) (*models.Settings, error) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	settings, err := s.storage.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *Service) SetSettings(ctx context.Context, settings models.Settings) error {
	if s == nil {
		return errdefs.NilCall()
	}
	if err := s.storage.SetSettings(ctx, settings); err != nil {
		return err
	}
	return nil
}

func (s *Service) EnsureSettings(ctx context.Context) error {
	if s == nil {
		return errdefs.NilCall()
	}
	if err := s.storage.EnsureSettings(ctx, models.Settings{}); err != nil {
		return err
	}
	return nil
}
