package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type ServerCfg interface {
	GetInbounds() ([]models.Inbound, error)
	GetApiURL() (string, error)
	GetUsersCfg(users []models.User) (string, error)
}

type ClientCfg interface {
	Get() (*models.ClientCfg, error)
}

type XRayServiceCtl interface {
	Start(ctx context.Context, config string) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.ServiceStatus, error)
}

type XRayAPI interface {
	EditUsers(ctx context.Context, add, remove []models.User) error
}
