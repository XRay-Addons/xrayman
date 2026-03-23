package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/nodesyncer"
)

type Client interface {
	GetNodeClient(ctx context.Context,
		conn models.NodeConnectionInfo) (nodesyncer.Client, error)
}
