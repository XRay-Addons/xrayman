package metrics

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Collector struct {
	s       Storage
	timeout time.Duration
	log     *zap.Logger
}

var _ prometheus.Collector = (*Collector)(nil)

func New(s Storage, log *zap.Logger) (*Collector, error) {
	if s == nil {
		return nil, xerr.NilArg("s")
	}
	if log == nil {
		return nil, xerr.NilArg("log")
	}
	return &Collector{
		s:       s,
		timeout: 1 * time.Second,
		log:     log,
	}, nil
}

func (c *Collector) Close() {
}

const namespace = "nodeman"

var labelNames = []string{"node_id", "endpoint"}

var (
	uplinkDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "upload"),
		"Total outbound traffic for the node",
		labelNames, nil,
	)
	downlinkDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "download"),
		"Total inbound traffic for the node",
		labelNames, nil,
	)
	openConnectionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "open_connections"),
		"Open connections on node",
		labelNames, nil,
	)
	cpuLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "cpu_load"),
		"CPU load of the node",
		labelNames, nil,
	)
	ramLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "ram_load"),
		"RAM load of the node",
		labelNames, nil,
	)
	memLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "mem_load"),
		"Memory load of the node",
		labelNames, nil,
	)
	scrapeErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "scrape_error"),
		"1 if the last scrape of node metrics failed, 0 otherwise",
		nil, nil,
	)
)

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- uplinkDesc
	ch <- downlinkDesc
	ch <- openConnectionsDesc
	ch <- cpuLoadDesc
	ch <- ramLoadDesc
	ch <- memLoadDesc
	ch <- scrapeErrorDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	m, err := c.s.GetNodeMetrics(ctx)
	if err != nil {
		c.log.Warn(fmt.Sprintf("metrics: failed to get node metrics: %v", err))
		c.sendConstMetric(ch, scrapeErrorDesc, prometheus.GaugeValue, 1)
		return
	}
	c.sendConstMetric(ch, scrapeErrorDesc, prometheus.GaugeValue, 0)

	for _, nm := range m {
		id := strconv.Itoa(nm.ID)
		endpoint := nm.Endpoint

		c.sendConstMetric(ch, uplinkDesc, prometheus.CounterValue, float64(nm.Traffic.Upload), id, endpoint)
		c.sendConstMetric(ch, downlinkDesc, prometheus.CounterValue, float64(nm.Traffic.Download), id, endpoint)
		c.sendConstMetric(ch, openConnectionsDesc, prometheus.GaugeValue, float64(nm.Perf.OpenConnections), id, endpoint)
		c.sendConstMetric(ch, cpuLoadDesc, prometheus.GaugeValue, float64(nm.Perf.CpuLoad), id, endpoint)
		c.sendConstMetric(ch, ramLoadDesc, prometheus.GaugeValue, float64(nm.Perf.RamLoad), id, endpoint)
		c.sendConstMetric(ch, memLoadDesc, prometheus.GaugeValue, float64(nm.Perf.MemLoad), id, endpoint)
	}
}

func (c *Collector) sendConstMetric(
	ch chan<- prometheus.Metric,
	desc *prometheus.Desc,
	valueType prometheus.ValueType,
	value float64,
	labelValues ...string,
) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		c.log.Warn(fmt.Sprintf("metrics: failed to create metric %v: %v", desc, err))
		return
	}
	ch <- m
}
