package syncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolClient interface {
	GetNodeClient(ctx context.Context, conn models.NodeConnectionInfo) (NodeClient, error)
}
