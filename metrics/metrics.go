package metrics

import (
	"fmt"
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
	secondarySubsystem  = "secondary"
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

	RecordRPCServerRequest(method string) func(status string, mode string, ver string)
	RecordSecondaryRequest(bt string, method string) func(status string)

	Document() []metrics.DocumentedMetric
}

// Metrics ... Metrics struct
type Metrics struct {
	Info *prometheus.GaugeVec
	Up   prometheus.Gauge

	// server metrics
	HTTPServerRequestsTotal          *prometheus.CounterVec
	HTTPServerBadRequestHeader       *prometheus.CounterVec
	HTTPServerRequestDurationSeconds *prometheus.HistogramVec

	// secondary metrics
	SecondaryRequestsTotal      *prometheus.CounterVec
	SecondaryRequestDurationSec *prometheus.HistogramVec

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
			"method", "status", "commitment_mode", "cert_version",
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
			"method", // no status on histograms because those are very expensive
		}),
		SecondaryRequestsTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: secondarySubsystem,
			Name:      "requests_total",
			Help:      "Total requests to the secondary storage",
		}, []string{
			"backend_type", "method", "status",
		}),
		SecondaryRequestDurationSec: factory.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: secondarySubsystem,
			Name:      "request_duration_seconds",
			Buckets:   prometheus.ExponentialBucketsRange(0.05, 1200, 20),
			Help:      "Histogram of secondary storage request durations",
		}, []string{
			"backend_type",
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
func (m *Metrics) RecordRPCServerRequest(method string) func(status, mode, ver string) {
	// we don't want to track the status code on the histogram because that would
	// create a huge number of labels, and cost a lot on cloud hosted services
	timer := prometheus.NewTimer(m.HTTPServerRequestDurationSeconds.WithLabelValues(method))
	return func(status, mode, ver string) {
		m.HTTPServerRequestsTotal.WithLabelValues(method, status, mode, ver).Inc()
		timer.ObserveDuration()
	}
}

// RecordSecondaryPut records a secondary put/get operation.
func (m *Metrics) RecordSecondaryRequest(bt string, method string) func(status string) {
	timer := prometheus.NewTimer(m.SecondaryRequestDurationSec.WithLabelValues(bt))

	return func(status string) {
		m.SecondaryRequestsTotal.WithLabelValues(bt, method, status).Inc()
		timer.ObserveDuration()
	}
}

// StartServer starts the metrics server on the given hostname and port.
// If port is 0, it automatically assigns an available port and returns the actual port.
func (m *Metrics) StartServer(hostname string, port int) (*ophttp.HTTPServer, error) {
	address := net.JoinHostPort(hostname, strconv.Itoa(port))

	h := promhttp.InstrumentMetricHandler(
		m.registry, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)

	server, err := ophttp.StartHTTPServer(address, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return server, nil
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

func (n *noopMetricer) RecordSecondaryRequest(string, string) func(status string) {
	return func(string) {}
}
