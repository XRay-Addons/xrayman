package users

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	// add new user to storage, assign UserID to user
	NewUser(ctx context.Context, user *models.User) error
	// get user by id, return (nil, false, nil) if not exists
	GetUser(ctx context.Context, id models.UserID) (*models.User, bool, error)
	// get all users
	ListUsers(ctx context.Context) ([]models.User, error)
	// change user target status
	SetTargetUserStatus(ctx context.Context, id models.UserID,
		status models.UserStatus) error
	// delete user
	DeleteUser(ctx context.Context,
		id models.UserID) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
