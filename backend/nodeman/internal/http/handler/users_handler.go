package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

func (h *Handler) NewUser(ctx context.Context, req *api.NewUserRequest) (*api.User, error) {
	if h == nil || h.users == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertNewUserRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.users.NewUser(ctx, *p)
	if err != nil {
		return nil, err
	}
	return converter.ConvertNewUserResult(res), nil
}

func (h *Handler) GetUser(ctx context.Context, req api.GetUserParams) (*api.UserView, error) {
	if h == nil || h.users == nil {
		return nil, errdefs.NilCall()
	}
	p, err := converter.ConvertGetUserRequest(&req)
	if err != nil {
		return nil, err
	}
	user, err := h.users.GetUserView(ctx, *p)
	if err != nil {
		return nil, err
	}
	userResponse := converter.ConvertGetUserResult(user)
	return userResponse, nil
}

func (h *Handler) ListUsers(ctx context.Context) (*api.ListUsersResponse, error) {
	if h == nil || h.users == nil {
		return nil, errdefs.NilCall()
	}
	res, err := h.users.ListUsers(ctx, models.ListUserParams{})
	if err != nil {
		return nil, err
	}
	return converter.ConvertListUsersResult(res), nil
}

func (h *Handler) EnableUser(ctx context.Context, req *api.EnableUserRequest) error {
	if h == nil || h.users == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertEnableUserRequest(req)
	if err != nil {
		return err
	}
	_, err = h.users.EnableUser(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DisableUser(ctx context.Context, req *api.DisableUserRequest) error {
	if h == nil || h.users == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDisableUserRequest(req)
	if err != nil {
		return err
	}
	_, err = h.users.DisableUser(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) error {
	if h == nil || h.users == nil {
		return errdefs.NilCall()
	}
	p, err := converter.ConvertDeleteUserRequest(req)
	if err != nil {
		return err
	}
	_, err = h.users.DeleteUser(ctx, *p)
	if err != nil {
		return err
	}
	return nil
}
