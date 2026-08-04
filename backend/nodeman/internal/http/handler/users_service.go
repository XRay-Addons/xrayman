package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=users_service.go -destination=./mocks/mock_users_service.go -package=mocks
type UsersService interface {
	NewUser(ctx context.Context, p models.NewUserParams) (*models.User, error)
	GetUserView(ctx context.Context, p models.GetUserParams) (*models.UserView, error)
	ListUsers(ctx context.Context) (*models.ListUsersResult, error)
	DisableUser(ctx context.Context, p models.DisableUserParams) error
	EnableUser(ctx context.Context, p models.EnableUserParams) error
	DeleteUser(ctx context.Context, p models.DeleteUserParams) error
}
