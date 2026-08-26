package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Performance interface {
	GetPerformance(context.Context) (*models.NodePerformance, error)
}
