package models

type ServiceStatus int

const (
	unknown ServiceStatus = iota
	ServiceStopped
	ServiceRunning
)
