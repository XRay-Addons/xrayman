package handler

import (
	"context"
	"errors"

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
	users    UsersService
	nodes    NodesService
	subscr   SubscrService
	auth     AuthService
	settings SettingsService
	log      *zap.Logger
}

var _ api.Handler = (*Handler)(nil)

func New(
	users UsersService,
	nodes NodesService,
	subscr SubscrService,
	settings SettingsService,
	auth AuthService,
	logger *zap.Logger,
) (*Handler, error) {
	if users == nil {
		return nil, errdefs.NilArg("users")
	}
	if nodes == nil {
		return nil, errdefs.NilArg("nodes")
	}
	if subscr == nil {
		return nil, errdefs.NilArg("subscr")
	}
	if settings == nil {
		return nil, errdefs.NilArg("settings")
	}
	if auth == nil {
		return nil, errdefs.NilArg("auth")
	}
	if logger == nil {
		return nil, errdefs.NilArg("logger")
	}
	return &Handler{
		users:    users,
		nodes:    nodes,
		subscr:   subscr,
		settings: settings,
		auth:     auth,
		log:      logger,
	}, nil
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

		headersResp.Set(h.Key, h.Value)
	}
	return nil
}
