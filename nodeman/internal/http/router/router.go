package router

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/spa"

	mw "github.com/XRay-Addons/xrayman/nodeman/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const DefaultRequestTimeout = 10 * time.Second
const DefaultCompressionLevel = 2

func WithHandler(path string, h http.Handler) Option {
	return func(r *routerOptions) {
		r.handlers = append(r.handlers,
			handler{
				path:    path,
				handler: h,
			},
		)
	}
}

func WithSPA(path string, content fs.FS, cfg any) Option {
	return func(r *routerOptions) {
		r.spas = append(r.spas,
			spaContent{
				path:    path,
				content: content,
				cfg:     cfg,
			},
		)
	}
}

func WithTimeout(d time.Duration) Option {
	return func(r *routerOptions) {
		r.requestTimeout = d
	}
}

func WithCompressionLevel(level int) Option {
	return func(r *routerOptions) {
		r.compressionLvl = level
	}
}

func WithLogger(log *zap.Logger) Option {
	return func(r *routerOptions) {
		if log != nil {
			r.log = log
		}
	}
}

func New(options ...Option) (http.Handler, error) {
	ro := &routerOptions{
		requestTimeout: DefaultRequestTimeout,
		compressionLvl: DefaultCompressionLevel,
		log:            zap.NewNop(),
	}
	for _, o := range options {
		o(ro)
	}

	// add middleware from chi
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mw.Logger(ro.log))
	r.Use(chimw.Timeout(ro.requestTimeout))
	r.Use(chimw.Recoverer)
	r.Use(chimw.NewCompressor(ro.compressionLvl).Handler)

	// add handler after middlewares
	for _, h := range ro.handlers {
		if h.handler == nil {
			return nil, errdefs.NewNilArg(fmt.Sprintf("%s handler", h.path))
		}
		chiMountHandler(r, h.path, h.handler)
	}

	// add SPAs after middlewares
	for _, spa := range ro.spas {
		if spa.content == nil {
			return nil, errdefs.NewNilArg(fmt.Sprintf("%s spa", spa.path))
		}
		if err := chiMountSPA(r, spa.path, spa.content, spa.cfg); err != nil {
			return nil, err
		}
	}

	return r, nil
}

type handler struct {
	path    string
	handler http.Handler
}

type spaContent struct {
	path    string
	content fs.FS
	cfg     any
}

type routerOptions struct {
	handlers       []handler
	spas           []spaContent
	requestTimeout time.Duration
	compressionLvl int
	log            *zap.Logger
}

type Option func(*routerOptions)

// Golang myass
func chiMountHandler(r chi.Router, prefix string, handler http.Handler) {
	if _, ok := handler.(*chi.Mux); ok {
		r.Mount(prefix, handler)
		return
	}
	r.Mount(prefix, http.StripPrefix(prefix, handler))
}

const configPath = "config.js"

func chiMountSPA(r chi.Router, prefix string, content fs.FS, cfg any) error {
	return spa.Mount(r, prefix, content, cfg)
}

/*// ensure prefix ends with "/"
	if !strings.HasSuffix(prefix, "/") {
		r.Get(prefix, http.RedirectHandler(prefix+"/",
			http.StatusPermanentRedirect).ServeHTTP)
		prefix += "/"
	}

	// host config on /config.js
	cfgData, err := getConfigData(cfg)
	if err != nil {
		return errdefs.WrapWithStack(err)
	}
	r.Get(prefix+configPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(cfgData)
	})

	// get list of embed fs files
	paths := make(map[string]struct{})
	err = fs.WalkDir(content, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return errdefs.WrapWithStack(err)
		}
		if p == "." {
			return nil
		}
		paths[p] = struct{}{}
		return nil
	})
	if err != nil {
		return err
	}

	// read index.html
	index, err := fs.ReadFile(content, "index.html")
	if err != nil {
		return errdefs.WrapWithStack(err)
	}

	var modTime time.Time
	if f, err := content.Open("index.html"); err == nil {
		if st, err := f.Stat(); err == nil {
			modTime = st.ModTime()
		}
		f.Close()
	}

	fileServer := http.FileServer(http.FS(content))
	r.Get(prefix+"*", func(w http.ResponseWriter, r *http.Request) {
		// fallback to
		path := strings.TrimPrefix(r.URL.Path, prefix)
		if path == "" {
			path = "index.html"
		}

		// file if exists
		if _, ok := paths[path]; ok {
			http.StripPrefix(prefix, fileServer).ServeHTTP(w, r)
			return
		}

		// serve index.html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(
			w,
			r,
			path,
			modTime,
			bytes.NewReader(index),
		)
	})

	return nil
}

func getConfigData(cfg any) ([]byte, error) {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	cfgData := []byte("window.__CONFIG__ = ")
	cfgData = append(cfgData, cfgJson...)
	cfgData = append(cfgData, ';')

	return cfgData, nil
}*/
