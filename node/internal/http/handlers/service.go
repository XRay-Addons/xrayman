package handlers

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

//go:generate mockgen -destination=./mocks/service_mock.go -package=mocks . Service

type Service interface {
	Start(ctx context.Context, users []models.User) (*models.NodeProperties, error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (*models.NodeStatus, error)
	EditUsers(ctx context.Context, add, remove []models.User) error
}
