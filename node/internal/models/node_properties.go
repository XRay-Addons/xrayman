package models

type NodeProperties struct {
	ClientCfgTemplate ClientCfgTemplate `json:"clientCfgTemplate" validate:"required"`
}
