package converter

import (
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./auth_generated.go
// goverter:extend ConvertExpireTime
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type AuthConverter interface {
	ConvertAuthRequest(r *api.AuthRequest) (*models.AuthParams, error)

	ConvertAuthResult(r *models.AuthResult) *api.AuthResponse
}

func ConvertExpireTime(i time.Duration) int {
	return int(i.Seconds())
}
