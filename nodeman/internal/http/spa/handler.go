package spa

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

//go:embed dist/*
var distFS embed.FS

const configPath = "/config.js"

const configTemplate = `
    window.__SPA_URL__ = "{{.SPAUrl}}";
    window.__API_URL__ = "{{.APIUrl}}";
`

func NewHandler(spaPrefix, apiPrefix string) (http.Handler, error) {
    distSub, err := fs.Sub(distFS, "dist")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }

    indexHTML, err := fs.ReadFile(distSub, "index.html")
    if err != nil {
        return nil, errdefs.WrapWithStack(err)
    }
    configContent, err := generateConfigScript(spaPrefix, apiPrefix)
    if err != nil {
        return nil, err
    }

    fileServer := http.FileServer(http.FS(distSub))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path

        // return config.js
        if path == configPath {
            w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
            w.WriteHeader(200)
            w.Write([]byte(configContent))
            return
        }

		// redirect "/spaPrefix" -> "/spaPrefix/"
		if path == "" {
			http.Redirect(w, r, spaPrefix+"/", http.StatusMovedPermanently)
			return;
		}

		// return 404 for "/spaPrefixxxx"
        if !strings.HasPrefix(path, "/") {
            http.NotFound(w, r)
            return
        }

		// return files for file requests
		filePath := strings.TrimPrefix(path, "/")
        if info, err := fs.Stat(distSub, filePath ); err == nil && !info.IsDir() {
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

func generateConfigScript(spaPrefix, apiPrefix string) (string, error) {
    data := map[string]string{
        "SPAUrl": spaPrefix,
        "APIUrl": apiPrefix,
    }
    
    tmpl, err := template.New("config").Parse(configTemplate)
    if err != nil {
        return "", errdefs.WrapWithStack(err)
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", errdefs.WrapWithStack(err)
    }
    
    return buf.String(), nil
}