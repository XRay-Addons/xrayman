package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayCtl interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.Status, error)
}
