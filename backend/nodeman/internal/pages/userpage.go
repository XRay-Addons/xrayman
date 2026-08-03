package pages

import (
	"context"
	"embed"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
)

//go:embed userpage/**
var userpageFS embed.FS

type UserCfgHandler = func(context.Context) (*pagecfg.UserPageCfg, error)

func NewUserPage(cfgHandler UserCfgHandler) (*Page, error) {
	if cfgHandler == nil {
		return nil, xerr.NilArg("cfgHandler")
	}

	pageCfgHandler := func(ctx context.Context) (any, error) {
		cfg, err := cfgHandler(ctx)
		if err != nil {
			return nil, err
		}
		cfgData := converter.ConvertUserPageCfg(*cfg)
		return cfgData, nil
	}

	return new(userpageFS, "userpage", pageCfgHandler)
}
