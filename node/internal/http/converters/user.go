package converters

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
)

func UserFromAPI(u api.User) models.User {
	return models.User{
		Name:      u.Name,
		VlessUUID: u.VlessUUID,
	}
}

func UserToAPI(u models.User) api.User {
	return api.User{
		Name:      u.Name,
		VlessUUID: u.VlessUUID,
	}
}
