package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./dynamic_config_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type DynamicConfig interface {
	ConvertDynamicConfigRequest(r api.DynamicConfig) (*models.DynamicConfig, error)
	ConvertDynamicConfigResult(r models.DynamicConfig) *api.DynamicConfig
}
