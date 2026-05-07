package password

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	GetAuth(ctx context.Context) (*models.Auth, error)
	SetAuth(ctx context.Context, auth models.Auth) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
