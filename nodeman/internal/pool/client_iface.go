package pool

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	GetNodeClient(ctx context.Context,
		conn models.NodeConnectionInfo) (NodeClient, error)
}
