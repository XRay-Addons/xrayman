package xraycfg

import (
	"text/template"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/cfgread"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type ClientConfig struct {
	cfg models.ClientConfigTemplate
}

func NewClientConfig(path string) (*ClientConfig, error) {
	rawTemplate, err := cfgread.ReadJSON(path)
	if err != nil {
		return nil, err
	}
	_, err = template.New("validate").Parse(rawTemplate)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	cfgTemplate, err := parseClientConfig(rawTemplate)
	if err != nil {
		return nil, err
	}
	emailField, err := extractVlessEmailField(rawTemplate)
	if err != nil {
		return nil, err
	}
	vlessUUIdField, err := extractVlessUUIDField(rawTemplate)
	if err != nil {
		return nil, err
	}

	clientCfg := models.ClientConfigTemplate{
		Template:        cfgTemplate,
		VlessEmailField: emailField,
		VlessUUIDField:  vlessUUIdField,
	}

	return &ClientConfig{cfg: clientCfg}, nil
}

func (cfg *ClientConfig) GetTemplate() (*models.ClientConfigTemplate, error) {
	if cfg == nil {
		return nil, errdefs.NewNilCall()
	}
	return &cfg.cfg, nil
}
