package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err == nil {
		return nil
	}
	if err := s.srv.Close(); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}
