package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
)

func (h *Handler) GetSettings(ctx context.Context) (*ogenserver.Settings, error) {
	if h == nil || h.settings == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.settings.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertSettingsResult(*res), nil
}

func (h *Handler) SetSettings(ctx context.Context, req *ogenserver.Settings) error {
	if h == nil || h.settings == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertSettingsRequest(*req)
	if err != nil {
		return err
	}
	return h.settings.SetSettings(ctx, *p)
}
