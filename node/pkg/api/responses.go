package xrayapi

type StartResponse struct {
	ClientCfg ClientCfg `json:"clientCfg" validate:"required"`
}

type StopResponse struct {
}

type StatusResponse struct {
	ServiceStatus ServiceStatus `json:"serviceStatus" validate:"required"`
}

type EditUsersResponse struct {
}
