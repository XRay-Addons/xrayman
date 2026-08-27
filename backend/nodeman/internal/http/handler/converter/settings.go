package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./settings_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type DynamicConfig interface {
	ConvertSettingsRequest(r ogenserver.Settings) (*models.Settings, error)
	ConvertSettingsResult(r models.Settings) *ogenserver.Settings
}
