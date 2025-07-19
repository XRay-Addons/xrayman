package converters

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
)

func StartParamsFromAPI(p api.StartRequest) models.StartParams {
	m := models.StartParams{
		Users: make([]models.User, 0, len(p.Users)),
	}
	for _, u := range p.Users {
		m.Users = append(m.Users, UserFromAPI(u))
	}
	return m
}

func StartResultToAPI(r models.StartResult) api.StartResponse {
	return api.StartResponse{
		ClientCfg: ClientCfgToAPI(r.ClientCfg),
	}
}

func StopParamsFromAPI(p api.StopRequest) models.StopParams {
	return models.StopParams{}
}

func StopResultToAPI(r models.StopResult) api.StopResponse {
	return api.StopResponse{}
}

func StatusParamsFromAPI(p api.StatusRequest) models.StatusParams {
	return models.StatusParams{}
}

func StatusResultToAPI(r models.StatusResult) api.StatusResponse {
	return api.StatusResponse{
		ServiceStatus: ServiceStatusToAPI(r.ServiceStatus),
	}
}

func EditUsersParamsFromAPI(p api.EditUsersRequest) models.EditUsersParams {
	m := models.EditUsersParams{
		Add:    make([]models.User, 0, len(p.Add)),
		Remove: make([]models.User, 0, len(p.Remove)),
	}
	for _, a := range p.Add {
		m.Add = append(m.Add, UserFromAPI(a))
	}
	for _, r := range p.Remove {
		m.Remove = append(m.Remove, UserFromAPI(r))
	}
	return m
}

func EditUsersResultToAPI(r models.EditUsersResult) api.EditUsersResponse {
	return api.EditUsersResponse{}
}
