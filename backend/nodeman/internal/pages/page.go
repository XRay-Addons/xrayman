package pages

import (
	"io/fs"

	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/spa"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Page struct {
	content        fs.FS
	cfgHandler     spa.CfgHandler
	fallbackFilter func(string) bool
}

var _ router.SPA = (*Page)(nil)

func new(
	contentFS fs.FS,
	contentDir string,
	cfgHandler spa.CfgHandler,
	fallbackFilter func(path string) bool,
) (*Page, error) {
	content, err := fs.Sub(contentFS, contentDir)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &Page{
		content:        content,
		cfgHandler:     cfgHandler,
		fallbackFilter: fallbackFilter,
	}, nil
}

func (p *Page) Mount(
	r chi.Router,
	prefix string,
	log *zap.Logger,
) error {
	return spa.Mount(r, prefix, p.content, p.cfgHandler, p.fallbackFilter, log)
}
