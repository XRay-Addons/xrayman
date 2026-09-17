package nodes

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type SyncService interface {
	SyncNodeState(ctx context.Context, id models.NodeID) error
}
