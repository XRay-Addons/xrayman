package node

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/models"
	"go.uber.org/zap"
)

// implement node emulator for tests
type APIMock struct {
	status models.NodeStatus
	users  map[models.UserID]struct{}

	rand            *rand.Rand
	failProb        float32
	turnOffProb     float32
	unavailableProb float32

	Unstable bool
}

func NewAPIMock(log *zap.Logger) *APIMock {
	return &APIMock{
		status:          models.NodeOff,
		users:           make(map[models.UserID]struct{}, 0),
		rand:            rand.New(rand.NewPCG(0, 0)),
		failProb:        0.25,
		turnOffProb:     0.25,
		unavailableProb: 0.25,
		Unstable:        true,
	}
}

var _ API = (*APIMock)(nil)

// EditUsers implements node.API.
func (a *APIMock) EditUsers(ctx context.Context, users []UserState) error {
	if _, err := a.applyUnstability(); err != nil {
		return fmt.Errorf("edit users: %w", err)
	}

	for _, u := range users {
		if u.Status == models.UserEnabled {
			a.users[u.User.ID] = struct{}{}
		} else {
			delete(a.users, u.User.ID)
		}
	}

	return nil
}

func (a *APIMock) GetStatus(ctx context.Context) (models.NodeStatus, error) {
	available, err := a.applyUnstability()
	if err != nil {
		return models.NodeStatusUnknown, fmt.Errorf("get status: %w", err)
	}
	if !available {
		return models.NodeStatusUnknown, nil
	}

	return a.status, nil
}

func (a *APIMock) Start(ctx context.Context, users []models.User) (*models.NodeProperties, error) {
	if _, err := a.applyUnstability(); err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	for u := range a.users {
		delete(a.users, u)
	}
	for _, u := range users {
		a.users[u.ID] = struct{}{}
	}
	a.status = models.NodeOn

	return &models.NodeProperties{}, nil
}

func (a *APIMock) Stop(ctx context.Context) error {
	if _, err := a.applyUnstability(); err != nil {
		return fmt.Errorf("get status: %w", err)
	}
	for u := range a.users {
		delete(a.users, u)
	}

	a.status = models.NodeOff
	return nil
}

func (a *APIMock) applyUnstability() (available bool, err error) {
	if !a.Unstable {
		return true, nil
	}
	if a.rand.Float32() < a.turnOffProb {
		return false, nil
	}

	if a.rand.Float32() < a.turnOffProb {
		a.status = models.NodeOff
		for u := range a.users {
			delete(a.users, u)
		}
	}

	if a.rand.Float32() < a.failProb {
		return true, fmt.Errorf("node api failed")
	}
	return true, nil
}
