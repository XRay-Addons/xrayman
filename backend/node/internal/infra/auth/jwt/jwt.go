package jwt

import (
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/models"

	jwtools "github.com/XRay-Addons/xrayman/common/http/jwt"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
)

type JWT struct {
	secret []byte
	config config
}

var _ (security.JWT) = (*JWT)(nil)

const defaultTTL = 72 * time.Hour
const defaultIssuer = "issuer"

const bearerTokenType = "Bearer"

type config struct {
	ttl    time.Duration
	issuer string
}

type option = func(o *config)

func WithTTL(ttl time.Duration) option {
	return func(o *config) {
		o.ttl = ttl
	}
}

func WithIssuer(issuer string) option {
	return func(o *config) {
		o.issuer = issuer
	}
}

func New(secret models.AccessSecret, opts ...option) (*JWT, error) {
	cfg := config{
		ttl:    defaultTTL,
		issuer: defaultIssuer,
	}
	for _, o := range opts {
		o(&cfg)
	}

	return &JWT{
		secret: secret[:],
		config: cfg,
	}, nil
}

func (j *JWT) ValidateToken(tokenString string) error {
	return jwtools.ValidateToken(tokenString, j.secret, j.config.issuer)
}
