package spahandler

import (
	"embed"
	"net/http"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"io/fs"
	"strings"
	"fmt"
)

//go:embed dist/*
var distFS embed.FS


func New() (http.Handler, error) {
    distSub, err := fs.Sub(distFS, "dist")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }

	// read index.html
    indexHTML, err := fs.ReadFile(distSub, "index.html")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }

	fileServer := http.FileServer(http.FS(distSub))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(r.URL.Path)
        path := strings.TrimPrefix(r.URL.Path, "/")
        if _, err := fs.Stat(distSub, path); err == nil {
            fileServer.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write(indexHTML)
    }), nil
}
