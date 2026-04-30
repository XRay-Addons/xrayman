package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
	GetUser(ctx context.Context, id models.UserID) (*models.User, bool, error)
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
