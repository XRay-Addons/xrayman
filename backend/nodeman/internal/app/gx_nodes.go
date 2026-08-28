package app

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"go.uber.org/zap"
)

type HttpClientParams struct {
	gx.In
	Lc      gx.Lifecycle
	Timeout time.Duration `name:"node-call-timeout"`
	Log     *zap.Logger
}

var httpClient = gx.ProvideAnnotated(
	func(p HttpClientParams) (*httpclient.ClientFactory, error) {
		f, err := httpclient.NewClientFactory(p.Timeout, p.Log)
		if err != nil {
			return nil, err
		}
		p.Lc.AppendCloser(gx.Closer{
			Name: "client factory",
			OnClose: func(context.Context) error {
				f.Close()
				return nil
			},
		})
		return f, nil
	},
	gx.As(new(node.HTTPClientFactory)),
	gx.As(gx.Self()),
)

var poolClient = gx.Provide(
	node.NewPoolClient,
	func(pc *node.PoolClient) poolsync.Client {
		return pc.PoolSyncClient()
	},
	func(pc *node.PoolClient) poolstats.Client {
		return pc.PoolStatsClient()
	},
)

var poolSync = gx.ProvideAnnotated(
	poolsync.New,
	gx.As(new(users.Syncer)),
	gx.As(new(nodes.Syncer)),
	gx.As(new(syncman.PoolSyncer)),
)

var poolStats = gx.ProvideAnnotated(
	poolstats.New,
	gx.As(new(statsman.StatsUpdater)),
)

var Nodes = gx.Module("nodes",
	httpClient,
	poolClient,
	poolSync,
	poolStats,
)
