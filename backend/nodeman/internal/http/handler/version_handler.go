package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func (h *Handler) GetVersion(ctx context.Context) (string, error) {
	if h == nil || h.version == nil {
		return "", errdefs.NilCall()
	}

	res, err := h.version.GetVersion(ctx)
	if err != nil {
		return "", err
	}
	return res, nil
}
