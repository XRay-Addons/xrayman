package handlers

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	apimodels "github.com/XRay-Addons/xrayman/node/pkg/api/models"
	apirequests "github.com/XRay-Addons/xrayman/node/pkg/api/requests"
)

func fromStartRequest(req apirequests.StartRequest) []models.User {
	users := make([]models.User, 0, len(req.Users))
	for _, ru := range req.Users {
		users = append(users, models.User{
			Name:      ru.Name,
			VlessUUID: ru.VlessUUID,
		})
	}
	return users
}

func toStartResponse(nodeProps models.NodeProperties) apirequests.StartResponse {
	clientCfg := nodeProps.ClientCfgTemplate

	return apirequests.StartResponse{
		NodeConfig: apimodels.NodeConfig{
			UserConfigTemplate: apimodels.ClientCfgTemplate{
				Template:       clientCfg.Template,
				UserNameField:  clientCfg.UserNameField,
				VlessUUIDField: clientCfg.VlessUUIDField,
			},
		},
	}
}
