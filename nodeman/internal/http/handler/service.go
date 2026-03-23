package handler

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

//go:generate mockgen -source=service.go -destination=./mocks/mock_service.go -package=mocks
type Service interface {
	NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error)
	StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error)
	StopNode(ctx context.Context, p models.StopNodeParams) (*models.StopNodeResult, error)
	ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error)

	NewUser(ctx context.Context, p models.NewUserParams) (*models.User, error)
	GetUser(ctx context.Context, p models.GetUserParams) (*models.User, bool, error)
	ListUsers(ctx context.Context, p models.ListUserParams) (*models.ListUsersResult, error)
	DisableUser(ctx context.Context, p models.DisableUserParams) (*models.DisableUserResult, error)
	EnableUser(ctx context.Context, p models.EnableUserParams) (*models.EnableUserResult, error)

	GetUserSub(ctx context.Context, p models.UserSubParams) (*models.UserSubResult, bool, error)
}
