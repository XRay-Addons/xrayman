package handler

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

type Handler struct {
	service Service
}

var _ api.Handler = (*Handler)(nil)

func New(s Service) (*Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilArgPassed)
	}
	return &Handler{service: s}, nil
}

func (h *Handler) StartPost(ctx context.Context, req *api.StartRequest) (*api.StartResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p := ConvertStartRequest(req)
	res, err := h.service.Start(ctx, *p)
	if err != nil {
		return nil, httperr.ErrInternalServerError
	}
	return ConvertStartResult(res), nil
}

func (h *Handler) StopPost(ctx context.Context) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	_, err := h.service.Stop(ctx, models.StopParams{})
	if err != nil {
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) GetStatus(ctx context.Context) (*api.StatusResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	status, err := h.service.Status(ctx, models.StatusParams{})
	if err != nil {
		return nil, httperr.ErrInternalServerError
	}
	return ConvertStatusResult(status), nil
}

func (h *Handler) EditUsers(ctx context.Context, req *api.EditUsersRequest) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p := ConvertEditUsersRequest(req)
	_, err := h.service.EditUsers(ctx, *p)
	if err != nil {
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// use passed HttpErr or default unknown
	httpErr := httperr.ErrUnknown
	if e, ok := err.(*httperr.HttpErr); ok {
		httpErr = e
	}
	statusCodeErr := api.ErrorStatusCode(*httpErr)
	return &statusCodeErr
}
