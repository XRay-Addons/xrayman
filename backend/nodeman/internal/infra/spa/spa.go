package spa

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Mount(r chi.Router, prefix string, content fs.FS, config any, log *zap.Logger) error {
	if r == nil {
		return errdefs.NilArg("r")
	}
	if log == nil {
		return errdefs.NilArg("log")
	}

	// prefix normalisation prefix -> prefix/
	if !strings.HasSuffix(prefix, "/") {
		mountPrefixNormalizer(r, prefix)
		prefix += "/"
	}

	// host config on /config.js
	if err := mountConfig(r, prefix, config, log); err != nil {
		return err
	}

	// host content
	if err := mountContent(r, prefix, content); err != nil {
		return err
	}

	return nil
}

func mountPrefixNormalizer(r chi.Router, prefix string) {
	r.Get(prefix, http.RedirectHandler(prefix+"/",
		http.StatusPermanentRedirect).ServeHTTP)
}

const configPath = "config.json"

func mountConfig(r chi.Router, prefix string, cfg any, log *zap.Logger) error {
	// make config js
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	r.Get(prefix+configPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err := w.Write(cfgData) // we have nothing to do with this error
		if err != nil {
			log.Warn("response writing", zap.String("path", prefix+configPath), zap.Error(err))
		}
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
			return xerr.WrapWithStack(err)
		}
		if p == "." {
			return nil
		}
		items[p] = struct{}{}
		return nil
	}); err != nil {
		return nil, xerr.WrapWithStack(err)
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
		return nil, xerr.WrapWithStack(err)
	}

	// fallback mod time is now
	return &fallbackContent{
		data:    index,
		modTime: time.Now(),
		name:    fallbackFile,
	}, nil
}
