package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=auth_service.go -destination=./mocks/mock_auth_service.go -package=mocks
type AuthService interface {
	Auth(ctx context.Context, p models.AuthParams) (*models.AuthResult, error)
}
