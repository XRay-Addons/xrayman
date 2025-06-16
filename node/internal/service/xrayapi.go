package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayApi interface {
	AddUsers(ctx context.Context, ins []models.Inbound, users []models.User) error
	DelUsers(ctx context.Context, ins []models.Inbound, users []models.User) error
	Ping(ctx context.Context) error
}
