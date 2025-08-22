package supervisor

import (
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
		return nil, errdefs.New("unsupported platform",
			errdefs.Withf("platform: %v", runtime.GOOS))
	}
}
