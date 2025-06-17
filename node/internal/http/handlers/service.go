package handlers

import (
	"context"

	"github.com/XRay-Addons/xrayman/shared/models"
)

type Service interface {
	Start(ctx context.Context, users []models.User) (*models.NodeProperties, error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (*models.NodeStatus, error)
	AddUsers(ctx context.Context, users []models.User) error
	DelUsers(ctx context.Context, users []models.User) error
}
