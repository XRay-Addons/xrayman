package requests

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

const EditUsersURLPath = "/users/edit"

type EditUsersRequest struct {
	Add    []models.User
	Remove []models.User
}
