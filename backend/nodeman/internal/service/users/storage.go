package users

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type UoWContext interface {
	// add new user to storage, assign UserID to user
	NewUser(ctx context.Context, user *models.User) error
	// get user by id, return ErrNotFound if not exists
	GetUserView(ctx context.Context, id models.UserID, name string) (*models.UserView, error)
	// get all users
	ListUserViews(ctx context.Context) ([]models.UserView, error)
	// change user target status
	SetTargetUserStatus(ctx context.Context, id models.UserID,
		status models.UserStatus) error
	// delete user
	DeleteUser(ctx context.Context,
		id models.UserID) error
}

type UoWFn = uow.Fn[UoWContext]
type Storage = uow.Storage[UoWContext]
