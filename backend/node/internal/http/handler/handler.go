package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/common/http/middleware"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperrdefs"
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

var _ ogenserver.Handler = (*Handler)(nil)

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

func (h *Handler) Start(ctx context.Context,
	req *ogenserver.StartRequest,
) (*ogenserver.StartResponse, error) {
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
	if err := h.service.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetStatus(ctx context.Context) (
	*ogenserver.StatusResponse, error,
) {
	if h == nil || h.service == nil {
		return nil, errdefs.NilCall()
	}
	status, err := h.service.Status(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertStatusResult(status), nil
}

func (h *Handler) EditUsers(ctx context.Context,
	req *ogenserver.EditUsersRequest,
) error {
	if h == nil || h.service == nil {
		return errdefs.NilCall()
	}
	p := converter.ConvertEditUsersRequest(req)
	if err := h.service.EditUsers(ctx, *p); err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetStats(ctx context.Context) (*ogenserver.StatsResponse, error) {
	if h == nil || h.service == nil {
		return nil, errdefs.NilCall()
	}
	stats, err := h.service.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertStatsResult(stats), nil
}

func (h *Handler) NewError(ctx context.Context, err error) *ogenserver.ErrorStatusCode {
	// log error
	h.logError(ctx, err)

	// if error contains status, return it
	var statusError *ogenserver.ErrorStatusCode
	if errors.As(err, &statusError) {
		return statusError
	}

	// translate error to status
	return h.translateError(err)
}

func (h *Handler) translateError(err error) *ogenserver.ErrorStatusCode {
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
