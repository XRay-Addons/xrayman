package requests

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

const EditUsersURLPath = "/users/edit"

type EditUsersRequest struct {
	Add    []models.User `json:"add" validate:"required"`
	Remove []models.User `json:"remove" validate:"required"`
}
