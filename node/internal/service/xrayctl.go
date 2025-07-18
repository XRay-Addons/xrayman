package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/shared/models"
)

type XRayCtl interface {
	Start(ctx context.Context, config string) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.Status, error)
}
