package requests

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

const StatusURLPath = "/status"

type StatusResponse struct {
	Status models.NodeStatus `json:"status" validate:"required"`
}
