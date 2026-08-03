package users

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Syncer interface {
	SyncPoolState(ctx context.Context) (*models.PoolOpResult, error)
}
