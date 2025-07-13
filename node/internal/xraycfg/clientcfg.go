package xraycfg

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

func ParseClientCfg(
	clientCfgPath string,
	userNameField string,
	vlessUUIDField string,
) (*models.ClientConfigTemplate, error) {
	clientCfg, err := readFile(clientCfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: client config file reading: %v", errdefs.ErrConfig, err)
	}
	_, err = template.New("validate").Parse(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid client config template", errdefs.ErrConfig)
	}

	cfgTemplate := models.ClientConfigTemplate{
		ConfigTemplate: clientCfg,
		UserNameField:  userNameField,
		VlessUUIDField: vlessUUIDField,
	}

	if err = validateClientConfig(&cfgTemplate); err != nil {
		return nil, fmt.Errorf("validate client config template: %w", err)
	}

	return &cfgTemplate, nil

}

func validateClientConfig(cfg *models.ClientConfigTemplate) error {
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
