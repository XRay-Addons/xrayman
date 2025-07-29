package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeClient interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfig, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUserStates(ctx context.Context, update models.NodeUsersUpdate) error
}

type NodeClientFactory interface {
	Get(ctx context.Context, endpoint string, secret []byte) (NodeClient, error)
}
