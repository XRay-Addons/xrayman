package handler

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
type Converter interface {
	ConvertNewNodeRequest(source *api.NewNodeParams) *models.NewNodeParams
	ConvertNewNodeResult(source *models.NewNodeResult) *api.NewNodeResult

	ConvertStartNodeRequest(source *api.StopNodeParams) *models.StartNodeParams
	ConvertStopNodeRequest(source *api.StopNodeParams) *models.StartNodeParams
	ConvertListNodesResult(source *api.ListNodeResult) *models.ListNodeResult
}
