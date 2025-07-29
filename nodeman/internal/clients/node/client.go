package client

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
)

type NodeClient struct {
	client *api.Client
}

var _ service.NodeClient = (*NodeClient)(nil)

func NewNodeClient(endpoint string, security api.SecuritySource, httpClient *http.Client) (
	*NodeClient, error,
) {
	opts := make([]api.ClientOption, 0)
	if httpClient != nil {
		opts = append(opts, api.WithClient(httpClient))
	}

	client, err := api.NewClient(endpoint, security, opts...)
	if err != nil {
		return nil, fmt.Errorf("node client init: %w", err)
	}

	return &NodeClient{
		client: client,
	}, nil
}

func (c *NodeClient) Start(ctx context.Context, users []models.UserProfile) (
	*models.ClientConfig, error,
) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("node client: start: %w", errdefs.ErrNilObjectCall)
	}

	startRequest := api.StartRequest{Users: converter.ConvertUsers(users)}
	startResponse, err := c.client.StartPost(ctx, &startRequest)
	if err != nil {
		return nil, fmt.Errorf("node client: start: %w", err)
	}

	clientTemplate := converter.ConvertClientCfg(startResponse.GetClientCfg())
	return &clientTemplate, nil
}

func (c *NodeClient) Stop(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("node client: stop: %w", errdefs.ErrNilObjectCall)
	}

	if err := c.client.StopPost(ctx); err != nil {
		return fmt.Errorf("node client: stop: %w", err)
	}
	return nil
}

func (c *NodeClient) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if c == nil || c.client == nil {
		return models.NodeStatusUnknown,
			fmt.Errorf("node client: status: %w", errdefs.ErrNilObjectCall)
	}

	status, err := c.client.GetStatus(ctx)
	if err != nil {
		return models.NodeStatusUnknown, fmt.Errorf("node client: status: %w", err)
	}
	return converter.ConvertNodeStatus(status.ServiceStatus), nil
}

func (c *NodeClient) UpdateUserStates(ctx context.Context, update models.NodeUsersUpdate) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("node client: status: %w", errdefs.ErrNilObjectCall)
	}

	editRequest := converter.ConvertUsersUpdate(update)
	if err := c.client.EditUsers(ctx, &editRequest); err != nil {
		return fmt.Errorf("node client: update users: %w", err)
	}
	return nil
}
