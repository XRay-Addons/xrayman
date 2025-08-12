package pool

import (
	"context"
)

type NodeSyncer interface {
	SyncNodeState(ctx context.Context, client NodeClient, uow NodeUoW) error
}
