package nodesyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfigTemplate, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUsers(ctx context.Context, upd models.NodeUsersUpdate) error
}
