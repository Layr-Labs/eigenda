package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"sync"
)

var _ Metrics = &metrics{}

// metrics is a standard implementation of the Metrics interface via prometheus.
type metrics struct {
	// logger is the logger used to log messages.
	logger logging.Logger

	// config is the configuration for the metrics.
	config *Config

	// registry is the prometheus registry used to report metrics.
	registry *prometheus.Registry

	// counterVecMap is a map of metric names to prometheus counter vectors.
	// These are used to create new counter metrics.
	counterVecMap map[string]*prometheus.CounterVec

	// summaryVecMap is a map of metric names to prometheus summary vectors.
	// These are used to create new latency metrics.
	summaryVecMap map[string]*prometheus.SummaryVec

	// gaugeVecMap is a map of metric names to prometheus gauge vectors.
	// These are used to create new gauge metrics.
	gaugeVecMap map[string]*prometheus.GaugeVec

	// A map from metricID to Metric instance. If a metric is requested but that metric
	// already exists, the existing metric will be returned instead of a new one being created.
	metricMap map[metricID]Metric

	// creationLock is a lock used to ensure that metrics are not created concurrently.
	creationLock sync.Mutex
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(logger logging.Logger, config *Config) Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	return &metrics{
		logger:        logger,
		config:        config,
		registry:      reg,
		counterVecMap: make(map[string]*prometheus.CounterVec),
		summaryVecMap: make(map[string]*prometheus.SummaryVec),
		gaugeVecMap:   make(map[string]*prometheus.GaugeVec),
		metricMap:     make(map[metricID]Metric),
	}
}

// metricID is a unique identifier for a metric.
type metricID struct {
	name  string
	label string
}

// newMetricID creates a new metricID instance.
func newMetricID(name string, label string) (metricID, error) {
	// TODO check for illegal characters
	return metricID{
		name:  name,
		label: label,
	}, nil
}

// String returns a string representation of the metricID.
func (i *metricID) String() string {
	if i.label != "" {
		return fmt.Sprintf("%s:%s", i.name, i.label)
	}
	return i.name
}

// Start starts the metrics server.
func (m *metrics) Start() {
	m.creationLock.Lock()
	defer m.creationLock.Unlock()

	m.logger.Infof("Starting metrics server at port %d", m.config.HTTPPort)
	addr := fmt.Sprintf(":%d", m.config.HTTPPort)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			m.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		panic(fmt.Sprintf("Prometheus server failed: %s", err)) // TODO wrong way to handle this
	}()
}

// Stop stops the metrics server.
func (m *metrics) Stop() {
	m.creationLock.Lock()
	defer m.creationLock.Unlock()
	// TODO
}

// NewLatencyMetric creates a new LatencyMetric instance.
func (m *metrics) NewLatencyMetric(name string, label string, quantiles ...*Quantile) (LatencyMetric, error) {
	m.creationLock.Lock()
	defer m.creationLock.Unlock()

	id, err := newMetricID(name, label)
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(LatencyMetric), nil
	}

	if m.isBlacklisted(id) {
		metric := newLatencyMetric(name, label, nil)
		m.metricMap[id] = metric
	}

	objectives := make(map[float64]float64, len(quantiles))
	for _, q := range quantiles {
		objectives[q.Quantile] = q.Error
	}

	vec, ok := m.summaryVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  m.config.Namespace,
				Name:       name,
				Objectives: objectives,
			},
			[]string{"label"},
		)
		m.summaryVecMap[name] = vec
	}

	metric := newLatencyMetric(name, label, vec)
	m.metricMap[id] = metric
	return metric, nil
}

// NewCountMetric creates a new CountMetric instance.
func (m *metrics) NewCountMetric(name string, label string) (CountMetric, error) {
	m.creationLock.Lock()
	defer m.creationLock.Unlock()

	// TODO
	return nil, nil
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (m *metrics) NewGaugeMetric(name string, label string) (GaugeMetric, error) {
	m.creationLock.Lock()
	defer m.creationLock.Unlock()

	// TODO
	return nil, nil
}

// isBlacklisted returns true if the metric name is blacklisted.
func (m *metrics) isBlacklisted(id metricID) bool {
	metric := id.String()

	if m.config.MetricsBlacklist != nil {
		for _, blacklisted := range m.config.MetricsBlacklist {
			if metric == blacklisted {
				return true
			}
		}
	}
	if m.config.MetricsFuzzyBlacklist != nil {
		for _, blacklisted := range m.config.MetricsFuzzyBlacklist {
			if strings.Contains(metric, blacklisted) {
				return true
			}
		}
	}
	return false
}
