package pages

import (
	"embed"

	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pages/pagecfg"
)

//go:embed userpage/**
var userpageFS embed.FS

func NewUserPage(cfg pagecfg.UserPageCfg) (*Page, error) {
	userPageCfg := converter.ConvertUserPageCfg(&cfg)
	return new(userpageFS, "userpage", userPageCfg)
}
