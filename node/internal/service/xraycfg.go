package service

import "github.com/XRay-Addons/xrayman/node/internal/models"

type XRayCfg interface {
	GetInbounds() []models.Inbound
	SetUsers(users []models.User) error
}
