package router

import (
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(key string, handlers Handlers, log *zap.Logger) (http.Handler, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: router init: logger", errdefs.ErrNilArgPassed)
	}

	// add middlewares
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logger(log))
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Compression(log))
	r.Use(middleware.Auth([]byte(key), log))
	r.Use(middleware.Encryption([]byte(key), log))

	r.Post("/start", func(w http.ResponseWriter, r *http.Request) {
		handlers.Start(log)(w, r)
	})
	r.Post("/stop", func(w http.ResponseWriter, r *http.Request) {
		handlers.Stop(log)(w, r)
	})
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		handlers.Status(log)(w, r)
	})
	r.Post("/users/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.EditUsers(log)(w, r)
	})

	return r, nil
}
