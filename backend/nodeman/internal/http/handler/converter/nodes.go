package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./nodes_generated.go
// goverter:extend ConvertAccessKey RConvertAccessKey
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Nodes interface {
	ConvertNewNodeRequest(r *api.NewNodeRequest) (*models.NewNodeParams, error)

	ConvertNewNodeResult(r *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(r *api.StartNodeRequest) (*models.StartNodeParams, error)

	ConvertStopNodeRequest(r *api.StopNodeRequest) (*models.StopNodeParams, error)

	ConvertListNodesResult(r *models.ListNodeResult) *api.ListNodeResponse

	ConvertDeleteNodeRequest(r *api.DeleteNodeRequest) (*models.DeleteNodeParams, error)
}

func ConvertAccessKey(s string) (models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := accessKey.UnmarshalText([]byte(s)); err != nil {
		return accessKey, errdefs.InvalidPayload(err.Error())
	}
	return accessKey, nil
}

func RConvertAccessKey(key models.AccessKey) string {
	return key.String()
}
