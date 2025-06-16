package service

import "github.com/XRay-Addons/xrayman/node/internal/models"

type XRayCfg interface {
	GetInbounds() []models.Inbound
	GetServerConfig(users []models.User) (string, error)
	GetClientConfig() string
}
