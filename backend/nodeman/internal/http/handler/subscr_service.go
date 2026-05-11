package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=subscr_service.go -destination=./mocks/mock_subscr_service.go -package=mocks
type SubscrService interface {
	GetUserSub(ctx context.Context, p models.UserSubParams) (*models.UserSubResult, bool, error)

	NewHeader(ctx context.Context, p models.NewSubHeaderParams) (*models.Header, error)
	ListHeaders(ctx context.Context, p models.ListSubHeadersParams) (*models.ListSubHeadersResult, error)
	DeleteHeader(ctx context.Context, p models.DeleteSubHeaderParams) (*models.DeleteSubHeaderResult, error)
}
