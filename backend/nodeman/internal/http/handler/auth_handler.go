package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"

	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) Auth(ctx context.Context, req *api.AuthRequest) (
	*api.AuthResponse, error,
) {
	if h == nil || h.auth == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertAuthRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.auth.Auth(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertAuthResult(res), nil
}
