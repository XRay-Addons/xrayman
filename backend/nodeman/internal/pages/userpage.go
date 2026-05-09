package pages

import "embed"

//go:embed userpage/**
var userpageFS embed.FS

type UserPageCfg struct {
	ApiPrefix  string `json:"api_prefix"`
	UserPrefix string `json:"user_prefix"`
}

func NewUserPage(apiPrefix, userPrefix string) (*Page, error) {
	cfg := UserPageCfg{
		ApiPrefix:  apiPrefix,
		UserPrefix: userPrefix,
	}
	return new(userpageFS, "userpage", cfg)
}
