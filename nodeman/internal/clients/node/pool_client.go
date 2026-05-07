package nodesync

import (
	"context"
	"net/http"
	"time"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	pool "github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
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

const (
	defaultCertExpiration = 10 * time.Minute
)

func NewPoolClient(opts ...Option) (*PoolClient, error) {
	pc := &PoolClient{
		sec: PoolSecurity{
			issuer:     "node manager",
			expiration: defaultCertExpiration,
		},
	}
	for _, o := range opts {
		o(pc)
	}
	return pc, nil
}

func (c *PoolClient) GetNodeClient(ctx context.Context,
	cfg models.NodeConnectionInfo,
) (nodesync.Client, error) {
	if c == nil {
		return nil, errdefs.NilCall()
	}

	var err error
	var httpClient *http.Client
	if c.httpClient != nil {
		if httpClient, err = c.httpClient.GetNodeClient(cfg.AccessKey.CertHash); err != nil {
			return nil, xerr.WrapWithStack(err)
		}
	}

	nodeSec, err := c.sec.GetNodeSecurity(cfg.AccessKey.AccessSecret)
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(cfg.Endpoint,
		nodeSec, api.WithClient(httpClient))
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &NodeClient{
		client: client,
	}, nil
}
