package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/common/http/httperr"
	"github.com/XRay-Addons/xrayman/common/http/middleware"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperrdefs"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.uber.org/zap"
)

type Handler struct {
	us  UsersService
	ns  NodesService
	ss  SubscrService
	as  AuthService
	log *zap.Logger
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

func New(
	us UsersService,
	ns NodesService,
	ss SubscrService,
	as AuthService,
	opts ...option,
) (*Handler, error) {
	if us == nil {
		return nil, errdefs.NilArg("us")
	}
	if ns == nil {
		return nil, errdefs.NilArg("ns")
	}
	if ss == nil {
		return nil, errdefs.NilArg("ss")
	}
	handler := &Handler{
		us:  us,
		ns:  ns,
		ss:  ss,
		as:  as,
		log: zap.NewNop(),
	}
	for _, o := range opts {
		o(handler)
	}
	return handler, nil
}

/*func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// if err = error + status (this is our error), return status, log error,
	// elsewise log err and return unknonw
	if e, s := httperr.ExtractStatus[api.ErrorStatusCode](err); s != nil {
		h.logError(ctx, e)
		return s
	}
	// ogen security errors
	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		h.logError(ctx, err)
		return httperrdefs.ErrAuthToken
	}
	h.logError(ctx, err)
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
}*/

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
	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return httperrdefs.ErrAuthToken
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
	if errors.Is(err, errdefs.ErrInvaildPayload) {
		return httperrdefs.ErrInvaildPayload
	}
	if errors.Is(err, errdefs.ErrNotFound) {
		return httperrdefs.ErrNotFound
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
