package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/pkg/api/models"
)

type XRayApi interface {
	EditUsers(ctx context.Context, ins []models.Inbound,
		add, remove []models.User) error

	Ping(ctx context.Context) error
}
