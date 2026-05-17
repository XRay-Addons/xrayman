package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) NewSubHeader(ctx context.Context, req *api.NewSubHeaderRequest) (*api.Header, error) {
	if h == nil || h.us == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertNewSubHeaderRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.shs.NewHeader(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertHeader(res), nil
}

func (h *Handler) DeleteSubHeader(ctx context.Context, req *api.DeleteSubHeaderRequest) error {
	if h == nil || h.ss == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDeleteSubHeaderRequest(req)
	if err != nil {
		return err
	}
	_, err = h.shs.DeleteHeader(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) ListSubHeaders(ctx context.Context) (*api.ListSubHeadersResponse, error) {
	if h == nil || h.ss == nil {
		return nil, errdefs.NilCall()
	}

	headers, err := h.shs.ListHeaders(ctx, models.ListSubHeadersParams{})
	if err != nil {
		return nil, err
	}

	subResponse := converter.ConvertListSubHeadersResult(headers)
	if err != nil {
		return nil, err
	}

	return subResponse, nil
}
