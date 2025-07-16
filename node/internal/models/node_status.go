package models

type NodeStatus struct {
	Status ServiceStatus `json:"clientCfgTemplate" validate:"required"`
}
g