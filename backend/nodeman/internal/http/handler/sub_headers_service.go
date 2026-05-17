package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=sub_headers_service.go -destination=./mocks/mock_sub_headers_service.go -package=mocks
type SubHeadersService interface {
	NewHeader(ctx context.Context, p models.NewSubHeaderParams) (*models.Header, error)
	DeleteHeader(ctx context.Context, p models.DeleteSubHeaderParams) (*models.DeleteSubHeaderResult, error)
	ListHeaders(ctx context.Context, p models.ListSubHeadersParams) (*models.ListSubHeadersResult, error)
}
