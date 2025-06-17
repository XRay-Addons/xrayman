package nodeclient

import (
	"bytes"
	"context"
	"encoding/json"

	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/taskpool"
	"github.com/XRay-Addons/xrayman/shared/models"
	"github.com/go-resty/resty/v2"
)

type NodeClient struct {
	client *resty.Client
	tp     *taskpool.TaskPool
}

const (
	requestTimeout = 10 * time.Second
)

func New(endpoint string, rateLimit int) (*NodeClient, error) {
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, fmt.Errorf("invalid node endpoint: %w", err)
	}

	client := resty.New()
	client.SetBaseURL(endpoint)
	client.SetTimeout(requestTimeout)

	return &NodeClient{
		client: client,
		tp:     taskpool.New(rateLimit),
	}, nil
}

func (c *NodeClient) Close() {
	if c == nil || c.tp == nil {
		return
	}
	c.tp.Close()
}

func (c *NodeClient) Start(ctx context.Context, users []models.User) (*models.NodeProperties, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("node client not initialized")
	}

	resp, err := c.callNodeAPI(ctx, func(ctx context.Context) (*resty.Response, error) {
		req := models.StartNodeRequest{Users: users}
		return c.postJSON(ctx, "/node/start", &req)
	})

	if err != nil {
		return nil, fmt.Errorf("call start node API: %w", err)
	}

	var response models.StartNodeResponse
	if err := c.parseResponseJSON(resp, &response); err != nil {
		return nil, fmt.Errorf("parse /node/start response json: %w", err)
	}

	return &response.NodeProperties, nil
}

func (c *NodeClient) Stop(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("node client not initialized")
	}

	resp, err := c.callNodeAPI(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.client.R().Post("/node/stop")
	})

	if err != nil {
		return fmt.Errorf("call stop node API: %w", err)
	}

	if err := c.checkResponse(resp); err != nil {
		return fmt.Errorf("/node/stop response: %w", err)
	}

	return nil
}

func (c *NodeClient) Status(ctx context.Context) (*models.NodeStatus, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("node client not initialized")
	}

	resp, err := c.callNodeAPI(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.client.R().Get("/node/status")
	})

	if err != nil {
		return nil, fmt.Errorf("call status node API: %w", err)
	}

	var response models.NodeStatusResponse
	if err := c.parseResponseJSON(resp, &response); err != nil {
		return nil, fmt.Errorf("/node/status response: %w", err)
	}

	return &response.NodeStatus, nil
}

type ApiFunc = func(ctx context.Context) (*resty.Response, error)

func (c *NodeClient) callNodeAPI(ctx context.Context, apiFn ApiFunc) (*resty.Response, error) {
	const retries = 3

	var err error
	for range retries {
		delay := 1 * time.Second
		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("node api call attempts cancelled")
		case <-timer.C:
			resp, err := apiFn(ctx)
			if err == nil {
				return resp, nil
			}
			if !isServerUnavailableErr(err) {
				return nil, fmt.Errorf("call node API: %w", err)
			}
		}
	}
	return nil, fmt.Errorf("api call last try: %w", err)
}

func (c *NodeClient) postJSON(ctx context.Context, url string, object interface{}) (*resty.Response, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("json marshalling body: %v", err)
	}

	// !important: we must post only raw data as io.Reader,
	// all other variants like SetBody(object) or SetBody(body)
	// doesn't guarantee that request bodt will be equal to 'body'
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewReader(body)).
		SetContext(ctx)

	return req.Post(url)
}

func isServerUnavailableErr(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		// too long for request timeout
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}

func (c *NodeClient) parseResponseJSON(resp *resty.Response, v interface{}) error {
	if err := c.checkResponse(resp); err != nil {
		return fmt.Errorf("check json response: %w", err)
	}
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return fmt.Errorf("unmarshal response body: %w", err)
	}
	return nil
}

func (c *NodeClient) checkResponse(resp *resty.Response) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return nil
}
