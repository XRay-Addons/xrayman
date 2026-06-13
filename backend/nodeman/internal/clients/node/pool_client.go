package node

import (
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/nodestats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type PoolClient struct {
	sec        PoolSecurity
	httpClient HTTPClientFactory
	log        *zap.Logger
}

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

func WithLogger(l *zap.Logger) Option {
	return func(s *PoolClient) {
		if l != nil {
			s.log = l
		}
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

func (c *PoolClient) GetNodeClient(cfg models.NodeConnectionInfo) (*NodeClient, error) {
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

	client, err := api.NewClient(cfg.Endpoint, nodeSec,
		api.WithClient(httpClient))
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &NodeClient{
		client: client,
	}, nil
}

// poolsync client impl
type poolsyncClient struct {
	c *PoolClient
}

var _ poolsync.Client = (*poolsyncClient)(nil)

func (p *poolsyncClient) GetNodeClient(conn models.NodeConnectionInfo) (nodesync.Client, error) {
	return p.c.GetNodeClient(conn)
}

func (c *PoolClient) PoolSyncClient() poolsync.Client {
	return &poolsyncClient{c: c}
}

// poolstats client impl
type poolstatsClient struct {
	c *PoolClient
}

var _ poolstats.Client = (*poolstatsClient)(nil)

func (p *poolstatsClient) GetNodeClient(conn models.NodeConnectionInfo) (nodestats.Client, error) {
	return p.c.GetNodeClient(conn)
}

func (c *PoolClient) PoolStatsClient() poolstats.Client {
	return &poolstatsClient{c: c}
}
