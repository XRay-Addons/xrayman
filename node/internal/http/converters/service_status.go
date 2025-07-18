package converters

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
)

func ServiceStatusFromAPI(s api.ServiceStatus) models.ServiceStatus {
	switch s {
	case api.ServiceRunning:
		return models.ServiceRunning
	case api.ServiceStopped:
		return models.ServiceStopped
	default:
		panic("unknown status")
	}
}

func ServiceStatusToAPI(s models.ServiceStatus) api.ServiceStatus {
	switch s {
	case models.ServiceRunning:
		return api.ServiceRunning
	case models.ServiceStopped:
		return api.ServiceStopped
	default:
		panic("unknown status")
	}
}
