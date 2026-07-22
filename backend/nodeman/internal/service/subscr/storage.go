package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
	GetUserView(ctx context.Context, id models.UserID, name string) (*models.UserView, error)

	GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error)
}
