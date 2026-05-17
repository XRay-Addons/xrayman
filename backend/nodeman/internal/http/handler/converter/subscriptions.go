package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./subscriptions_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Subscriptions interface {
	ConvertUserSubRequest(r *api.UserSubParams) (*models.UserSubParams, error)

	ConvertUserSubResultBody(r []models.ClientConfigItem) (api.UserSubContent, error)
}
