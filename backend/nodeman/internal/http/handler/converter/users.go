package converter

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./users_generated.go
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Users interface {
	ConvertNewUserRequest(r *api.NewUserRequest) (*models.NewUserParams, error)
	ConvertNewUserResult(r *models.User) *api.User

	ConvertGetUserRequest(r *api.GetUserParams) (*models.GetUserParams, error)
	ConvertGetUserResult(r *models.UserView) *api.UserView

	ConvertEnableUserRequest(r *api.EnableUserRequest) (*models.EnableUserParams, error)

	ConvertDisableUserRequest(r *api.DisableUserRequest) (*models.DisableUserParams, error)

	ConvertListUsersResult(r *models.ListUsersResult) *api.ListUsersResponse

	ConvertDeleteUserRequest(r *api.DeleteUserRequest) (*models.DeleteUserParams, error)

	// goverter:map . SubscriptionPath | GetUserSubscription
	ConvertProfile(r models.UserProfile) api.UserProfile
}

func GetUserSubscription(source models.UserProfile) string {
	return source.SubscriptionURL()
}
