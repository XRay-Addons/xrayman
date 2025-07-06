package requests

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

const StartURLPath = "/start"

type StartRequest struct {
	Users []models.User
}

type StartResponse struct {
	NodeConfig models.NodeConfig
}
