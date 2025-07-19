package converters

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
)

func ClientCfgFromAPI(cfg api.ClientCfg) models.ClientCfg {
	return models.ClientCfg{
		Template:       cfg.Template,
		UserNameField:  cfg.UserNameField,
		VlessUUIDField: cfg.VlessUUIDField,
	}
}

func ClientCfgToAPI(cfg models.ClientCfg) api.ClientCfg {
	return api.ClientCfg{
		Template:       cfg.Template,
		UserNameField:  cfg.UserNameField,
		VlessUUIDField: cfg.VlessUUIDField,
	}
}
