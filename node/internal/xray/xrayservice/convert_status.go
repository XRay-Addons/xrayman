package xrayservice

import (
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

func convertStatus(s supervisorapi.ServiceStatus) (models.ServiceStatus, error) {
	switch s {
	case supervisorapi.StatusStopped:
		return models.ServiceStopped, nil
	case supervisorapi.StatusRunning:
		return models.ServiceRunning, nil
	default:
		return models.ServiceStopped, errdefs.New("unknown service status",
			errdefs.Withf("status: %v", s))
	}
}
