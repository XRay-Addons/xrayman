package xraycfg

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
	cfg models.ClientCfg
}

func NewClientCfg(path string) (*ClientCfg, error) {
	cfgTemplate, err := cfgread.ReadJSON(path)
	if err != nil {
		return nil, fmt.Errorf("%w: client config file reading: %v", errdefs.ErrConfig, err)
	}
	_, err = template.New("validate").Parse(cfgTemplate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid client config template", errdefs.ErrConfig)
	}

	nameField, err := extractNameField(cfgTemplate)
	if err != nil {
		return nil, fmt.Errorf("client cfg init: %w: %v", errdefs.ErrConfig, err)
	}
	vlessUUIdField, err := extractVlessUUIDField(cfgTemplate)
	if err != nil {
		return nil, fmt.Errorf("client cfg init: %w: %v", errdefs.ErrConfig, err)
	}

	clientCfg := models.ClientCfg{
		Template:       cfgTemplate,
		UserNameField:  nameField,
		VlessUUIDField: vlessUUIdField,
	}

	if err = validateClientConfig(&clientCfg); err != nil {
		return nil, fmt.Errorf("validate client config template: %w", err)
	}

	return &ClientCfg{cfg: clientCfg}, nil
}

func validateClientConfig(cfg *models.ClientCfg) error {
	t, err := template.New("json").Parse(cfg.Template)
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

func (cfg *ClientCfg) Get() (*models.ClientCfg, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: client cfg: get", errdefs.ErrNilObjectCall)
	}
	return &cfg.cfg, nil
}
