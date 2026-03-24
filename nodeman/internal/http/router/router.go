package router

import (
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	mw "github.com/XRay-Addons/xrayman/nodeman/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const DefaultRequestTimeout = 10 * time.Second
const DefaultCompressionLevel = 2

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

func New(apiHandler, staticHandler http.Handler, options ...Option) (http.Handler, error) {
	if apiHandler == nil {
		return nil, errdefs.NewNilArg("apiHandler")
	}
	if staticHandler == nil {
		return nil, errdefs.NewNilArg("staticHandler")
	}

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
	chiMount(r, "/api", apiHandler)
	chiMount(r, "/u", staticHandler)

	return r, nil
}

type routerOptions struct {
	requestTimeout time.Duration
	compressionLvl int
	log            *zap.Logger
}

type Option func(*routerOptions)

// Golang myass
func chiMount(r chi.Router, prefix string, handler http.Handler) {
	if _, ok := handler.(*chi.Mux); ok {
		r.Mount(prefix, handler)
		return
	}
	r.Mount(prefix, http.StripPrefix(prefix, handler))
}
