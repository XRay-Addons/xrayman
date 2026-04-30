package poolsync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	GetNodeClient(ctx context.Context,
		conn models.NodeConnectionInfo) (nodesync.Client, error)
}
