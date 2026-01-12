package handler

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
type Converter interface {
	ConvertStartRequest(source *api.StartRequest) *models.StartParams
	ConvertStartResult(source *models.StartResult) *api.StartResponse
	ConvertEditUsersRequest(source *api.EditUsersRequest) *models.EditUsersParams
}

// gonverter can't generate it :((
func ConvertStatusResult(source *models.StatusResult) *api.StatusResponse {
	var response api.StatusResponse
	switch source.ServiceStatus {
	case models.ServiceRunning:
		response.ServiceStatus = api.ServiceStatusRunning
	case models.ServiceStopped:
		response.ServiceStatus = api.ServiceStatusStopped
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", source))
	}
	return &response
}
