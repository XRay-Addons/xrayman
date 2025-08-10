package converter

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:extend ConvertNodeID RConvertNodeID ConvertNodeStatusResult
//
//go:generate goverter gen .
type Converter interface {
	ConvertNewNodeRequest(r *api.NewNodeRequest) *models.NewNodeParams
	ConvertNewNodeResult(r *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(r *api.StartNodeRequest) *models.StartNodeParams
	ConvertStopNodeRequest(r *api.StopNodeRequest) *models.StopNodeParams

	ConvertListNodesResult(r *models.ListNodeResult) *api.ListNodeResponse
}

func ConvertNodeID(i models.NodeID) api.NodeID {
	return api.NodeID(i)
}

func RConvertNodeID(i api.NodeID) models.NodeID {
	return models.NodeID(i)
}

func ConvertNodeStatusResult(source models.NodeStatus) api.NodeStatus {
	var response api.NodeStatus
	switch source {
	case models.NodeStatusStopped:
		response = api.NodeStatusStopped
	case models.NodeStatusRunning:
		response = api.NodeStatusRunning
	case models.NodeStatusUnknown:
		response = api.NodeStatusUnknown
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", source))
	}
	return response
}
