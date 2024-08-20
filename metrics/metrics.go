package metrics

import (
	"net"
	"strconv"

	ophttp "github.com/ethereum-optimism/optimism/op-service/httputil"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Namespace = "eigenda_proxy"
)

// Config ... Metrics server configuration
type Config struct {
	Host              string
	Port              int
	Enabled           bool
	ReadHeaderTimeout int
}

// Metricer ... Interface for metrics
type Metricer interface {
	RecordInfo(version string)
	RecordUp()
	RecordRPCServerRequest(method string) func()
	RecordRPCClientResponse(method string, err error)

	Document() []metrics.DocumentedMetric
}

// Metrics ... Metrics struct
type Metrics struct {
	Info *prometheus.GaugeVec
	Up   prometheus.Gauge

	metrics.RPCMetrics

	registry *prometheus.Registry
	factory  metrics.Factory
}

var _ Metricer = (*Metrics)(nil)

func NewMetrics(procName string) *Metrics {
	if procName == "" {
		procName = "default"
	}
	ns := Namespace + "_" + procName

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	factory := metrics.With(registry)

	return &Metrics{
		Up: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "up",
			Help:      "1 if the proxy server has finished starting up",
		}),
		Info: factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "info",
			Help:      "Pseudo-metric tracking version and config info",
		}, []string{
			"version",
		}),
		RPCMetrics: metrics.MakeRPCMetrics(ns, factory),
		registry:   registry,
		factory:    factory,
	}
}

// RecordInfo sets a pseudo-metric that contains versioning and
// config info for the proxy DA node.
func (m *Metrics) RecordInfo(version string) {
	m.Info.WithLabelValues(version).Set(1)
}

// RecordUp sets the up metric to 1.
func (m *Metrics) RecordUp() {
	prometheus.MustRegister()
	m.Up.Set(1)
}

// StartServer starts the metrics server on the given hostname and port.
func (m *Metrics) StartServer(hostname string, port int) (*ophttp.HTTPServer, error) {
	addr := net.JoinHostPort(hostname, strconv.Itoa(port))
	h := promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
	return ophttp.StartHTTPServer(addr, h)
}

func (m *Metrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopMetricer struct {
	metrics.NoopRPCMetrics
}

var NoopMetrics Metricer = new(noopMetricer)

func (n *noopMetricer) Document() []metrics.DocumentedMetric {
	return nil
}

func (n *noopMetricer) RecordInfo(_ string) {
}

func (n *noopMetricer) RecordUp() {
}
