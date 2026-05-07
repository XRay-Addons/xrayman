package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayAPI interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	EditUsers(ctx context.Context, add, remove []models.User) error
}
