package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
)

func (s *Storage) GetAdmin(ctx context.Context) (*auth.Auth, error) {
	return &auth.Auth{
		PasswordHash: s.adminPass,
	}, nil
}

func (s *Storage) SetAdmin(ctx context.Context, a *auth.Auth) error {
	s.adminPass = a.PasswordHash
	return nil
}
