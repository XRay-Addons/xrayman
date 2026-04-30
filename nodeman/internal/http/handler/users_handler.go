package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func (h *Handler) NewUser(ctx context.Context, req *api.NewUserRequest) (*api.User, error) {
	if h == nil || h.us == nil {
		return nil, errdefs.NewNilCall()
	}
	p, err := converter.ConvertNewUserRequest(req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	res, err := h.us.NewUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertUser(res), nil
}

func (h *Handler) GetUser(ctx context.Context, req api.GetUserParams) (*api.User, error) {
	if h == nil || h.us == nil {
		return nil, errdefs.NewNilCall()
	}
	p, err := converter.ConvertGetUserRequest(&req)
	if err != nil {
		return nil, httperr.ErrInvaildPayload
	}
	user, exists, err := h.us.GetUser(ctx, *p)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, httperr.ErrUserNotFound
	}
	userResponse := converter.ConvertUser(user)
	return userResponse, nil
}

func (h *Handler) ListUsers(ctx context.Context) (*api.ListUsersResponse, error) {
	if h == nil || h.us == nil {
		return nil, errdefs.NewNilCall()
	}
	res, err := h.us.ListUsers(ctx, models.ListUserParams{})
	if err != nil {
		h.logError(ctx, err)
		return nil, httperr.ErrInternalServerError
	}
	return converter.ConvertListUsersResult(res), nil
}

func (h *Handler) EnableUser(ctx context.Context, req *api.EnableUserRequest) error {
	if h == nil || h.us == nil {
		return errdefs.NewNilCall()
	}
	p, err := converter.ConvertEnableUserRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.us.EnableUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) DisableUser(ctx context.Context, req *api.DisableUserRequest) error {
	if h == nil || h.us == nil {
		return errdefs.NewNilCall()
	}
	p, err := converter.ConvertDisableUserRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.us.DisableUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) error {
	if h == nil || h.us == nil {
		return errdefs.NewNilCall()
	}
	p, err := converter.ConvertDeleteUserRequest(req)
	if err != nil {
		return httperr.ErrInvaildPayload
	}
	_, err = h.us.DeleteUser(ctx, *p)
	if err != nil {
		h.logError(ctx, err)
		return httperr.ErrInternalServerError
	}
	return nil
}
