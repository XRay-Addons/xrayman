package tlscfg

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func New(crt, key, caCrt string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	caCert, err := os.ReadFile(caCrt)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
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
