package service

import "github.com/XRay-Addons/xrayman/node/internal/models"

type ServerCfg interface {
	GetUsersCfg(users []models.User) (string, error)
}
