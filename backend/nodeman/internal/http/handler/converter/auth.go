package converter

import (
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./auth_generated.go
// goverter:extend ConvertExpireTime
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type AuthConverter interface {
	ConvertAuthRequest(r *ogenserver.AuthRequest) (*models.AuthParams, error)

	ConvertAuthResult(r *models.AuthResult) *ogenserver.AuthResponse
}

func ConvertExpireTime(i time.Duration) int {
	return int(i.Seconds())
}
