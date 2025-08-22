package tlscfg

import (
	"crypto/tls"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Load(nodeCrt, nodeKey []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(nodeCrt, nodeKey)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
