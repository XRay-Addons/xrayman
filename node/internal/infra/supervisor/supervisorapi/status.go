package supervisorapi

type ServiceStatus int

const (
	StatusUnknown ServiceStatus = iota
	StatusStopped
	StatusRunning
)
