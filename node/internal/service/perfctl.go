package service

import "github.com/XRay-Addons/xrayman/node/pkg/api/models"

type PerfCtl interface {
	GetSystemUsage() (*models.SystemUsage, error)
}
