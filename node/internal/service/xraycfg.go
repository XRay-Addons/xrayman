package service

import "github.com/XRay-Addons/xrayman/shared/models"

type XRayCfg interface {
	GetInbounds() []models.Inbound
	GetServerConfig(users []models.User) (string, error)
	GetClientConfig() string
}
