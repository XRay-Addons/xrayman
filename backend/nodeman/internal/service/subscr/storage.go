package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
	GetUser(ctx context.Context, id models.UserID) (*models.User, bool, error)

	NewSubHeader(ctx context.Context, header *models.Header) error
	DeleteSubHeader(ctx context.Context, id models.HeaderID) error
	ListSubHeaders(ctx context.Context) ([]models.Header, error)
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
