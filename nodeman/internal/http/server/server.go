package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

type HttpServer struct {
	server http.Server
}

func New(endpoint string, handler http.Handler) (*HttpServer, error) {
	if handler == nil {
		return nil, errdefs.NewNilArg("handler")
	}

	return &HttpServer{
		server: http.Server{
			Addr:    endpoint,
			Handler: handler,
		},
	}, nil
}

func (s *HttpServer) Listen() error {
	if s == nil {
		return errdefs.NewNilCall()
	}

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: run: %w", err)
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
		return fmt.Errorf("http server: shutdown: %w", err)
	}
	return nil
}
