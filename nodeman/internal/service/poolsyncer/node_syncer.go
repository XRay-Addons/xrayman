package poolsyncer

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
)

type NodesSyncer interface {
	SyncNode(ctx context.Context,
		uow nodesyncer.UoW,
		client nodesyncer.NodeClient) error
}
