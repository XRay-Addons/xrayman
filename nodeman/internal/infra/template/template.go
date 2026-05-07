package template

import (
	"bytes"
	"text/template"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func RenderTemplate(tmpl string, data any) ([]byte, error) {
	t, err := template.New("inline").Parse(tmpl)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, errdefs.WrapWithStack(err)
	}
	return buf.Bytes(), nil
}
