package security

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

type Handler struct {
	jwt JWT
}

var _ api.SecurityHandler = (*Handler)(nil)

func New(jwt JWT) (*Handler, error) {
	if jwt == nil {
		return nil, errdefs.NewNilArg("jwt")
	}
	return &Handler{jwt: jwt}, nil
}

func (h *Handler) HandleBearerAuth(ctx context.Context,
	operationName api.OperationName, t api.BearerAuth,
) (context.Context, error) {
	if h == nil || h.jwt == nil {
		return ctx, errdefs.NewNilCall()
	}
	if err := h.jwt.ValidateToken(t.GetToken()); err != nil {
		return ctx, httperr.ErrAuthToken
	}
	return ctx, nil
}
