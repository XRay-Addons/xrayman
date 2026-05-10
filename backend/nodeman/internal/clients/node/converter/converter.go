package converter

import (
	"fmt"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
//
//go:generate goverter gen .
type Converter interface {
	ConvertUsers(users []models.UserProfile) []api.User
	ConvertClientConfig(cfg api.ClientConfigTemplate) models.ClientConfigTemplate
	ConvertUsersUpdate(users models.NodeUsersUpdate) api.EditUsersRequest
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
