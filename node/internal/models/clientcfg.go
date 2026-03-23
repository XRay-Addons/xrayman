package models

import "github.com/go-faster/jx"

type ClientConfigTemplateItem = jx.Raw

type ClientConfigTemplate struct {
	Template        []ClientConfigTemplateItem
	VlessEmailField string
	VlessUUIDField  string
}
