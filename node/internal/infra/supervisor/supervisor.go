package supervisor

import (
	"fmt"
	"runtime"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/launchctl"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"
	"go.uber.org/zap"
)

func New(serviceName string, command []string, log *zap.Logger) (supervisorapi.Supervisor, error) {
	switch runtime.GOOS {
	case "darwin":
		return launchctl.New(serviceName, command, log)
	default:
		return nil, fmt.Errorf("supervisor: %w", errdefs.ErrUnsupportedPlatform)
	}
}
