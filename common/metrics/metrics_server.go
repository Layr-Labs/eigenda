package metrics

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

	// autoGaugesToStart is a list of functions that will start auto-gauges. If an auto-gauge is created
	// before the metrics server is started, we don't actually start the goroutine until the server is started.
	autoGaugesToStart []func()

	// lock is a lock used to ensure that metrics are not created concurrently.
	lock sync.Mutex

	// started is true if the metrics server has been started.
	started bool

	// isAlize is true if the metrics server has not been stopped.
	isAlive atomic.Bool

	// server is the metrics server
	server *http.Server
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(logger logging.Logger, config *Config) Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", config.HTTPPort)
	addr := fmt.Sprintf(":%d", config.HTTPPort)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{},
	))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	m := &metrics{
		logger:        logger,
		config:        config,
		registry:      reg,
		counterVecMap: make(map[string]*prometheus.CounterVec),
		summaryVecMap: make(map[string]*prometheus.SummaryVec),
		gaugeVecMap:   make(map[string]*prometheus.GaugeVec),
		metricMap:     make(map[metricID]Metric),
		isAlive:       atomic.Bool{},
		server:        server,
	}
	m.isAlive.Store(true)
	return m
}

// metricID is a unique identifier for a metric.
type metricID struct {
	name  string
	label string
}

var legalCharactersRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// newMetricID creates a new metricID instance.
func newMetricID(name string, label string) (metricID, error) {
	if !legalCharactersRegex.MatchString(name) {
		return metricID{}, fmt.Errorf("invalid metric name: %s", name)
	}
	if label != "" && !legalCharactersRegex.MatchString(label) {
		return metricID{}, fmt.Errorf("invalid metric label: %s", label)
	}
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
func (m *metrics) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.started {
		return errors.New("metrics server already started")
	}
	m.started = true

	go func() {
		err := m.server.ListenAndServe()
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			m.logger.Errorf("metrics server error: %v", err)
		}
	}()

	// start the auto-gauges that were created before the server was started
	for _, autoGauge := range m.autoGaugesToStart {
		go autoGauge()
	}

	return nil
}

// Stop stops the metrics server.
func (m *metrics) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.started {
		return errors.New("metrics server not started")
	}

	if !m.isAlive.Load() {
		return errors.New("metrics server already stopped")
	}

	m.isAlive.Store(false)
	return m.server.Close()
}

// NewLatencyMetric creates a new LatencyMetric instance.
func (m *metrics) NewLatencyMetric(name string, label string, quantiles ...*Quantile) (LatencyMetric, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

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
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

	id, err := newMetricID(name, label)
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(CountMetric), nil
	}

	if m.isBlacklisted(id) {
		metric := newLatencyMetric(name, label, nil)
		m.metricMap[id] = metric
	}

	vec, ok := m.counterVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: m.config.Namespace,
				Name:      name,
			},
			[]string{"label"},
		)
		m.counterVecMap[name] = vec
	}

	metric := newCountMetric(name, label, vec)
	m.metricMap[id] = metric

	return metric, nil
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (m *metrics) NewGaugeMetric(name string, label string) (GaugeMetric, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.newGaugeMetricUnsafe(name, label)
}

// newGaugeMetricUnsafe creates a new GaugeMetric instance without locking.
func (m *metrics) newGaugeMetricUnsafe(name string, label string) (GaugeMetric, error) {
	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

	id, err := newMetricID(name, label)
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(GaugeMetric), nil
	}

	if m.isBlacklisted(id) {
		metric := newLatencyMetric(name, label, nil)
		m.metricMap[id] = metric
	}

	vec, ok := m.gaugeVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: m.config.Namespace,
				Name:      name,
			},
			[]string{"label"},
		)
		m.gaugeVecMap[name] = vec
	}

	metric := newGaugeMetric(name, label, vec)
	m.metricMap[id] = metric

	return metric, nil
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

func (m *metrics) NewAutoGauge(name string, label string, pollPeriod time.Duration, source func() float64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return errors.New("metrics server is not alive")
	}

	gauge, err := m.newGaugeMetricUnsafe(name, label)
	if err != nil {
		return err
	}

	pollingAgent := func() {
		for m.isAlive.Load() {
			value := source()
			gauge.Set(value)
			time.Sleep(pollPeriod)
		}
	}

	if m.started {
		// start the polling agent immediately
		go pollingAgent()
	} else {
		// the polling agent will be started when the metrics server is started
		m.autoGaugesToStart = append(m.autoGaugesToStart, pollingAgent)
	}

	return nil
}
