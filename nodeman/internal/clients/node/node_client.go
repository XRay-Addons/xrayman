package nodesync

import (
	"context"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/converter"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeClient struct {
	client *api.Client
}

var _ nodesync.Client = (*NodeClient)(nil)

func (c *NodeClient) Start(ctx context.Context, users []models.UserProfile) (
	*models.ClientConfigTemplate, error,
) {
	if c == nil || c.client == nil {
		return nil, errdefs.NewNilCall()
	}

	startRequest := api.StartRequest{Users: converter.ConvertUsers(users)}
	startResponse, err := c.client.Start(ctx, &startRequest)
	if err != nil {
		return nil, wrapOgenErr(err)
	}
	clientTemplate := converter.ConvertClientConfig(startResponse.GetClientConfigTemplate())
	return &clientTemplate, nil
}

func (c *NodeClient) Stop(ctx context.Context) error {
	if c == nil || c.client == nil {
		return errdefs.NewNilCall()
	}

	if err := c.client.Stop(ctx); err != nil {
		return wrapOgenErr(err)
	}
	return nil
}

func (c *NodeClient) CheckStatus(ctx context.Context) (models.NodeStatus, error) {
	if c == nil || c.client == nil {
		return models.NodeStatusUnknown, errdefs.NewNilCall()
	}

	status, err := c.client.GetStatus(ctx)
	if err != nil {
		return models.NodeStatusUnknown, wrapOgenErr(err)
	}
	return converter.ConvertNodeStatus(status.ServiceStatus), nil
}

func (c *NodeClient) UpdateUsers(ctx context.Context, update models.NodeUsersUpdate) error {
	if c == nil || c.client == nil {
		return errdefs.NewNilCall()
	}

	editRequest := converter.ConvertUsersUpdate(update)
	if err := c.client.EditUsers(ctx, &editRequest); err != nil {
		return wrapOgenErr(err)
	}
	return nil
}

func wrapOgenErr(err error) error {
	return errdefs.Wrap(err, errdefs.WithOgen(), errdefs.WithStack())
}
