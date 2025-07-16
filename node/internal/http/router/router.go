package router

import (
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handlers"
	"github.com/XRay-Addons/xrayman/node/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(log *zap.Logger) (http.Handler, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: router init: logger", errdefs.ErrNilArgPassed)
	}

	// add middleware
	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logger(log))
	r.Use(middleware.Compression(log))

	/*// register routes
	if err := addStatusHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("status handlers: %v", err)
	}
	if err := addUsersHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("users handlers: %v", err)
	}*/

	return r, nil
}

func addStatusHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	statusHandler, err := handlers.NewStatusHandler(s, log)
	if err != nil {
		return fmt.Errorf("status handler creation: %v", err)
	}
	r.Route("/node", func(r chi.Router) {
		r.Post("/start", statusHandler.StartHandler())
		r.Post("/stop", statusHandler.StopHandler())
		r.Get("/status", statusHandler.StatusHandler())
	})
	return nil
}

func addUsersHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	usersHandler, err := handlers.NewUsersHandler(s, log)
	if err != nil {
		return fmt.Errorf("users handler creation: %v", err)
	}
	r.Route("/users", func(r chi.Router) {
		r.Post("/add", usersHandler.AddUsersHandler())
		r.Post("/del", usersHandler.DelUsersHandler())
	})

	return nil
}
