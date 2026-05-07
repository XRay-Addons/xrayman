package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/node/internal/infra/common/http/middleware"
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
		return nil, errdefs.NewNilArg("s")
	}
	if log == nil {
		return nil, errdefs.NewNilArg("log")
	}
	return &Handler{
		service: s,
		log:     log,
	}, nil
}

func (h *Handler) Start(ctx context.Context, req *api.StartRequest) (_ *api.StartResponse, err error) {
	if h == nil || h.service == nil {
		return nil, errdefs.NewNilCall()
	}

	p := converter.ConvertStartRequest(req)
	res, err := h.service.Start(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertStartResult(res), nil
}

func (h *Handler) Stop(ctx context.Context) error {
	if h == nil || h.service == nil {
		return errdefs.NewNilCall()
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
		return nil, errdefs.NewNilCall()
	}
	status, err := h.service.Status(ctx, models.StatusParams{})
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertStatusResult(status), nil
}

func (h *Handler) EditUsers(ctx context.Context, req *api.EditUsersRequest) error {
	if h == nil || h.service == nil {
		return errdefs.NewNilCall()
	}
	p := converter.ConvertEditUsersRequest(req)
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
	if ok := errors.As(err, &httpErr); !ok {
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
		zap.String(middleware.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
}
