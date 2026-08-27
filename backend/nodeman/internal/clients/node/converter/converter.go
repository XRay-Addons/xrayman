package converter

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/ogenclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
//
//go:generate goverter gen .
type Converter interface {
	ConvertUsers(users []models.UserProfile) []ogenclient.User
	ConvertStartResponse(cfg ogenclient.StartResponse) models.NodeSettings
	ConvertUsersUpdate(users models.NodeUsersUpdate) ogenclient.EditUsersRequest
	ConvertNodeStats(stats *ogenclient.StatsResponse) *models.NodeStats
}

func ConvertNodeStatus(s ogenclient.ServiceStatus) models.NodeStatus {
	switch s {
	case ogenclient.ServiceStatusUnknown:
		return models.NodeStatusUnknown
	case ogenclient.ServiceStatusRunning:
		return models.NodeStatusRunning
	case ogenclient.ServiceStatusStopped:
		return models.NodeStatusStopped
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", s))
	}
}
