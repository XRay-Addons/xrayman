package security

import (
	"context"
	"fmt"
	"time"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/golang-jwt/jwt/v4"
)

type SecuritySource struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

var _ api.SecuritySource = (*SecuritySource)(nil)

func (s *SecuritySource) BearerAuth(ctx context.Context,
	operationName api.OperationName,
) (api.BearerAuth, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return api.BearerAuth{}, fmt.Errorf("failed to sign JWT: %w", err)
	}

	return api.BearerAuth{
		Token: signed,
	}, nil
}

type SecurityFactory struct {
	issuer     string
	expiration time.Duration
}

type SecurityOption = func(*SecurityFactory)

func WithIssuer(issuer string) SecurityOption {
	return func(f *SecurityFactory) {
		f.issuer = issuer
	}
}

func WithExpiration(expiration time.Duration) SecurityOption {
	return func(f *SecurityFactory) {
		f.expiration = expiration
	}
}

func NewFactory(options ...SecurityOption) *SecurityFactory {
	f := &SecurityFactory{
		issuer:     "node manager",
		expiration: 10 * time.Minute,
	}
	for _, o := range options {
		o(f)
	}
	return f
}

func (f *SecurityFactory) GetSecurity(secret []byte) *SecuritySource {
	return &SecuritySource{
		secret:     secret,
		issuer:     f.issuer,
		expiration: f.expiration,
	}
}
