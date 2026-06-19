package subheaders

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	NewSubHeader(ctx context.Context, header *models.Header) error
	DeleteSubHeader(ctx context.Context, id models.HeaderID) error
	ListSubHeaders(ctx context.Context) ([]models.Header, error)
}
