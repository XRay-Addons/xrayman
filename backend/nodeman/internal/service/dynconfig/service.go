package dynconfig

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	storage Storage
}

//var _ handler.NodesService = (*Service)(nil)

func New(storage Storage) (*Service, error) {
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}

	return &Service{
		storage: storage,
	}, nil
}

func (s *Service) GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	cfg, err := s.storage.GetDynamicConfig(ctx)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Service) SetDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error {
	if s == nil {
		return errdefs.NilCall()
	}
	if err := s.storage.SetDynamicConfig(ctx, cfg); err != nil {
		return err
	}
	return nil
}

func (s *Service) EnsureDefaultConfig(ctx context.Context) error {
	if s == nil {
		return errdefs.NilCall()
	}
	if err := s.storage.EnsureDynamicConfig(ctx, models.DynamicConfig{}); err != nil {
		return err
	}
	return nil
}
