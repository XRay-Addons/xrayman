package models

type ClientConfigTemplateItem = string

type ClientConfigTemplate struct {
	Template        []ClientConfigTemplateItem
	VlessEmailField string
	VlessUUIDField  string
}
