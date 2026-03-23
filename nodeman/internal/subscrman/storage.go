package subscrman

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
