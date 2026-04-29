package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) GetAdmin(ctx context.Context) (*models.Auth, error) {
	return &models.Auth{
		PasswordHash: s.adminPass,
	}, nil
}

func (s *Storage) SetAdmin(ctx context.Context, a *models.Auth) error {
	s.adminPass = a.PasswordHash
	return nil
}
