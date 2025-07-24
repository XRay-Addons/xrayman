package handler

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	log     *zap.Logger
}

var _ api.Handler = (*Handler)(nil)

func New(s Service, log *zap.Logger) (*Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("handler init: service: %w", errdefs.ErrNilArgPassed)
	}
	if log == nil {
		return nil, fmt.Errorf("handler init: logger: %w", errdefs.ErrNilArgPassed)
	}
	return &Handler{
		service: s,
		log:     log,
	}, nil
}

func (h *Handler) StartPost(ctx context.Context, req *api.StartRequest) (_ *api.StartResponse, err error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}

	p := ConvertStartRequest(req)
	res, err := h.service.Start(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
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
		h.logError(ctx, err)
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
		h.logError(ctx, err)
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
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// use passed HttpErr or default unknown
	httpErr := httperr.ErrUnknown
	if e, ok := err.(*httperr.HttpErr); ok {
		httpErr = e
	} else {
		// all errors pass to this handler, many of them are consequences
		// of errors processed and logged before, others come here
		h.logError(ctx, err)
	}
	statusCodeErr := api.ErrorStatusCode(*httpErr)
	return &statusCodeErr
}

func (h *Handler) logError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	h.log.Error("handler request",
		zap.String(constants.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
}
