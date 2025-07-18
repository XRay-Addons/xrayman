package router

import (
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	mw "github.com/XRay-Addons/xrayman/node/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(key string, handlers Handlers, log *zap.Logger) (http.Handler, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: router init: logger", errdefs.ErrNilArgPassed)
	}

	// add middlewares
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mw.Logger(log))
	r.Use(chimw.Recoverer)
	r.Use(mw.Auth([]byte(key), log))
	r.Use(chimw.NewCompressor(2).Handler)
	r.Use(mw.Encryption([]byte(key), log))

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
