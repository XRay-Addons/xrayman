package tlscfg

import (
	"crypto/tls"
	"fmt"
)

func Load(nodeCrt, nodeKey string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(nodeCrt, nodeKey)
	if err != nil {
		return nil, fmt.Errorf("loading server cert: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
