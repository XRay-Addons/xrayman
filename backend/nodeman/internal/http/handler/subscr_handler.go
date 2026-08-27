package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
	"github.com/go-faster/jx"
)

func (h *Handler) UserSub(ctx context.Context, req ogenserver.UserSubParams) (
	[]jx.Raw, error,
) {
	if h == nil || h.subscr == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertUserSubRequest(&req)
	if err != nil {
		return nil, err
	}
	sub, err := h.subscr.GetUserSub(ctx, *p)
	if err != nil {
		return nil, err
	}
	subResponse, err := converter.ConvertUserSubResultBody(sub.ClientConfigs)
	if err != nil {
		return nil, err
	}

	// write to context header with key = "k" and value "v"
	if err := h.writeHeaders(ctx, sub.Headers); err != nil {
		return nil, err
	}

	return subResponse, nil
}
