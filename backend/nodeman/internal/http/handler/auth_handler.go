package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"

	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func (h *Handler) Auth(ctx context.Context, req *api.AuthRequest) (
	*api.AuthResponse, error,
) {
	if h == nil || h.ns == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertAuthRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.as.Auth(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertAuthResult(res), nil
}
