package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/constants"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
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

func (h *Handler) NewNode(ctx context.Context, req *api.NewNodeRequest) (*api.NewNodeResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertNewNodeRequest(req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	res, err := h.service.NewNode(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertNewNodeResult(res), nil
}

func (h *Handler) StartNode(ctx context.Context, req *api.StartNodeRequest) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertStartNodeRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.service.StartNode(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) StopNode(ctx context.Context, req *api.StopNodeRequest) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertStopNodeRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.service.StopNode(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) ListNodes(ctx context.Context) (*api.ListNodeResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	res, err := h.service.ListNodes(ctx, models.ListNodeParams{})
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertListNodesResult(res), nil
}

func (h *Handler) NewUser(ctx context.Context, req *api.NewUserRequest) (*api.NewUserResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertNewUserRequest(req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	res, err := h.service.NewUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertNewUserResult(res), nil
}

func (h *Handler) EnableUser(ctx context.Context, req *api.EnableUserRequest) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertEnableUserRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.service.EnableUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) DisableUser(ctx context.Context, req *api.DisableUserRequest) error {
	if h == nil || h.service == nil {
		return fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertDisableUserRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.service.DisableUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) ListUsers(ctx context.Context) (*api.ListUsersResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	res, err := h.service.ListUsers(ctx, models.ListUserParams{})
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertListUsersResult(res), nil
}

func (h *Handler) GetUserSub(ctx context.Context, req api.GetUserSubParams) (api.GetUserSubResponse, error) {
	if h == nil || h.service == nil {
		return nil, fmt.Errorf("handler impl: %w", errdefs.ErrNilObjectCall)
	}
	p, err := converter.ConvertUserSubRequest(&req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	sub, err := h.service.GetUserSub(ctx, *p)
	if err != nil {
		return nil, err
	}
	subResponse, err := converter.ConvertUserSubResult(sub)
	if err != nil {
		return nil, err
	}
	return *subResponse, nil
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
		zap.String(constants.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
}
