package models

type ServiceStatus string

const (
	ServiceRunning ServiceStatus = "Running"
	ServiceStopped ServiceStatus = "Stopped"
)
