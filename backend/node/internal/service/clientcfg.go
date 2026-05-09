package service

import "github.com/XRay-Addons/xrayman/node/internal/models"

type ClientConfig interface {
	GetTemplate() (*models.ClientConfigTemplate, error)
}
