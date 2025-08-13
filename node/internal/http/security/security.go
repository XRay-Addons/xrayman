package security

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/golang-jwt/jwt"
)

type Handler struct {
	secret models.AccessSecret
}

var _ api.SecurityHandler = (*Handler)(nil)

func New(secret models.AccessSecret) *Handler {
	return &Handler{secret: secret}
}

func (s *Handler) HandleBearerAuth(ctx context.Context,
	operationName api.OperationName, t api.BearerAuth,
) (context.Context, error) {
	parsedToken, err := jwt.Parse(
		t.GetToken(),
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return s.secret[:], nil
		},
	)

	if err != nil || !parsedToken.Valid {
		return nil, httperr.ErrAuthToken
	}

	return ctx, nil
}
