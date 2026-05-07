package jwt

import (
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	jwtools "github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/http/jwt"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
)

type JWT struct {
	secret []byte
	config config
}

var _ (auth.JWT) = (*JWT)(nil)
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

func New(secret string, opts ...option) (*JWT, error) {
	if secret == "" {
		return nil, errdefs.NilArg("secret")
	}
	cfg := config{
		ttl:    defaultTTL,
		issuer: defaultIssuer,
	}
	for _, o := range opts {
		o(&cfg)
	}

	return &JWT{
		secret: []byte(secret),
		config: cfg,
	}, nil
}

func (j *JWT) GenerateToken(subject string) (*models.AuthResult, error) {
	token, err := jwtools.GenerateToken(j.secret, j.config.issuer,
		jwtools.WithSubject(subject))
	if err != nil {
		return nil, err
	}
	return &models.AuthResult{
		AccessToken: token,
		TokenType:   bearerTokenType,
		ExpiresIn:   j.config.ttl,
	}, nil
}

func (j *JWT) ValidateToken(tokenString string) error {
	return jwtools.ValidateToken(tokenString, j.secret, j.config.issuer)
}
