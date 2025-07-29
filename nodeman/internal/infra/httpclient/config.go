package httpclient

import "time"

type tlsConfig struct {
	crt   string
	key   string
	caCrt string
}

type httpClientConfig struct {
	tlsConfig           *tlsConfig
	dialerTimeout       time.Duration
	keepAlive           time.Duration
	maxIdleConns        int
	idleConnTimeout     time.Duration
	tlsHandshakeTimeout time.Duration
}
