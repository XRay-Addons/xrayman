package pool

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeClient interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfig, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUsers(ctx context.Context, upd models.NodeUsersUpdate) error
}
