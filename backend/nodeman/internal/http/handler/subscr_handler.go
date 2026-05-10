package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func (h *Handler) UserSub(ctx context.Context, req api.UserSubParams) (
	*api.UserSubResponseHeaders, error,
) {
	if h == nil || h.ss == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertUserSubRequest(&req)
	if err != nil {
		return nil, err
	}
	sub, exists, err := h.ss.GetUserSub(ctx, *p)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errdefs.NotFound("user")
	}
	subResponse, err := converter.ConvertUserSubResult(sub)
	if err != nil {
		return nil, err
	}
	return subResponse, nil
}
