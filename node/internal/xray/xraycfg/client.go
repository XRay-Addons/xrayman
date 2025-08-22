package xraycfg

import (
	"bytes"
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
		return nil, err
	}
	_, err = template.New("validate").Parse(cfgTemplate)
	if err != nil {
		return nil, err
	}

	nameField, err := extractNameField(cfgTemplate)
	if err != nil {
		return nil, err
	}
	vlessUUIdField, err := extractVlessUUIDField(cfgTemplate)
	if err != nil {
		return nil, err
	}

	clientCfg := models.ClientCfg{
		Template:       cfgTemplate,
		UserNameField:  nameField,
		VlessUUIDField: vlessUUIdField,
	}

	if err = validateClientConfig(&clientCfg); err != nil {
		return nil, err
	}

	return &ClientCfg{cfg: clientCfg}, nil
}

func validateClientConfig(cfg *models.ClientCfg) error {
	t, err := template.New("json").Parse(cfg.Template)
	if err != nil {
		return errdefs.WrapWithStack(err)
	}

	testTemplateData := map[string]string{
		cfg.UserNameField:  "field",
		cfg.VlessUUIDField: "name",
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, &testTemplateData); err != nil {
		return errdefs.WrapWithStack(err)
	}

	if !gjson.Valid(buf.String()) {
		return errdefs.WrapWithStack(err)
	}

	return nil
}

func (cfg *ClientCfg) Get() (*models.ClientCfg, error) {
	if cfg == nil {
		return nil, errdefs.NewNilCall()
	}
	return &cfg.cfg, nil
}
