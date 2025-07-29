package httpclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type HttpClientOption func(opts *httpClientConfig)

func WithTLS(crt, key, caCrt string) HttpClientOption {
	return func(cfg *httpClientConfig) {
		cfg.tlsConfig = &tlsConfig{
			crt:   crt,
			key:   key,
			caCrt: caCrt,
		}
	}
}

func New(options ...HttpClientOption) (*http.Client, error) {
	cfg := httpClientConfig{
		dialerTimeout:       30 * time.Second,
		keepAlive:           30 * time.Second,
		maxIdleConns:        24,
		idleConnTimeout:     90 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
	}

	for _, o := range options {
		o(&cfg)
	}

	var tlsConfig *tls.Config
	if cfg.tlsConfig != nil {
		var err error
		tlsConfig, err = createTLSConfig(*cfg.tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("http client init: %w", err)
		}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
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
