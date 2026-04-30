package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	mw "github.com/XRay-Addons/xrayman/nodeman/internal/http/middleware"
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

	// CORS middleware - разрешаем все
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // разрешаем все домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // разрешаем все заголовки
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false, // при AllowedOrigins: ["*"] должен быть false
		MaxAge:           300,   // кэширование preflight запросов на 5 минут
	}))

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
		chiMount(r, h.path, h.handler)
	}

	return r, nil
}

type handler struct {
	path    string
	handler http.Handler
}

type routerOptions struct {
	handlers       []handler
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
