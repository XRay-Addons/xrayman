package subscr

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
	GetUserView(ctx context.Context, id models.UserID, name string) (*models.UserView, error)

	GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error)

	//NewSubHeader(ctx context.Context, header *models.Header) error
	//DeleteSubHeader(ctx context.Context, id models.HeaderID) error
	//ListSubHeaders(ctx context.Context) ([]models.Header, error)
}
