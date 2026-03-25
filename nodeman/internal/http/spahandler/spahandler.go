package spahandler

import (
	"embed"
	"net/http"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"io/fs"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

func New() (http.Handler, error) {    
    distSub, err := fs.Sub(distFS, "dist")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }

    indexHTML, err := fs.ReadFile(distSub, "index.html")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }

    fileServer := http.FileServer(http.FS(distSub))

	basePath := "/u"

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
       	
        if path != "" && !strings.HasPrefix(path, "/") {
            http.NotFound(w, r)
            return
        }

		if path == "" {
			http.Redirect(w, r, basePath+"/", http.StatusMovedPermanently)
			return;
		}

        cleanPath := strings.TrimPrefix(path, "/")

        if info, err := fs.Stat(distSub, cleanPath); err == nil && !info.IsDir() {
            r.URL.Path = path
            fileServer.ServeHTTP(w, r)
            return
        }

        // SPA fallback
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write(indexHTML)
    }), nil
}
