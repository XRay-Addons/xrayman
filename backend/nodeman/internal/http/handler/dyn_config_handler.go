package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) GetDynamicConfig(ctx context.Context) (*api.DynamicConfig, error) {
	if h == nil || h.ns == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.dcs.GetDynamicConfig(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertDynamicConfigResult(*res), nil
}

func (h *Handler) SetDynamicConfig(ctx context.Context, req *api.DynamicConfig) error {
	if h == nil || h.ns == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDynamicConfigRequest(*req)
	if err != nil {
		return err
	}
	return h.dcs.SetDynamicConfig(ctx, *p)
}
