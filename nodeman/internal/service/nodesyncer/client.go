package nodesyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfig, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUserStates(ctx context.Context, update models.NodeUsersUpdate) error
}
