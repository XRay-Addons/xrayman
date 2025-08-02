package handler

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:typedef-to-base:true

//go:generate goverter gen .
type Converter interface {
	ConvertNewNodeRequest(source *api.NewNodeRequest) *models.NewNodeParams
	ConvertNewNodeResult(source *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(source *api.StartNodeRequest) *models.StartNodeParams
	ConvertStopNodeRequest(source *api.StopNodeRequest) *models.StartNodeParams
	ConvertListNodesResult(source *api.ListNodeResponse) *models.ListNodeResult

	//ConvertNodeID(models.NodeID) int
	//RConvertNodeID(int) models.NodeID
}
