package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type ServerCfg interface {
	GetUsersCfg(users []models.User) (string, error)
}

type ClientCfg interface {
	Get() (*models.ClientCfg, error)
}

type XRayService interface {
	Start(ctx context.Context, config string) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.ServiceStatus, error)
}

type XRayAPI interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	EditUsers(ctx context.Context, add, remove []models.User) error
}
