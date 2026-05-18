package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) UserSub(ctx context.Context, req api.UserSubParams) (
	api.UserSubContent, error,
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
