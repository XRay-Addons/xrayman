package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	GetAuth(ctx context.Context) (*models.Auth, error)
	SetAuth(ctx context.Context, auth *models.Auth) error
}
