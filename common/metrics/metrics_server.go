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
	"os"
	"regexp"
	"slices"
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
func (m *metrics) NewLatencyMetric(
	name string,
	label string,
	description string,
	quantiles ...*Quantile) (LatencyMetric, error) {

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
		metric := newLatencyMetric(name, label, description, nil)
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
			[]string{"label", "unit"},
		)
		m.summaryVecMap[name] = vec
	}

	metric := newLatencyMetric(name, label, description, vec)
	m.metricMap[id] = metric
	return metric, nil
}

// NewCountMetric creates a new CountMetric instance.
func (m *metrics) NewCountMetric(
	name string,
	label string,
	description string) (CountMetric, error) {

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
		metric := newCountMetric(name, label, description, nil)
		m.metricMap[id] = metric
	}

	vec, ok := m.counterVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: m.config.Namespace,
				Name:      name,
			},
			[]string{"label", "unit"},
		)
		m.counterVecMap[name] = vec
	}

	metric := newCountMetric(name, label, description, vec)
	m.metricMap[id] = metric

	return metric, nil
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (m *metrics) NewGaugeMetric(
	name string,
	label string,
	unit string,
	description string) (GaugeMetric, error) {

	m.lock.Lock()
	defer m.lock.Unlock()
	return m.newGaugeMetricUnsafe(name, label, unit, description)
}

// newGaugeMetricUnsafe creates a new GaugeMetric instance without locking.
func (m *metrics) newGaugeMetricUnsafe(
	name string,
	label string,
	unit string,
	description string) (GaugeMetric, error) {

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
		metric := newGaugeMetric(name, label, unit, description, nil)
		m.metricMap[id] = metric
	}

	vec, ok := m.gaugeVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: m.config.Namespace,
				Name:      name,
			},
			[]string{"label", "unit"},
		)
		m.gaugeVecMap[name] = vec
	}

	metric := newGaugeMetric(name, label, unit, description, vec)
	m.metricMap[id] = metric

	return metric, nil
}

func (m *metrics) NewAutoGauge(
	name string,
	label string,
	unit string,
	description string,
	pollPeriod time.Duration,
	source func() float64) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return errors.New("metrics server is not alive")
	}

	gauge, err := m.newGaugeMetricUnsafe(name, label, unit, description)
	if err != nil {
		return err
	}

	if !gauge.Enabled() {
		return nil
	}

	pollingAgent := func() {
		ticker := time.NewTicker(pollPeriod)
		for m.isAlive.Load() {
			value := source()
			gauge.Set(value)
			<-ticker.C
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

func (m *metrics) GenerateMetricsDocumentation() string {
	sb := &strings.Builder{}

	enabledCount := 0
	metricIDs := make([]*metricID, 0, len(m.metricMap))
	for id, metric := range m.metricMap {
		if metric.Enabled() {
			enabledCount++
		}
		boundID := id
		metricIDs = append(metricIDs, &boundID)
	}

	// sort the metric IDs alphabetically
	sortFunc := func(a *metricID, b *metricID) int {
		if a.name != b.name {
			return strings.Compare(a.name, b.name)
		}
		return strings.Compare(a.label, b.label)
	}
	slices.SortFunc(metricIDs, sortFunc)

	sb.Write([]byte(fmt.Sprintf("# Metrics Documentation for namespace '%s'\n\n", m.config.Namespace)))
	sb.Write([]byte("This documentation is automatically generated. " +
		"It reflects the metrics that were registered at the time that this document was generated.\n\n"))

	sb.Write([]byte(fmt.Sprintf("There are %d metrics registered. Of these, %d are enabled.\n\n",
		len(m.metricMap), enabledCount)))

	for _, id := range metricIDs {
		metric := m.metricMap[*id]

		sb.Write([]byte("---\n\n"))
		sb.Write([]byte(fmt.Sprintf("## %s\n\n", id.String())))
		sb.Write([]byte(fmt.Sprintf("%s\n\n", metric.Description())))
		sb.Write([]byte("| Field | Value |\n"))
		sb.Write([]byte("|---|---|\n"))
		sb.Write([]byte(fmt.Sprintf("| **Name** | '%s' |\n", metric.Name())))
		sb.Write([]byte(fmt.Sprintf("| **Label** | '%s' |\n", metric.Label())))
		sb.Write([]byte(fmt.Sprintf("| **Type** | %s |\n", metric.Type())))
		sb.Write([]byte(fmt.Sprintf("| **Unit** | %s |\n", metric.Unit())))
		enabledString := "enabled"
		if !metric.Enabled() {
			enabledString = "disabled"
		}
		sb.Write([]byte(fmt.Sprintf("| **Status** | %s |\n\n", enabledString)))
	}

	return sb.String()
}

func (m *metrics) WriteMetricsDocumentation(fileName string) error {
	doc := m.GenerateMetricsDocumentation()

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	_, err = file.Write([]byte(doc))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %v", err)
	}

	return nil
}
