package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayService interface {
	Start(ctx context.Context, config string) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.ServiceStatus, error)
}
