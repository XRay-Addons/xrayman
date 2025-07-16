package models

type ClientCfgTemplate struct {
	Template       string `json:"template" validate:"required"`
	UserNameField  string `json:"userNameField" validate:"required"`
	VlessUUIDField string `json:"vlessUuidField" validate:"required"`
}
