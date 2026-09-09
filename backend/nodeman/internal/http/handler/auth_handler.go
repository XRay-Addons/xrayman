package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"

	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter/convauth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
)

func (h *Handler) Auth(ctx context.Context, req *ogenserver.AuthRequest) (
	*ogenserver.AuthResponse, error,
) {
	if h == nil || h.auth == nil {
		return nil, errdefs.NilCall()
	}
	p, err := convauth.ConvertAuthRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.auth.Auth(ctx, *p)
	if err != nil {
		return nil, err
	}
	return convauth.ConvertAuthResult(res), nil
}
