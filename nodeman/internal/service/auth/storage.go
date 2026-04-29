package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
)


type Auth struct {
	PasswordHash []byte
}

type UoWContext interface {
	GetAdmin(context.Context) (*Auth, error)
	SetAdmin(context.Context, *Auth) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
