package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
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

func New(endpoint string, handler http.Handler, tls *tls.Config) (*HttpServer, error) {
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
			Addr:      endpoint,
			Handler:   handler,
			TLSConfig: tls,

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

	var err error
	if s.server.TLSConfig != nil {
		// keys are already in cfg
		err = s.server.ListenAndServeTLS("", "")
	} else {
		err = s.server.ListenAndServe()
	}
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
