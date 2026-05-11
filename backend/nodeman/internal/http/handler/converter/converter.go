package converter

import (
	"fmt"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:extend ConvertExpireTime
// goverter:extend ConvertNodeID RConvertNodeID
// goverter:extend ConvertAccessKey RConvertAccessKey
// goverter:extend ConvertUserID RConvertUserID
// goverter:extend ConvertUserStatusResult
// goverter:enum:unknown @panic
//
//go:generate goverter gen .
type Converter interface {
	ConvertAuthRequest(r *api.AuthRequest) (*models.AuthParams, error)
	ConvertAuthResult(r *models.AuthResult) *api.AuthResponse

	ConvertNewNodeRequest(r *api.NewNodeRequest) (*models.NewNodeParams, error)
	ConvertNewNodeResult(r *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(r *api.StartNodeRequest) (*models.StartNodeParams, error)

	ConvertStopNodeRequest(r *api.StopNodeRequest) (*models.StopNodeParams, error)

	ConvertListNodesResult(r *models.ListNodeResult) *api.ListNodeResponse

	ConvertDeleteNodeRequest(r *api.DeleteNodeRequest) (*models.DeleteNodeParams, error)

	ConvertNewUserRequest(r *api.NewUserRequest) (*models.NewUserParams, error)

	ConvertUser(r *models.User) *api.User

	ConvertGetUserRequest(r *api.GetUserParams) (*models.GetUserParams, error)

	ConvertEnableUserRequest(r *api.EnableUserRequest) (*models.EnableUserParams, error)

	ConvertDisableUserRequest(r *api.DisableUserRequest) (*models.DisableUserParams, error)

	ConvertListUsersResult(r *models.ListUsersResult) *api.ListUsersResponse

	ConvertDeleteUserRequest(r *api.DeleteUserRequest) (*models.DeleteUserParams, error)

	ConvertUserSubRequest(r *api.UserSubParams) (*models.UserSubParams, error)
	// goverter:map ClientConfigs Response
	// goverter:map . Routing | GetSubscriptionRouting
	ConvertUserSubResult(r *models.UserSubResult) (*api.UserSubResponseHeaders, error)

	//ConvertUserSubResultBody(r []models.ClientConfigItem) (api.UserSubContent, error)

	// goverter:map . SubscriptionPath | GetUserSubscription
	ConvertProfile(r models.UserProfile) api.UserProfile

	ConvertNodeStatusResult(source models.NodeStatus) api.NodeStatus
}

func ConvertExpireTime(i time.Duration) int {
	return int(i.Seconds())
}

func ConvertNodeID(i models.NodeID) api.NodeID {
	return api.NodeID(i)
}

func RConvertNodeID(i api.NodeID) models.NodeID {
	return models.NodeID(i)
}

func ConvertAccessKey(s string) (models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := accessKey.UnmarshalText([]byte(s)); err != nil {
		return accessKey, errdefs.InvalidPayload(err.Error())
	}
	return accessKey, nil
}

func RConvertAccessKey(key models.AccessKey) string {
	return key.String()
}

func ConvertUserID(i models.UserID) api.UserID {
	return api.UserID(i)
}

func RConvertUserID(i api.UserID) models.UserID {
	return models.UserID(i)
}

func ConvertUserStatusResult(source models.UserStatus) api.UserStatus {
	var response api.UserStatus
	switch source {
	case models.UserStatusDisabled:
		response = api.UserStatusDisabled
	case models.UserStatusEnabled:
		response = api.UserStatusEnabled
	case models.UserStatusUnknown:
		response = api.UserStatusUnknown
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", source))
	}
	return response
}

func GetUserSubscription(source models.UserProfile) string {
	return source.SubscriptionURL()
}

func GetSubscriptionRouting(source *models.UserSubResult) api.OptString {
	if source.Headers.Routing != nil {
		return api.NewOptString(*source.Headers.Routing)
	}
	return api.OptString{}
}
