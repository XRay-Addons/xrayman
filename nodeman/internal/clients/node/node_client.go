package node

import (
	"context"
	"fmt"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
)

type NodeClient struct {
	client *api.Client
}

var _ pool.NodeClient = (*NodeClient)(nil)

func (c *NodeClient) Start(ctx context.Context, users []models.UserProfile) (
	*models.ClientConfig, error,
) {
	if c == nil || c.client == nil {
		return nil, errdefs.NewNilCall()
	}

	startRequest := api.StartRequest{Users: converter.ConvertUsers(users)}
	startResponse, err := c.client.Start(ctx, &startRequest)
	if err != nil {
		return nil, fmt.Errorf("node client: start: %w", err)
	}

	clientTemplate := converter.ConvertClientCfg(startResponse.GetClientCfg())
	return &clientTemplate, nil
}

func (c *NodeClient) Stop(ctx context.Context) error {
	if c == nil || c.client == nil {
		return errdefs.NewNilCall()
	}

	if err := c.client.Stop(ctx); err != nil {
		return fmt.Errorf("node client: stop: %w", err)
	}
	return nil
}

func (c *NodeClient) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if c == nil || c.client == nil {
		return models.NodeStatusUnknown, errdefs.NewNilCall()
	}

	status, err := c.client.GetStatus(ctx)
	if err != nil {
		return models.NodeStatusUnknown, fmt.Errorf("node client: status: %w", err)
	}
	return converter.ConvertNodeStatus(status.ServiceStatus), nil
}

func (c *NodeClient) UpdateUsers(ctx context.Context, update models.NodeUsersUpdate) error {
	if c == nil || c.client == nil {
		return errdefs.NewNilCall()
	}

	editRequest := converter.ConvertUsersUpdate(update)
	if err := c.client.EditUsers(ctx, &editRequest); err != nil {
		return fmt.Errorf("node client: update users: %w", err)
	}
	return nil
}
