package statichandler

import (
	"embed"
	"net/http"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"io/fs"
	
)

//go:embed dist/*
var distFS embed.FS

func New() (http.Handler, error) {
	distSub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}
	return http.FileServer(http.FS(distSub)), nil
}