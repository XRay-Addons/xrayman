package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=service.go -destination=./mocks/mock_users_service.go -package=mocks
type UsersService interface {
	NewUser(ctx context.Context, p models.NewUserParams) (*models.User, error)
	GetUser(ctx context.Context, p models.GetUserParams) (*models.User, bool, error)
	ListUsers(ctx context.Context, p models.ListUserParams) (*models.ListUsersResult, error)
	DisableUser(ctx context.Context, p models.DisableUserParams) (*models.DisableUserResult, error)
	EnableUser(ctx context.Context, p models.EnableUserParams) (*models.EnableUserResult, error)
	DeleteUser(ctx context.Context, p models.DeleteUserParams) (*models.DeleteUserResult, error)
}
