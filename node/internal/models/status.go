package models

type Status string

const (
	NotRunning Status = "NotRunning"
	Running    Status = "Running"
)

type NodeStatus struct {
	Status Status `json:"status"`
	CPULoad int `json:"cpu_load"`
	RAMLoad int `json:"ram_load"`
}