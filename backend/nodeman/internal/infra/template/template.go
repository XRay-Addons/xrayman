package template

import (
	"bytes"
	"text/template"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func RenderTemplate(tmpl string, data any) ([]byte, error) {
	t, err := template.New("inline").Parse(tmpl)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return buf.Bytes(), nil
}
