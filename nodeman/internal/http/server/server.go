package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/rs/cors"
)

type HttpServer struct {
	server http.Server
}

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 10 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultIdleTimeout       = 120 * time.Second
	defaultMaxHeaderBytes    = 1 << 20 // 1 MB
)

func New(endpoint string, handler http.Handler) (*HttpServer, error) {
	if handler == nil {
		return nil, errdefs.NewNilArg("handler")
	}

	// middleware for CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler = c.Handler(handler)

	return &HttpServer{
		server: http.Server{
			Addr:    endpoint,
			Handler: handler,

			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
			MaxHeaderBytes:    defaultMaxHeaderBytes,
		},
	}, nil
}

func (s *HttpServer) Listen() error {
	if s == nil {
		return errdefs.NewNilCall()
	}

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errdefs.WrapWithStack(err)
	}
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if err := s.server.Shutdown(ctx); err == nil {
		return nil
	}
	if err := s.server.Close(); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}
