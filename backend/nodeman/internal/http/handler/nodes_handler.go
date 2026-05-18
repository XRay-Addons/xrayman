package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) NewNode(ctx context.Context, req *api.NewNodeRequest) (
	*api.NewNodeResponse, error,
) {
	if h == nil || h.ns == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertNewNodeRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.ns.NewNode(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertNewNodeResult(res), nil
}

func (h *Handler) StartNode(ctx context.Context, req *api.StartNodeRequest) error {
	if h == nil || h.ns == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertStartNodeRequest(req)
	if err != nil {
		return err
	}
	_, err = h.ns.StartNode(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) StopNode(ctx context.Context, req *api.StopNodeRequest) error {
	if h == nil || h.ns == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertStopNodeRequest(req)
	if err != nil {
		return err
	}
	_, err = h.ns.StopNode(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) ListNodes(ctx context.Context) (*api.ListNodeResponse, error) {
	if h == nil || h.ns == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.ns.ListNodes(ctx, models.ListNodeParams{})
	if err != nil {
		return nil, err
	}
	return converter.ConvertListNodesResult(res), nil
}

func (h *Handler) DeleteNode(ctx context.Context, req *api.DeleteNodeRequest) error {
	if h == nil || h.ns == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDeleteNodeRequest(req)
	if err != nil {
		return err
	}
	_, err = h.ns.DeleteNode(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}
