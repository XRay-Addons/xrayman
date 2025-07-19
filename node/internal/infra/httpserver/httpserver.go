package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	srv    *http.Server
	cancel context.CancelFunc
}

func New(addr string, handler http.Handler) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
		cancel: cancel,
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
	// pass cancel to all handlers
	s.cancel()

	// shutdown server
	if err := s.srv.Shutdown(ctx); err == nil {
		return nil
	}

	if err := s.srv.Close(); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}
