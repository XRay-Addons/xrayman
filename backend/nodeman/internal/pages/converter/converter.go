package converter

import (
	cfgs "github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
	schemas "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/schemas-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./schemes_generated.go
//
//go:generate goverter gen .
type AuthConverter interface {
	ConvertUserPageCfg(r *cfgs.UserPageCfg) *schemas.UserpagecfgJson
	ConvertAdminPageCfg(r *cfgs.AdminPageCfg) *schemas.AdminpagecfgJson
}
