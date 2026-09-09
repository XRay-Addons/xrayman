package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type options struct {
	tls *tls.Config
	log *zap.Logger
}

type option = func(o *options)

func WithTLS(tls *tls.Config) option {
	return func(o *options) {
		o.tls = tls
	}
}

func WithLog(log *zap.Logger) option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

func New(endpoint string, handler http.Handler, opts ...option) (*HttpServer, error) {
	cfg := options{
		log: zap.NewNop(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	if handler == nil {
		return nil, xerr.NilArg("handler")
	}

	errlog, err := zap.NewStdLogAt(cfg.log, zapcore.ErrorLevel)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &HttpServer{
		server: http.Server{
			Addr:      endpoint,
			Handler:   handler,
			TLSConfig: cfg.tls,
			ErrorLog:  errlog,

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
		return xerr.NilCall()
	}

	var err error
	if s.server.TLSConfig != nil {
		// keys are already in cfg
		err = s.server.ListenAndServeTLS("", "")
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return xerr.WrapWithStack(err)
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
		return xerr.WrapWithStack(err)
	}
	return nil
}
