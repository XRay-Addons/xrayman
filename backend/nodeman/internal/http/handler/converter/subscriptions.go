package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	jx "github.com/go-faster/jx"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./subscriptions_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Subscriptions interface {
	ConvertUserSubRequest(r *ogenserver.UserSubParams) (*models.UserSubParams, error)

	ConvertUserSubResultBody(r []models.ClientConfigItem) ([]jx.Raw, error)
}
