package pooljob

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolOp = func(ctx context.Context) (*models.PoolOpResult, error)
