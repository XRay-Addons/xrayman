package spa

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/go-chi/chi/v5"
)

func Mount(r chi.Router, prefix string, content fs.FS, config any) error {
	if r == nil {
		return errdefs.NewNilArg("r")
	}

	// prefix normalisation prefix -> prefix/
	if !strings.HasSuffix(prefix, "/") {
		mountPrefixNormalizer(r, prefix, content)
		prefix += "/"
	}

	// host config on /config.js
	if err := mountConfig(r, prefix, config); err != nil {
		return err
	}

	// host content
	if err := mountContent(r, prefix, content); err != nil {
		return err
	}

	return nil
}

func mountPrefixNormalizer(r chi.Router, prefix string, content fs.FS) {
	r.Get(prefix, http.RedirectHandler(prefix+"/",
		http.StatusPermanentRedirect).ServeHTTP)
}

const configPath = "config.js"

func mountConfig(r chi.Router, prefix string, cfg any) error {
	// make config js
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		return errdefs.WrapWithStack(err)
	}
	cfgJs := []byte("window.__CONFIG__ = ")
	cfgJs = append(cfgJs, cfgData...)
	cfgJs = append(cfgJs, ';')

	r.Get(prefix+configPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(cfgJs)
	})

	return nil
}

func mountContent(r chi.Router, prefix string, content fs.FS) error {
	contentItems, err := listContent(content)
	if err != nil {
		return err
	}
	contentFallback, err := getFallback(content)
	if err != nil {
		return err
	}
	contentFS := http.FileServer(http.FS(content))

	r.Get(prefix+"*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, prefix)

		// serve files
		if _, ok := contentItems[path]; ok {
			http.StripPrefix(prefix, contentFS).ServeHTTP(w, r)
			return
		}

		// serve fallback
		http.ServeContent(w, r,
			contentFallback.name,
			contentFallback.modTime,
			bytes.NewReader(contentFallback.data),
		)
	})

	return nil
}

func listContent(content fs.FS) (map[string]struct{}, error) {
	items := make(map[string]struct{})
	if err := fs.WalkDir(content, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return errdefs.WrapWithStack(err)
		}
		if p == "." {
			return nil
		}
		items[p] = struct{}{}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

type fallbackContent struct {
	data    []byte
	modTime time.Time
	name    string
}

const fallbackFile = "index.html"

func getFallback(content fs.FS) (*fallbackContent, error) {
	// read fallback content
	index, err := fs.ReadFile(content, fallbackFile)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	// fallback mod time is now
	return &fallbackContent{
		data:    index,
		modTime: time.Now(),
		name:    fallbackFile,
	}, nil
}
