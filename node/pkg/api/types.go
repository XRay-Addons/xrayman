package xrayapi

type User struct {
	Name      string `json:"name" validate:"required"`
	VlessUUID string `json:"vlessUuid" validate:"required"`
}

type ClientCfg struct {
	Template       string `json:"template" validate:"required"`
	UserNameField  string `json:"userNameField" validate:"required"`
	VlessUUIDField string `json:"vlessUuidField" validate:"required"`
}

type ServiceStatus string

const (
	ServiceRunning ServiceStatus = "Running"
	ServiceStopped ServiceStatus = "Stopped"
)
