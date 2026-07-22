package handler

import (
	"context"
	"errors"
	"fmt"

	mw "github.com/XRay-Addons/xrayman/common/http/middleware"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperrdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.uber.org/zap"
)

type Handler struct {
	us  UsersService
	ns  NodesService
	ss  SubscrService
	dcs DynConfigService
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
	dcs DynConfigService,
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
	if dcs == nil {
		return nil, errdefs.NilArg("dcs")
	}
	if as == nil {
		return nil, errdefs.NilArg("as")
	}
	handler := &Handler{
		us:  us,
		ns:  ns,
		ss:  ss,
		dcs: dcs,
		as:  as,
		log: zap.NewNop(),
	}
	for _, o := range opts {
		o(handler)
	}
	return handler, nil
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	// log error
	h.logError(ctx, err)

	// if error contains status, return it
	var statusError *api.ErrorStatusCode
	if errors.As(err, &statusError) {
		return statusError
	}

	// translate error to status
	return h.translateError(err)
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
		zap.String(mw.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
}

func (h *Handler) writeHeaders(ctx context.Context, headers models.SubHeaders) error {
	headersResp := mw.GetHeaders(ctx)
	if headersResp == nil {
		return xerr.New("request context doesn't support headers")
	}
	for _, h := range headers {
		fmt.Println(h.Key, h.Value)
		headersResp.Set(h.Key, h.Value)
	}
	return nil
}
