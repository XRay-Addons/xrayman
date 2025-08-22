package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
)

type PoolClient struct {
	sec        PoolSecurity
	httpClient HTTPClientFactory
}

var _ pool.Client = (*PoolClient)(nil)

type Option = func(pc *PoolClient)

func WithSecIssuer(iss string) Option {
	return func(s *PoolClient) {
		s.sec.issuer = iss
	}
}

func WithSecExpiration(exp time.Duration) Option {
	return func(s *PoolClient) {
		s.sec.expiration = exp
	}
}

func WithHTTPClient(h HTTPClientFactory) Option {
	return func(s *PoolClient) {
		s.httpClient = h
	}
}

func NewPoolClient(opts ...Option) (*PoolClient, error) {
	pc := &PoolClient{
		sec: PoolSecurity{
			issuer:     "node manager",
			expiration: 10 * time.Minute,
		},
	}
	for _, o := range opts {
		o(pc)
	}
	return pc, nil
}

func (c *PoolClient) GetNodeClient(ctx context.Context,
	cfg models.NodeConnectionInfo,
) (pool.NodeClient, error) {
	if c == nil {
		return nil, errdefs.NewNilCall()
	}

	var err error
	var httpClient *http.Client
	if c.httpClient != nil {
		if httpClient, err = c.httpClient.GetNodeClient(cfg.AccessKey.CertHash); err != nil {
			return nil, fmt.Errorf("node client factory: get: %w", err)
		}
	}

	nodeSec, err := c.sec.GetNodeSecurity(cfg.AccessKey.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("node client factory: get: %w", err)
	}

	client, err := api.NewClient(cfg.Endpoint,
		nodeSec, api.WithClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("node client init: %w", err)
	}
	return &NodeClient{
		client: client,
	}, nil
}
