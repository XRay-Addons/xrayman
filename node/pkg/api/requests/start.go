package requests

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

const StartURLPath = "/start"

type StartRequest struct {
	Users []models.User `json:"users" validate:"required"`
}

type StartResponse struct {
	NodeConfig models.NodeConfig `json:"nodeConfig" validate:"required"`
}
