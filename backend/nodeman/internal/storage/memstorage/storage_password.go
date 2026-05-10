package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) GetAuth(ctx context.Context) (*models.Auth, error) {
	return &models.Auth{
		PasswordHash: s.adminPass,
	}, nil
}

func (s *Storage) SetAuth(ctx context.Context, a models.Auth) error {
	s.adminPass = a.PasswordHash
	return nil
}
