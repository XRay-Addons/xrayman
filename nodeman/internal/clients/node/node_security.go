package node

import (
	"context"
	"fmt"
	"time"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/golang-jwt/jwt/v4"
)

type NodeSecurity struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

func (s *NodeSecurity) BearerAuth(ctx context.Context,
	op api.OperationName,
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
