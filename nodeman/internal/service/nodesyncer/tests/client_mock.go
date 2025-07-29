package tests

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
)

// implement node emulator for tests
type ClientMock struct {
	Status models.NodeStatus
	Users  map[models.UserID]struct{}
}

func NewClientMock() *ClientMock {
	return &ClientMock{
		Status: models.NodeStatusStopped,
		Users:  make(map[models.UserID]struct{}, 0),
	}
}

var _ nodesyncer.Client = (*ClientMock)(nil)

func (c *ClientMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	return c.Status, nil
}

func (c *ClientMock) Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfig, error) {
	for u := range c.Users {
		delete(c.Users, u)
	}
	for _, u := range users {
		c.Users[u.ID] = struct{}{}
	}
	c.Status = models.NodeStatusRunning
	return &models.ClientConfig{}, nil
}

func (c *ClientMock) Stop(ctx context.Context) error {
	for u := range c.Users {
		delete(c.Users, u)
	}
	c.Status = models.NodeStatusStopped
	return nil
}

// UpdateUserStates implements nodesyncer.NodeClient.
func (c *ClientMock) UpdateUserStates(ctx context.Context, update models.NodeUsersUpdate) error {
	if c.Status != models.NodeStatusRunning {
		return fmt.Errorf("node not running")
	}
	for _, u := range update.Add {
		c.Users[u.ID] = struct{}{}
	}
	for _, u := range update.Remove {
		delete(c.Users, u.ID)
	}
	return nil
}

// storage mock with external faults or edit state modifications
type UnstableClientMock struct {
	BaseAPI     *ClientMock
	Instability float32
	rand        *rand.Rand
}

func NewUnstableClientMock() *UnstableClientMock {
	return &UnstableClientMock{
		BaseAPI: NewClientMock(),
		rand:    rand.New(rand.NewPCG(0, 0)),
	}
}

func (c *UnstableClientMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if c.rand.Float32() < c.Instability {
		return models.NodeStatusUnknown, fmt.Errorf("random c fail")
	}
	if c.rand.Float32() < c.Instability {
		return models.NodeStatusUnknown, nil
	}
	return c.BaseAPI.CheckStatus(ctx)
}

func (c *UnstableClientMock) Start(ctx context.Context, users []models.UserProfile) (*models.ClientConfig, error) {
	if c.rand.Float32() < c.Instability {
		return nil, fmt.Errorf("random c fail")
	}
	return c.BaseAPI.Start(ctx, users)
}

func (c *UnstableClientMock) Stop(ctx context.Context) error {
	if c.rand.Float32() < c.Instability {
		return fmt.Errorf("random c fail")
	}
	return c.BaseAPI.Stop(ctx)
}

func (c *UnstableClientMock) UpdateUserStates(ctx context.Context, update models.NodeUsersUpdate) error {
	if c.rand.Float32() < c.Instability {
		return fmt.Errorf("random c fail")
	}
	return c.BaseAPI.UpdateUserStates(ctx, update)
}
