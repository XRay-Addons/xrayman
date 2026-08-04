package converter

import (
	"fmt"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/openapi-gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
//
//go:generate goverter gen .
type Converter interface {
	ConvertUsers(users []models.UserProfile) []api.User
	ConvertStartResponse(cfg api.StartResponse) models.NodeSettings
	ConvertUsersUpdate(users models.NodeUsersUpdate) api.EditUsersRequest
	ConvertNodeStats(stats *api.StatsResponse) *models.NodeStats
}

func ConvertNodeStatus(s api.ServiceStatus) models.NodeStatus {
	switch s {
	case api.ServiceStatusUnknown:
		return models.NodeStatusUnknown
	case api.ServiceStatusRunning:
		return models.NodeStatusRunning
	case api.ServiceStatusStopped:
		return models.NodeStatusStopped
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", s))
	}
}
