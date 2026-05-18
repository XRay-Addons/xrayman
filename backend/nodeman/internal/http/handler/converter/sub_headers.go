package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./sub_headers_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type SubHeaders interface {
	ConvertNewSubHeaderRequest(r *api.NewSubHeaderRequest) (*models.NewSubHeaderParams, error)

	ConvertDeleteSubHeaderRequest(r *api.DeleteSubHeaderRequest) (*models.DeleteSubHeaderParams, error)

	ConvertListSubHeadersResult(r *models.ListSubHeadersResult) *api.ListSubHeadersResponse

	ConvertHeader(r *models.Header) *api.Header
}
