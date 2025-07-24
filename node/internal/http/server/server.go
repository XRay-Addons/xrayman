package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

type HttpServer struct {
	server http.Server
}

func New(endpoint string, handler http.Handler, tls *tls.Config) (*HttpServer, error) {
	if handler == nil {
		return nil, fmt.Errorf("http server init: handler: %w", errdefs.ErrNilArgPassed)
	}

	return &HttpServer{
		server: http.Server{
			Addr:      endpoint,
			Handler:   handler,
			TLSConfig: tls,
		},
	}, nil
}

func (s *HttpServer) Listen() error {
	if s == nil {
		return fmt.Errorf("%w: http server", errdefs.ErrNilObjectCall)
	}

	var err error
	if s.server.TLSConfig != nil {
		// keys are already in cfg
		err = s.server.ListenAndServeTLS("", "")
	} else {
		err = s.server.ListenAndServe()
	}
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
