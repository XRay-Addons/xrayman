package converter

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Converter interface {
	ConvertStartRequest(source *api.StartRequest) (*models.StartParams, error)
	ConvertStartResult(source *models.StartResult) *api.StartResponse
	ConvertEditUsersRequest(source *api.EditUsersRequest) (*models.EditUsersParams, error)
	ConvertStatusResult(source *models.StatusResult) *api.StatusResponse
	ConvertStatus(source models.ServiceStatus) api.ServiceStatus
}
