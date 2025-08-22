package supervisor

import (
	"runtime"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/launchctl"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) option {
	return func(o *options) {
		if logger == nil {
			return
		}
		o.log = logger
	}
}

type option func(o *options)

type options struct {
	log *zap.Logger
}

func New(serviceName string, command []string, opts ...option) (supervisorapi.Supervisor, error) {
	o := &options{
		log: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(o)
	}

	switch runtime.GOOS {
	case "darwin":
		return launchctl.New(serviceName, command, launchctl.WithLogger(o.log))
	default:
		return nil, errdefs.New("unsupported platform",
			errdefs.Withf("platform: %v", runtime.GOOS))
	}
}
