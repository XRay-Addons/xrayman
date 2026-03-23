package converter

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:extend ConvertNodeID RConvertNodeID
// goverter:extend ConvertNodeStatusResult
// goverter:extend ConvertAccessKey RConvertAccessKey
// goverter:extend ConvertUserID RConvertUserID
// goverter:extend ConvertUserStatusResult
//
//go:generate goverter gen .
type Converter interface {
	ConvertNewNodeRequest(r *api.NewNodeRequest) (*models.NewNodeParams, error)
	ConvertNewNodeResult(r *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(r *api.StartNodeRequest) (*models.StartNodeParams, error)

	ConvertStopNodeRequest(r *api.StopNodeRequest) (*models.StopNodeParams, error)

	ConvertListNodesResult(r *models.ListNodeResult) *api.ListNodeResponse

	ConvertNewUserRequest(r *api.NewUserRequest) (*models.NewUserParams, error)

	ConvertUser(r *models.User) *api.User

	ConvertGetUserRequest(r *api.GetUserParams) (*models.GetUserParams, error)

	ConvertEnableUserRequest(r *api.EnableUserRequest) (*models.EnableUserParams, error)

	ConvertDisableUserRequest(r *api.DisableUserRequest) (*models.DisableUserParams, error)

	ConvertListUsersResult(r *models.ListUsersResult) *api.ListUsersResponse

	ConvertUserSubRequest(r *api.UserSubParams) (*models.UserSubParams, error)
	// goverter:map ClientConfigs Response
	ConvertUserSubResult(r *models.UserSubResult) (*api.UserSubResponseHeaders, error)
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
		return accessKey, errdefs.WrapWithStack(err)
	}
	return accessKey, nil
}

func RConvertAccessKey(key models.AccessKey) string {
	return key.String()
}

func ConvertNodeStatusResult(source models.NodeStatus) api.NodeStatus {
	var response api.NodeStatus
	switch source {
	case models.NodeStatusStopped:
		response = api.NodeStatusStopped
	case models.NodeStatusRunning:
		response = api.NodeStatusRunning
	case models.NodeStatusUnknown:
		response = api.NodeStatusUnknown
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", source))
	}
	return response
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

/*func ConvertClientConfig(source models.Subscription) (api.Subscription, error) {
	var s api.Subscription
	if err := json.Unmarshal([]byte(source), &s); err != nil {
		return nil, errdefs.WrapWithStack(err)
	}
	return s, nil
}*/
