package poolstats

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/nodestats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Client interface {
	GetNodeClient(conn models.NodeConnectionInfo) (nodestats.Client, error)
}
