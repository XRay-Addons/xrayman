package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
)

type PoolClient interface {
	GetNodeClient(ctx context.Context, conn models.NodeConnectionInfo) (
		nodesyncer.NodeClient, error)
}
