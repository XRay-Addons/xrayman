package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"

	mw "github.com/XRay-Addons/xrayman/common/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

type SPA interface {
	Mount(r chi.Router, prefix string, log *zap.Logger) error
}

func WithSPA(path string, spa SPA) Option {
	return func(r *routerOptions) {
		r.spas = append(r.spas,
			spaItem{
				path: path,
				page: spa,
			},
		)
	}
}

func WithCrossOrigin(origins []string) Option {
	return func(r *routerOptions) {
		r.origins = append(r.origins, origins...)
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
	r.Use(mw.Headers())
	r.Use(chimw.Timeout(ro.requestTimeout))
	r.Use(chimw.Recoverer)
	r.Use(chimw.NewCompressor(ro.compressionLvl).Handler)

	// allow cross-origin
	if len(ro.origins) != 0 {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   ro.origins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// add handler after middlewares
	for _, h := range ro.handlers {
		if h.handler == nil {
			return nil, xerr.NilArg(fmt.Sprintf("%s handler", h.path))
		}
		chiMountHandler(r, h.path, h.handler)
	}

	// add SPAs after middlewares
	for _, spa := range ro.spas {
		if spa.page == nil {
			return nil, xerr.NilArg(fmt.Sprintf("%s spa", spa.path))
		}
		if err := spa.page.Mount(r, spa.path, ro.log); err != nil {
			return nil, err
		}
	}

	return r, nil
}

type handler struct {
	path    string
	handler http.Handler
}

type spaItem struct {
	path string
	page SPA
}

type routerOptions struct {
	handlers       []handler
	spas           []spaItem
	origins        []string
	requestTimeout time.Duration
	compressionLvl int
	log            *zap.Logger
}

type Option func(*routerOptions)

// Golang myass
func chiMountHandler(r chi.Router, prefix string, handler http.Handler) {
	if _, ok := handler.(*chi.Mux); ok || prefix == "/" {
		r.Mount(prefix, handler)
		return
	}
	r.Mount(prefix, http.StripPrefix(prefix, handler))
}
