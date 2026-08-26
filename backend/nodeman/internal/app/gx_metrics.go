package app

import (
	"context"
	"net/http"

	"math/rand/v2"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsStorage struct {
}

var _ metrics.Storage = (*MetricsStorage)(nil)

func (m *MetricsStorage) GetNodeMetrics(ctx context.Context) ([]models.NodeMetrics, error) {
	return []models.NodeMetrics{
		models.NodeMetrics{
			ID:            models.NodeID(1),
			Endpoint:      "xo.stepka.co.uk",
			TotalInbount:  int64(rand.IntN(100)),
			TotalOutbound: int64(rand.IntN(100)),
			CpuLoad:       rand.Float32(),
			MemLoad:       rand.Float32(),
			RamLoad:       rand.Float32(),
		},
		models.NodeMetrics{
			ID:            models.NodeID(2),
			Endpoint:      "ox.stepka.co.uk",
			TotalInbount:  int64(rand.IntN(100)),
			TotalOutbound: int64(rand.IntN(100)),
			CpuLoad:       rand.Float32(),
			MemLoad:       rand.Float32(),
			RamLoad:       rand.Float32(),
		},
	}, nil
}

func NewMetricsStorage() *MetricsStorage {
	return &MetricsStorage{}
}

var metricsStorage = gx.ProvideAnnotated(
	NewMetricsStorage,
	gx.As(new(metrics.Storage)),
)

var metricsCollector = gx.ProvideAnnotated(
	metrics.New,
	gx.As(new(prometheus.Collector)),
)

type MetricsServerParams struct {
	gx.In
	PC       prometheus.Collector
	Endpoint string `name:"metrics-endpoint"`
}

var metricsServer = gx.ProvideNamed(
	func(p MetricsServerParams) (*server.HttpServer, error) {
		reg := prometheus.NewRegistry()
		if err := reg.Register(p.PC); err != nil {
			return nil, err
		}

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

		return server.New(p.Endpoint, h)
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
	metricsStorage,
	metricsCollector,
	metricsServer,
	metricsServerJob,
)
