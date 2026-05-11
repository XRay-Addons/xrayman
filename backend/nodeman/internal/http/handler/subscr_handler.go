package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
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

func (h *Handler) ListSubHeaders(ctx context.Context) (*api.ListSubHeadersResponse, error) {
	if h == nil || h.ss == nil {
		return nil, errdefs.NilCall()
	}

	headers, err := h.ss.ListHeaders(ctx, models.ListSubHeadersParams{})
	if err != nil {
		return nil, err
	}

	subResponse := converter.ConvertListSubHeadersResult(headers)
	if err != nil {
		return nil, err
	}

	return subResponse, nil
}

func (h *Handler) NewSubHeader(ctx context.Context, req *api.NewSubHeaderRequest) (*api.Header, error) {
	if h == nil || h.us == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertNewSubHeaderRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.ss.NewHeader(ctx, *p)
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
	_, err = h.ss.DeleteHeader(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}
