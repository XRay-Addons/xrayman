package pages

import (
	"embed"

	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
)

//go:embed admpage/**
var admpageFS embed.FS

type AdmPageCfg struct {
	ApiPrefix   string `json:"api_prefix"`
	AdminPrefix string `json:"admin_prefix"`
	UserPrefix  string `json:"user_prefix"`
}

func NewAdmPage(cfg pagecfg.AdminPageCfg) (*Page, error) {
	adminPageCfg := converter.ConvertAdminPageCfg(&cfg)
	return new(admpageFS, "admpage", adminPageCfg)
}
