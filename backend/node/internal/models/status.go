package models

type ServiceStatus int

const (
	ServiceStatusUnknown ServiceStatus = iota + 1
	ServiceStatusStopped
	ServiceStatusRunning
)
