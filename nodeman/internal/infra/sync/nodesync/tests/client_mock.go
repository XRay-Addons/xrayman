package tests

import (
	"context"
	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// implement node emulator for tests
type ClientMock struct {
	Status models.NodeStatus
	Users  map[models.UserProfile]struct{}
}

func NewClientMock() *ClientMock {
	return &ClientMock{
		Status: models.NodeStatusStopped,
		Users:  make(map[models.UserProfile]struct{}, 0),
	}
}

var _ nodesync.Client = (*ClientMock)(nil)

func (c *ClientMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	return c.Status, nil
}

func (c *ClientMock) Start(ctx context.Context, users []models.UserProfile) (
	*models.ClientConfigTemplate, error,
) {
	for u := range c.Users {
		delete(c.Users, u)
	}
	for _, u := range users {
		c.Users[u] = struct{}{}
	}
	c.Status = models.NodeStatusRunning
	return &models.ClientConfigTemplate{}, nil
}

func (c *ClientMock) Stop(ctx context.Context) error {
	for u := range c.Users {
		delete(c.Users, u)
	}
	c.Status = models.NodeStatusStopped
	return nil
}

func (c *ClientMock) UpdateUsers(ctx context.Context,
	upd models.NodeUsersUpdate,
) error {
	if c.Status != models.NodeStatusRunning {
		return errdefs.New("node not running")
	}
	for _, u := range upd.Add {
		c.Users[u] = struct{}{}
	}
	for _, u := range upd.Remove {
		delete(c.Users, u)
	}
	return nil
}

// storage mock with external faults or edit state modifications
type UnstableClientMock struct {
	BaseClient  *ClientMock
	Instability float32
	rand        *rand.Rand
}

func NewUnstableClientMock() *UnstableClientMock {
	return &UnstableClientMock{
		BaseClient: NewClientMock(),
		rand:       rand.New(rand.NewPCG(0, 0)), // #nosec
	}
}

func (c *UnstableClientMock) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if c.rand.Float32() < c.Instability {
		return models.NodeStatusUnknown, errdefs.New("random client fail")
	}
	if c.rand.Float32() < c.Instability {
		if err := c.Stop(ctx); err != nil {
			return models.NodeStatusUnknown, err
		}
	}
	return c.BaseClient.CheckStatus(ctx)
}

func (c *UnstableClientMock) Start(ctx context.Context, users []models.UserProfile) (
	*models.ClientConfigTemplate, error,
) {
	if c.rand.Float32() < c.Instability {
		return nil, errdefs.New("random client fail")
	}
	return c.BaseClient.Start(ctx, users)
}

func (c *UnstableClientMock) Stop(ctx context.Context) error {
	if c.rand.Float32() < c.Instability {
		return errdefs.New("random client fail")
	}
	return c.BaseClient.Stop(ctx)
}

func (c *UnstableClientMock) UpdateUsers(ctx context.Context,
	upd models.NodeUsersUpdate,
) error {
	if c.rand.Float32() < c.Instability {
		return errdefs.New("random client fail")
	}
	return c.BaseClient.UpdateUsers(ctx, upd)
}
