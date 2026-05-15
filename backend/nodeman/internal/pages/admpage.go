package pages

import "embed"

//go:embed admpage/**
var admpageFS embed.FS

type AdmPageCfg struct {
	ApiPrefix   string `json:"api_prefix"`
	AdminPrefix string `json:"admin_prefix"`
	UserPrefix  string `json:"user_prefix"`
}

func NewAdmPage(apiServiceUrl, adminSpaUrl, userSpaUrl string) (*Page, error) {
	cfg := AdmPageCfg{
		ApiPrefix:   apiServiceUrl,
		AdminPrefix: adminSpaUrl,
		UserPrefix:  userSpaUrl,
	}
	return new(admpageFS, "admpage", cfg)
}
