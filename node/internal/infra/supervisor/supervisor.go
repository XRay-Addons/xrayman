package supervisor

import (
	"context"
	"fmt"
	"runtime"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/launchctl"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/zap"
)

type Supervisor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (models.ServiceStatus, error)
	Close(ctx context.Context) error
}

func New(serviceName string, command []string, log *zap.Logger) (Supervisor, error) {
	switch runtime.GOOS {
	case "darwin":
		return launchctl.New(serviceName, command, log)
	default:
		return nil, fmt.Errorf("supervisor: %w", errdefs.ErrUnsupportedPlatform)
	}
}
