package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
)

func (h *Handler) NewNode(ctx context.Context, req *ogenserver.NewNodeRequest) (
	*ogenserver.NewNodeResponse, error,
) {
	if h == nil || h.nodes == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertNewNodeRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.nodes.NewNode(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertNewNodeResult(res), nil
}

func (h *Handler) StartNode(ctx context.Context, req *ogenserver.StartNodeRequest) error {
	if h == nil || h.nodes == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertStartNodeRequest(req)
	if err != nil {
		return err
	}
	if err = h.nodes.StartNode(ctx, *p); err != nil {
		return err
	}
	return nil
}

func (h *Handler) StopNode(ctx context.Context, req *ogenserver.StopNodeRequest) error {
	if h == nil || h.nodes == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertStopNodeRequest(req)
	if err != nil {
		return err
	}
	if err = h.nodes.StopNode(ctx, *p); err != nil {
		return err
	}
	return nil
}

func (h *Handler) ListNodes(ctx context.Context) (*ogenserver.ListNodeResponse, error) {
	if h == nil || h.nodes == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.nodes.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertListNodesResult(res), nil
}

func (h *Handler) DeleteNode(ctx context.Context, req *ogenserver.DeleteNodeRequest) error {
	if h == nil || h.nodes == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDeleteNodeRequest(req)
	if err != nil {
		return err
	}
	if err = h.nodes.DeleteNode(ctx, *p); err != nil {
		return err
	}
	return nil
}
