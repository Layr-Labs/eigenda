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
	namespace           = "eigenda_proxy"
	httpServerSubsystem = "http_server"
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
	RecordRPCServerRequest(method string) func(status string, commitmentMode string, version string)

	Document() []metrics.DocumentedMetric
}

// Metrics ... Metrics struct
type Metrics struct {
	Info *prometheus.GaugeVec
	Up   prometheus.Gauge

	HTTPServerRequestsTotal          *prometheus.CounterVec
	HTTPServerBadRequestHeader       *prometheus.CounterVec
	HTTPServerRequestDurationSeconds *prometheus.HistogramVec

	registry *prometheus.Registry
	factory  metrics.Factory
}

var _ Metricer = (*Metrics)(nil)

func NewMetrics(subsystem string) *Metrics {
	if subsystem == "" {
		subsystem = "default"
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	factory := metrics.With(registry)

	return &Metrics{
		Up: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "up",
			Help:      "1 if the proxy server has finished starting up",
		}),
		Info: factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "info",
			Help:      "Pseudo-metric tracking version and config info",
		}, []string{
			"version",
		}),
		HTTPServerRequestsTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: httpServerSubsystem,
			Name:      "requests_total",
			Help:      "Total requests to the HTTP server",
		}, []string{
			"method", "status", "commitment_mode", "DA_cert_version",
		}),
		HTTPServerBadRequestHeader: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: httpServerSubsystem,
			Name:      "requests_bad_header_total",
			Help:      "Total requests to the HTTP server with bad headers",
		}, []string{
			"method", "error_type",
		}),
		HTTPServerRequestDurationSeconds: factory.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: httpServerSubsystem,
			Name:      "request_duration_seconds",
			// TODO: we might want different buckets for different routes?
			// also probably different buckets depending on the backend (memstore, s3, and eigenda have different latencies)
			Buckets: prometheus.ExponentialBucketsRange(0.05, 1200, 20),
			Help:    "Histogram of HTTP server request durations",
		}, []string{
			"method", "commitment_mode", "DA_cert_version", // no status on histograms because those are very expensive
		}),
		registry: registry,
		factory:  factory,
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

// RecordRPCServerRequest is a helper method to record an incoming HTTP request.
// It bumps the requests metric, and tracks how long it takes to serve a response,
// including the HTTP status code.
func (m *Metrics) RecordRPCServerRequest(method string) func(status string, mode string, ver string) {
	// we don't want to track the status code on the histogram because that would
	// create a huge number of labels, and cost a lot on cloud hosted services
	timer := prometheus.NewTimer(m.HTTPServerRequestDurationSeconds.WithLabelValues(method))
	return func(status, mode, ver string) {
		m.HTTPServerRequestsTotal.WithLabelValues(method, status, mode, ver).Inc()
		timer.ObserveDuration()
	}
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
}

var NoopMetrics Metricer = new(noopMetricer)

func (n *noopMetricer) Document() []metrics.DocumentedMetric {
	return nil
}

func (n *noopMetricer) RecordInfo(_ string) {
}

func (n *noopMetricer) RecordUp() {
}

func (n *noopMetricer) RecordRPCServerRequest(string) func(status, mode, ver string) {
	return func(string, string, string) {}
}
