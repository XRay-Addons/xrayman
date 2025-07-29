package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func createTLSConfig(cfg tlsConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.crt, cfg.key)
	if err != nil {
		return nil, fmt.Errorf("loading server cert: %w", err)
	}

	caCert, err := os.ReadFile(cfg.caCrt)
	if err != nil {
		return nil, fmt.Errorf("reading ca cert: %w", err)
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}, nil
}
