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

func New(h api.Handler, sec api.SecurityHandler, log *zap.Logger) (http.Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("router init: handler: %w", errdefs.ErrNilArgPassed)
	}
	if log == nil {
		return nil, fmt.Errorf("router init: log: %w", errdefs.ErrNilArgPassed)
	}

	// create api handler
	apiHandler, err := api.NewServer(h, sec)
	if err != nil {
		return nil, fmt.Errorf("router: init: %w", err)
	}

	// add middleware from chi
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mw.Logger(log))
	r.Use(chimw.Timeout(10 * time.Second))
	r.Use(chimw.Recoverer)
	r.Use(chimw.NewCompressor(2).Handler)

	// add handler after middlewares
	r.Mount("/", apiHandler)

	return r, nil
}
