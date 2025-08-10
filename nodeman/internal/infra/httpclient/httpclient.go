package httpclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Option func(opts *config)

func WithTLS(tlscfg *tls.Config) Option {
	return func(cfg *config) {
		cfg.tlsConfig = tlscfg
	}
}

func New(options ...Option) (*http.Client, error) {
	cfg := config{
		dialerTimeout:       30 * time.Second,
		keepAlive:           30 * time.Second,
		maxIdleConns:        24,
		idleConnTimeout:     90 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
	}

	for _, o := range options {
		o(&cfg)
	}

	transport := &http.Transport{
		TLSClientConfig: cfg.tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   cfg.dialerTimeout,
			KeepAlive: cfg.keepAlive,
		}).DialContext,
		MaxIdleConns:        cfg.maxIdleConns,
		IdleConnTimeout:     cfg.idleConnTimeout,
		TLSHandshakeTimeout: cfg.tlsHandshakeTimeout,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}

type config struct {
	tlsConfig           *tls.Config
	dialerTimeout       time.Duration
	keepAlive           time.Duration
	maxIdleConns        int
	idleConnTimeout     time.Duration
	tlsHandshakeTimeout time.Duration
}
