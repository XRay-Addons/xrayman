package subheaders

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	NewSubHeader(ctx context.Context, header *models.Header) error
	DeleteSubHeader(ctx context.Context, id models.HeaderID) error
	ListSubHeaders(ctx context.Context) ([]models.Header, error)
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
