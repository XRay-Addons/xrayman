package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./users_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Users interface {
	ConvertNewUserRequest(r *ogenserver.NewUserRequest) (*models.NewUserParams, error)
	ConvertNewUserResult(r *models.User) *ogenserver.User

	ConvertGetUserRequest(r *ogenserver.GetUserParams) (*models.GetUserParams, error)
	ConvertGetUserResult(r *models.UserView) *ogenserver.UserView

	ConvertEnableUserRequest(r *ogenserver.EnableUserRequest) (*models.EnableUserParams, error)

	ConvertDisableUserRequest(r *ogenserver.DisableUserRequest) (*models.DisableUserParams, error)

	ConvertListUsersResult(r *models.ListUsersResult) *ogenserver.ListUsersResponse

	ConvertDeleteUserRequest(r *ogenserver.DeleteUserRequest) (*models.DeleteUserParams, error)

	// goverter:map . SubscriptionPath | GetUserSubscription
	ConvertProfile(r models.UserProfile) ogenserver.UserProfile
}

func GetUserSubscription(source models.UserProfile) string {
	return source.SubscriptionURL()
}
