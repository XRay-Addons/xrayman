package poolop

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Storage interface {
	ListNodes(ctx context.Context) ([]models.Node, error)
	GetNode(ctx context.Context, id models.NodeID) (*models.Node, error)
}
