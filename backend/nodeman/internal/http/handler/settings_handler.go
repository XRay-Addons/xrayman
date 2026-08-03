package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) GetSettings(ctx context.Context) (*api.Settings, error) {
	if h == nil || h.settings == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.settings.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertSettingsResult(*res), nil
}

func (h *Handler) SetSettings(ctx context.Context, req *api.Settings) error {
	if h == nil || h.settings == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertSettingsRequest(*req)
	if err != nil {
		return err
	}
	return h.settings.SetSettings(ctx, *p)
}
