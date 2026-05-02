package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperr"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func (h *Handler) Auth(ctx context.Context, req *api.AuthRequest) (
	*api.AuthResponse, error,
) {
	if h == nil || h.ns == nil {
		return nil, errdefs.NewNilCall()
	}
	p, err := converter.ConvertAuthRequest(req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	res, err := h.as.Auth(ctx, *p)
	if errors.Is(err, errdefs.ErrAccessDenied) {
		return nil, httperr.ErrAuthToken
	}
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrUnknown
	}
	return converter.ConvertAuthResult(res), nil
}
