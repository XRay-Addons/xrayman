package node

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeAPI interface {
	Start(ctx context.Context, users []models.UserProfile) (*models.ClientTemplate, error)
	Stop(ctx context.Context) error
	CheckStatus(ctx context.Context) (models.NodeStatus, error)
	UpdateUserStates(ctx context.Context, transitions []models.UserTargetState) error
}
