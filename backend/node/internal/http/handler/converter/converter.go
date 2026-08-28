package converter

import (
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Converter interface {
	ConvertStartRequest(source *ogenserver.StartRequest) *models.StartParams
	ConvertStartResult(source *models.StartResult) *ogenserver.StartResponse
	ConvertEditUsersRequest(source *ogenserver.EditUsersRequest) *models.EditUsersParams
	ConvertStatusResult(source *models.StatusResult) *ogenserver.StatusResponse
	ConvertStatus(source models.ServiceStatus) ogenserver.ServiceStatus
	ConvertStatsResult(source *models.StatsResult) *ogenserver.StatsResponse
}
