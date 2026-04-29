package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=service.go -destination=./mocks/mock_subscr_service.go -package=mocks
type SubscrService interface {
	GetUserSub(ctx context.Context, p models.UserSubParams) (*models.UserSubResult, bool, error)
}
