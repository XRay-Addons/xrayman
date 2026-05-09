package clientcfg

import (
	"os"
	"text/template"

	"github.com/XRay-Addons/xrayman/common/jsonval"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Config struct {
	cfg models.ClientConfigTemplate
}

func New(path string) (cfg *Config, err error) {
	defer func() {
		if err != nil {
			err = errdefs.WrapWithFile(err, path)
		}
	}()

	rawTemplate, err := os.ReadFile(path)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	if err := jsonval.ValidateJsonData(rawTemplate); err != nil {
		return nil, err
	}
	rawTemplateStr := string(rawTemplate)

	_, err = template.New("validate").Parse(rawTemplateStr)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	cfgTemplate, err := parseClientConfig(rawTemplateStr)
	if err != nil {
		return nil, err
	}
	emailField, err := extractVlessEmailField(rawTemplateStr)
	if err != nil {
		return nil, err
	}
	vlessUUIdField, err := extractVlessUUIDField(rawTemplateStr)
	if err != nil {
		return nil, err
	}

	clientCfg := models.ClientConfigTemplate{
		Template:        cfgTemplate,
		VlessEmailField: emailField,
		VlessUUIDField:  vlessUUIdField,
	}

	return &Config{cfg: clientCfg}, nil
}

func (cfg *Config) GetTemplate() (*models.ClientConfigTemplate, error) {
	if cfg == nil {
		return nil, errdefs.NilCall()
	}
	return &cfg.cfg, nil
}
