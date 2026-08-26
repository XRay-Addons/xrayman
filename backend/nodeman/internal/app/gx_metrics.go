package app

import (
	"context"
	"net/http"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var metricsCollector = gx.ProvideAnnotated(
	metrics.New,
	gx.As(new(prometheus.Collector)),
)

type MetricsRouterParams struct {
	gx.In
	PC  prometheus.Collector
	Log *zap.Logger
}

var metricsRouter = gx.ProvideNamed(
	func(p MetricsRouterParams) (http.Handler, error) {
		reg := prometheus.NewRegistry()
		if err := reg.Register(p.PC); err != nil {
			return nil, err
		}

		return router.New(
			router.WithHandler("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})),
			router.WithLogger(p.Log))
	},
	"metrics-router",
)

type MetricsServerParams struct {
	gx.In
	Router   http.Handler `name:"metrics-router"`
	Endpoint string       `name:"metrics-endpoint"`
}

var metricsServer = gx.ProvideNamed(
	func(p MetricsServerParams) (*server.HttpServer, error) {
		return server.New(p.Endpoint, p.Router)
	},
	"metrics-server",
)

type MetricsServerJobParams struct {
	gx.In
	Endpoint string             `name:"metrics-endpoint"`
	S        *server.HttpServer `name:"metrics-server"`
}

var metricsServerJob = gx.Invoke(
	func(p MetricsServerJobParams, lc gx.Lifecycle) {
		if p.Endpoint == "" {
			return
		}

		lc.AppendJob(gx.Job{
			Name: "metrics server",
			OnStart: func(context.Context) error {
				return p.S.Listen()
			},
			OnStop: func(ctx context.Context) error {
				return p.S.Shutdown(ctx)
			},
		})
	},
)

var MetricsServer = gx.Module("metrics",
	metricsCollector,
	metricsRouter,
	metricsServer,
	metricsServerJob,
)
