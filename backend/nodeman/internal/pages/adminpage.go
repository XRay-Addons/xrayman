package pages

import (
	"context"
	"embed"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
)

//go:embed admpage/**
var admpageFS embed.FS

type AdminCfgHandler = func(context.Context) (*pagecfg.AdminPageCfg, error)

func NewAdmPage(cfgHandler AdminCfgHandler) (*Page, error) {
	if cfgHandler == nil {
		return nil, xerr.NilArg("cfgHandler")
	}

	pageCfgHandler := func(ctx context.Context) (any, error) {
		cfg, err := cfgHandler(ctx)
		if err != nil {
			return nil, err
		}
		cfgData := converter.ConvertAdminPageCfg(*cfg)
		return cfgData, nil
	}
	pageFallbackFilter := func(path string) bool {
		return path == ""
	}

	return new(admpageFS, "admpage", pageCfgHandler, pageFallbackFilter)
}
