package secrets

import (
	"fmt"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Option = func(cfg *config)

type config struct {
	issuer string
	expire time.Duration
}

func WithIssuer(issuer string) Option {
	return func(cfg *config) {
		cfg.issuer = issuer
	}
}

func WithExpiration(expire time.Duration) Option {
	return func(cfg *config) {
		cfg.expire = expire
	}
}

type Secrets struct {
	Cert      []byte
	Key       []byte
	AccessKey models.AccessKey
}

const (
	CertFile   = "cert.pem"
	KeyFile    = "key.pem"
	AccessFile = "access.json"
)

// load or create secrets to node access
func Init(dir string, opts ...Option) (*Secrets, error) {
	cfg := &config{
		issuer: "xrayman node",
		expire: 10 * 365 * 24 * time.Hour,
	}
	for _, o := range opts {
		o(cfg)
	}
	cert, key, err := ensureTLS(dir, cfg.issuer, cfg.expire)
	if err != nil {
		return nil, fmt.Errorf("secrets: init: %w", err)
	}
	accessKey, err := ensureAccessKey(dir)
	if err != nil {
		return nil, fmt.Errorf("secrets: init: %w", err)
	}
	return &Secrets{
		Cert:      cert,
		Key:       key,
		AccessKey: *accessKey,
	}, nil
}
