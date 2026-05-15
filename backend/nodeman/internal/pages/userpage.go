package pages

import "embed"

//go:embed userpage/**
var userpageFS embed.FS

type UserPageCfg struct {
	ApiPrefix  string `json:"api_prefix"`
	UserPrefix string `json:"user_prefix"`
}

func NewUserPage(apiServiceUrl, userSpaUrl string) (*Page, error) {
	cfg := UserPageCfg{
		ApiPrefix:  apiServiceUrl,
		UserPrefix: userSpaUrl,
	}
	return new(userpageFS, "userpage", cfg)
}
