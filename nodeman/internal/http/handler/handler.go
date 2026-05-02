package handler

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/constants"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperr"
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
		return nil, errdefs.NewNilArg("us")
	}
	if ns == nil {
		return nil, errdefs.NewNilArg("ns")
	}
	if ss == nil {
		return nil, errdefs.NewNilArg("ss")
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

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// use passed HttpErr or default unknown
	httpErr := httperr.ErrUnknown
	if errors.As(err, &httpErr) {
		// all errors pass to this handler, many of them are consequences
		// of errors processed and logged before, others come here
		h.logError(ctx, err)
	}
	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		// non-httperr ogen errors
		httpErr = httperr.ErrAuthToken
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
