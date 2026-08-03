package poolop

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type NodeOp interface {
	Exec(ctx context.Context, node models.Node, log *zap.Logger) error
}
