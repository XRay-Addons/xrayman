package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	GetAdmin(context.Context) (*models.Auth, error)
	SetAdmin(context.Context, *models.Auth) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
