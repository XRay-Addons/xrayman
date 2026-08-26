package metrics

import (
	"context"
	"fmt"
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

var (
	totalInboundDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "total_inbound"),
		"Total inbound traffic/connections for the node",
		[]string{"endpoint"}, nil,
	)
	totalOutboundDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "total_outbound"),
		"Total outbound traffic/connections for the node",
		[]string{"endpoint"}, nil,
	)
	cpuLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "cpu_load"),
		"CPU load of the node",
		[]string{"endpoint"}, nil,
	)
	ramLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "ram_load"),
		"RAM load of the node",
		[]string{"endpoint"}, nil,
	)
	memLoadDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "mem_load"),
		"Memory load of the node",
		[]string{"endpoint"}, nil,
	)
	scrapeErrorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "node", "scrape_error"),
		"1 if the last scrape of node metrics failed, 0 otherwise",
		nil, nil,
	)
)

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- totalInboundDesc
	ch <- totalOutboundDesc
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
		endpoint := nm.Endpoint
		c.sendConstMetric(ch, totalInboundDesc, prometheus.CounterValue, float64(nm.TotalInbount), endpoint)
		c.sendConstMetric(ch, totalOutboundDesc, prometheus.CounterValue, float64(nm.TotalOutbound), endpoint)
		c.sendConstMetric(ch, cpuLoadDesc, prometheus.GaugeValue, float64(nm.CpuLoad), endpoint)
		c.sendConstMetric(ch, ramLoadDesc, prometheus.GaugeValue, float64(nm.RamLoad), endpoint)
		c.sendConstMetric(ch, memLoadDesc, prometheus.GaugeValue, float64(nm.MemLoad), endpoint)
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
