package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayAPI interface {
	EditUsers(ctx context.Context, add, remove []models.User) error
}
