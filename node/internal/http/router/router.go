package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	mw "github.com/XRay-Addons/xrayman/node/internal/http/middleware"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/go-chi/chi"
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

func New(h api.Handler, sec api.SecurityHandler, options ...Option) (http.Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("router init: handler: %w", errdefs.ErrNilArgPassed)
	}
	if sec == nil {
		return nil, fmt.Errorf("router init: security: %w", errdefs.ErrNilArgPassed)
	}

	ro := &routerOptions{
		requestTimeout: DefaultRequestTimeout,
		compressionLvl: DefaultCompressionLevel,
		log:            zap.NewNop(),
	}
	for _, o := range options {
		o(ro)
	}

	// create api handler
	apiHandler, err := api.NewServer(h, sec)
	if err != nil {
		return nil, fmt.Errorf("router: init: %w", err)
	}

	// add middleware from chi
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mw.Logger(ro.log))
	r.Use(chimw.Timeout(ro.requestTimeout))
	r.Use(chimw.Recoverer)
	r.Use(chimw.NewCompressor(ro.compressionLvl).Handler)

	// add handler after middlewares
	r.Mount("/", apiHandler)

	return r, nil
}

type routerOptions struct {
	requestTimeout time.Duration
	compressionLvl int
	log            *zap.Logger
}

type Option func(*routerOptions)
