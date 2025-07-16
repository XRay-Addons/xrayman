package clientcfg

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/cfgread"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

type ClientCfg struct {
	clientCfgTemplate models.ClientCfgTemplate
}

func New(
	clientCfgPath string,
	userNameField string,
	vlessUUIDField string,
) (*ClientCfg, error) {
	clientCfg, err := cfgread.ReadJSON(clientCfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: client config file reading: %v", errdefs.ErrConfig, err)
	}
	_, err = template.New("validate").Parse(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid client config template", errdefs.ErrConfig)
	}

	cfgTemplate := models.ClientCfgTemplate{
		ConfigTemplate: clientCfg,
		UserNameField:  userNameField,
		VlessUUIDField: vlessUUIDField,
	}

	if err = validateClientConfig(&cfgTemplate); err != nil {
		return nil, fmt.Errorf("validate client config template: %w", err)
	}

	return &ClientCfg{
		clientCfgTemplate: cfgTemplate,
	}, nil

}

func (cfg *ClientCfg) GetClientConfigTemplate() (*models.ClientCfgTemplate, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: client cfg: get client template", errdefs.ErrNilObjectCall)
	}
	return &cfg.clientCfgTemplate, nil
}

func validateClientConfig(cfg *models.ClientCfgTemplate) error {
	t, err := template.New("json").Parse(cfg.ConfigTemplate)
	if err != nil {
		return fmt.Errorf("%w: template syntax: %v", errdefs.ErrConfig, err)
	}

	testTemplateData := map[string]string{
		cfg.UserNameField:  "field",
		cfg.VlessUUIDField: "name",
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, &testTemplateData); err != nil {
		return fmt.Errorf("%w: template execution: %v", errdefs.ErrConfig, err)
	}

	if !gjson.Valid(buf.String()) {
		return fmt.Errorf("%w: invalid filled json", errdefs.ErrConfig)
	}

	return nil
}
