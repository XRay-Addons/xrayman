package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./settings_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type DynamicConfig interface {
	ConvertSettingsRequest(r api.Settings) (*models.Settings, error)
	ConvertSettingsResult(r models.Settings) *api.Settings
}
