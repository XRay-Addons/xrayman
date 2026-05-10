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

const bearerTokenType = "Bearer"

type config struct {
	issuer *string
}

type option = func(o *config)

func WithIssuerCheck(issuer string) option {
	return func(c *config) {
		c.issuer = &issuer
	}
}

func New(secret models.AccessSecret, opts ...option) (*JWT, error) {
	cfg := config{}
	for _, o := range opts {
		o(&cfg)
	}

	return &JWT{
		secret: secret[:],
		config: cfg,
	}, nil
}

func (j *JWT) ValidateToken(tokenString string) error {
	return jwtools.ValidateToken(tokenString, j.secret,
		jwtools.WithIssuerCheck(j.config.issuer))
}
