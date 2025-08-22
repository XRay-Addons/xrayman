package node

import (
	"context"
	"time"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

type NodeSecurity struct {
	secret     models.AccessSecret
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
	// secret is [32]byte, but token require []byte
	signed, err := token.SignedString(s.secret[:])
	if err != nil {
		return api.BearerAuth{}, errdefs.WrapWithStack(err)
	}

	return api.BearerAuth{
		Token: signed,
	}, nil
}
