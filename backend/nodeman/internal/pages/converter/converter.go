package converter

import (
	cfgs "github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/schemas"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./schemes_generated.go
//
//go:generate goverter gen .
type PageCfgConverter interface {
	ConvertUserPageCfg(r cfgs.UserPageCfg) schemas.UserpagecfgJson
	ConvertAdminPageCfg(r cfgs.AdminPageCfg) schemas.AdminpagecfgJson
}
