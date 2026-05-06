package pages

import (
	"io/fs"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/spa"
	"github.com/go-chi/chi/v5"
)

type Page struct {
	content fs.FS
	config  any
}

func new(contentFS fs.FS, contentDir string, config any) (*Page, error) {
	content, err := fs.Sub(contentFS, contentDir)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}
	return &Page{content: content, config: config}, nil
}

func (p *Page) Mount(r chi.Router, prefix string) error {
	return spa.Mount(r, prefix, p.content, p.config)
}
