package poolsync

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	GetNodeClient(conn models.NodeConnectionInfo) (nodesync.Client, error)
}
