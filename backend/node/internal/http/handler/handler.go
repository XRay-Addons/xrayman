package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/common/http/httperr"
	"github.com/XRay-Addons/xrayman/common/http/middleware"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperrdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	log     *zap.Logger
}

func WithLogger(log *zap.Logger) option {
	return func(h *Handler) {
		if log == nil {
			return
		}
		h.log = log
	}
}

type option = func(h *Handler)

var _ api.Handler = (*Handler)(nil)

func New(s Service, opts ...option) (*Handler, error) {
	if s == nil {
		return nil, errdefs.NilArg("s")
	}
	handler := &Handler{
		service: s,
		log:     zap.NewNop(),
	}
	for _, o := range opts {
		o(handler)
	}
	return handler, nil
}

func (h *Handler) Start(ctx context.Context, req *api.StartRequest) (_ *api.StartResponse, err error) {
	if h == nil || h.service == nil {
		return nil, errdefs.NilCall()
	}

	p := converter.ConvertStartRequest(req)
	res, err := h.service.Start(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertStartResult(res), nil
}

func (h *Handler) Stop(ctx context.Context) error {
	if h == nil || h.service == nil {
		return errdefs.NilCall()
	}
	_, err := h.service.Stop(ctx, models.StopParams{})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetStatus(ctx context.Context) (*api.StatusResponse, error) {
	if h == nil || h.service == nil {
		return nil, errdefs.NilCall()
	}
	status, err := h.service.Status(ctx, models.StatusParams{})
	if err != nil {
		return nil, err
	}
	return converter.ConvertStatusResult(status), nil
}

func (h *Handler) EditUsers(ctx context.Context, req *api.EditUsersRequest) error {
	if h == nil || h.service == nil {
		return errdefs.NilCall()
	}
	p := converter.ConvertEditUsersRequest(req)
	_, err := h.service.EditUsers(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// if err = error + status, return status, log error,
	// elsewise analyze errors and map to status codes
	if e, s := httperr.ExtractStatus[api.ErrorStatusCode](err); s != nil {
		h.logError(ctx, e)
		return s
	}
	// this is our error, translate it to status
	s := h.translateError(err)
	h.logError(ctx, err)
	return s
}

func (h *Handler) translateError(err error) *api.ErrorStatusCode {
	if err == nil {
		return nil
	}
	if errors.Is(err, errdefs.ErrAccessDenied) {
		return httperrdefs.ErrAccessDenied
	}
	if errors.Is(err, errdefs.ErrTemporaryUnavailable) {
		return httperrdefs.ErrTemporaryUnavailable
	}
	if errors.Is(err, errdefs.ErrConnection) {
		return httperrdefs.ErrConnection
	}
	return httperrdefs.ErrUnknown
}

func (h *Handler) logError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	h.log.Error("handle request",
		zap.String(middleware.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
}
