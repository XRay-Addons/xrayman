package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/node"
)

// implement node emulator for tests
type APIMock struct {
	Status models.NodeStatus
	Users  map[models.UserID]struct{}
}

func NewAPIMock() *APIMock {
	return &APIMock{
		Status: models.NodeStatusStopped,
		Users:  make(map[models.UserID]struct{}, 0),
	}
}

var _ node.NodeAPI = (*APIMock)(nil)

func (api *APIMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	return api.Status, nil
}

func (api *APIMock) Start(ctx context.Context, users []models.UserProfile) (*models.ClientTemplate, error) {
	for u := range api.Users {
		delete(api.Users, u)
	}
	for _, u := range users {
		api.Users[u.ID] = struct{}{}
	}
	api.Status = models.NodeStatusRunning
	return &models.ClientTemplate{}, nil
}

func (api *APIMock) Stop(ctx context.Context) error {
	for u := range api.Users {
		delete(api.Users, u)
	}
	api.Status = models.NodeStatusStopped
	return nil
}

func (api *APIMock) UpdateUserStates(ctx context.Context, transitions []models.UserTargetState) error {
	if api.Status != models.NodeStatusRunning {
		return fmt.Errorf("node not running")
	}
	for _, t := range transitions {
		if t.Target == models.UserStatusActive {
			api.Users[t.User.ID] = struct{}{}
		} else {
			delete(api.Users, t.User.ID)
		}
	}
	return nil
}

// storage mock with external faults or edit state modifications
type UnstableAPIMock struct {
	BaseAPI     *APIMock
	Instability float32
	rand        *rand.Rand
}

func NewUnstableAPIMock() *UnstableAPIMock {
	return &UnstableAPIMock{
		BaseAPI: NewAPIMock(),
		rand:    rand.New(rand.NewPCG(0, 0)),
	}
}

func (api *UnstableAPIMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if api.rand.Float32() < api.Instability {
		return models.NodeStatusUnknown, fmt.Errorf("random api fail")
	}
	if api.rand.Float32() < api.Instability {
		return models.NodeStatusUnknown, nil
	}
	return api.BaseAPI.CheckStatus(ctx)
}

func (api *UnstableAPIMock) Start(ctx context.Context, users []models.UserProfile) (*models.ClientTemplate, error) {
	if api.rand.Float32() < api.Instability {
		return nil, fmt.Errorf("random api fail")
	}
	return api.BaseAPI.Start(ctx, users)
}

func (api *UnstableAPIMock) Stop(ctx context.Context) error {
	if api.rand.Float32() < api.Instability {
		return fmt.Errorf("random api fail")
	}
	return api.BaseAPI.Stop(ctx)
}

func (api *UnstableAPIMock) UpdateUserStates(ctx context.Context, transitions []models.UserTargetState) error {
	if api.rand.Float32() < api.Instability {
		return fmt.Errorf("random api fail")
	}
	return api.BaseAPI.UpdateUserStates(ctx, transitions)
}
