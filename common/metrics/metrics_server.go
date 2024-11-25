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

	// TODO these maps are not needed any more!

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

	// quantilesMap contains a string describing the quantiles for each latency metric. Used to generate documentation.
	quantilesMap map[metricID]string
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
		quantilesMap:  make(map[metricID]string),
	}
	m.isAlive.Store(true)
	return m
}

// metricID is a unique identifier for a metric.
type metricID struct {
	name string
	unit string
}

var legalCharactersRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// containsLegalCharacters returns true if the string contains only legal characters (alphanumeric and underscore).
func containsLegalCharacters(s string) bool {
	return legalCharactersRegex.MatchString(s)
}

// newMetricID creates a new metricID instance.
func newMetricID(name string, unit string) (metricID, error) {
	if !containsLegalCharacters(name) {
		return metricID{}, fmt.Errorf("invalid metric name: %s", name)
	}
	if !containsLegalCharacters(unit) {
		return metricID{}, fmt.Errorf("invalid metric unit: %s", unit)
	}
	return metricID{
		name: name,
		unit: unit,
	}, nil
}

// NameWithUnit returns the name of the metric with the unit appended.
func (i *metricID) NameWithUnit() string {
	return fmt.Sprintf("%s_%s", i.name, i.unit)
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
	description string,
	labelTemplate any,
	quantiles ...*Quantile) (LatencyMetric, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

	id, err := newMetricID(name, "ms")
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(LatencyMetric), nil
	}

	quantilesString := ""
	objectives := make(map[float64]float64, len(quantiles))
	for i, q := range quantiles {
		objectives[q.Quantile] = q.Error

		quantilesString += fmt.Sprintf("`%.3f`", q.Quantile)
		if i < len(quantiles)-1 {
			quantilesString += ", "
		}
	}
	m.quantilesMap[id] = quantilesString

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec, ok := m.summaryVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  m.config.Namespace,
				Name:       id.NameWithUnit(),
				Objectives: objectives,
			},
			labeler.getKeys(),
		)
		m.summaryVecMap[name] = vec
	}

	metric := newLatencyMetric(name, description, vec, labeler)
	m.metricMap[id] = metric
	return metric, nil
}

// NewCountMetric creates a new CountMetric instance.
func (m *metrics) NewCountMetric(
	name string,
	description string,
	labelTemplate any) (CountMetric, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

	id, err := newMetricID(name, "count")
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(CountMetric), nil
	}

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec, ok := m.counterVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: m.config.Namespace,
				Name:      id.NameWithUnit(),
			},
			labeler.getKeys(),
		)
		m.counterVecMap[name] = vec
	}

	metric := newCountMetric(name, description, vec, labeler)
	m.metricMap[id] = metric

	return metric, nil
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (m *metrics) NewGaugeMetric(
	name string,
	unit string,
	description string,
	labelTemplate any) (GaugeMetric, error) {

	m.lock.Lock()
	defer m.lock.Unlock()
	return m.newGaugeMetricUnsafe(name, unit, description, labelTemplate)
}

// newGaugeMetricUnsafe creates a new GaugeMetric instance without locking.
func (m *metrics) newGaugeMetricUnsafe(
	name string,
	unit string,
	description string,
	labelTemplate any) (GaugeMetric, error) {

	if !m.isAlive.Load() {
		return nil, errors.New("metrics server is not alive")
	}

	id, err := newMetricID(name, unit)
	if err != nil {
		return nil, err
	}

	preExistingMetric, ok := m.metricMap[id]
	if ok {
		return preExistingMetric.(GaugeMetric), nil
	}

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec, ok := m.gaugeVecMap[name]
	if !ok {
		vec = promauto.With(m.registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: m.config.Namespace,
				Name:      id.NameWithUnit(),
			},
			labeler.getKeys(),
		)
		m.gaugeVecMap[name] = vec
	}

	metric := newGaugeMetric(name, unit, description, vec, labeler)
	m.metricMap[id] = metric

	return metric, nil
}

func (m *metrics) NewAutoGauge(
	name string,
	unit string,
	description string,
	pollPeriod time.Duration,
	source func() float64,
	label ...any) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.isAlive.Load() {
		return errors.New("metrics server is not alive")
	}

	if len(label) > 1 {
		return fmt.Errorf("too many labels provided, expected 1, got %d", len(label))
	}
	var l any
	if len(label) == 1 {
		l = label[0]
	}

	gauge, err := m.newGaugeMetricUnsafe(name, unit, description, l)
	if err != nil {
		return err
	}

	pollingAgent := func() {
		ticker := time.NewTicker(pollPeriod)
		for m.isAlive.Load() {
			value := source()

			_ = gauge.Set(value, l)
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

func (m *metrics) GenerateMetricsDocumentation() string {
	sb := &strings.Builder{}

	metricIDs := make([]*metricID, 0, len(m.metricMap))
	for id := range m.metricMap {
		boundID := id
		metricIDs = append(metricIDs, &boundID)
	}

	// sort the metric IDs alphabetically
	sortFunc := func(a *metricID, b *metricID) int {
		if a.name != b.name {
			return strings.Compare(a.name, b.name)
		}
		return strings.Compare(a.unit, b.unit)
	}
	slices.SortFunc(metricIDs, sortFunc)

	sb.Write([]byte(fmt.Sprintf("# Metrics Documentation for namespace '%s'\n\n", m.config.Namespace)))
	sb.Write([]byte(fmt.Sprintf("This documentation was automatically generated at time `%s`\n\n",
		time.Now().Format(time.RFC3339))))

	sb.Write([]byte(fmt.Sprintf("There are a total of `%d` registered metrics.\n\n", len(m.metricMap))))

	for _, id := range metricIDs {
		metric := m.metricMap[*id]

		sb.Write([]byte("---\n\n"))
		sb.Write([]byte(fmt.Sprintf("## %s\n\n", id.NameWithUnit())))

		sb.Write([]byte(fmt.Sprintf("%s\n\n", metric.Description())))

		sb.Write([]byte("|   |   |\n"))
		sb.Write([]byte("|---|---|\n"))
		sb.Write([]byte(fmt.Sprintf("| **Name** | `%s` |\n", metric.Name())))
		sb.Write([]byte(fmt.Sprintf("| **Unit** | `%s` |\n", metric.Unit())))
		labels := metric.LabelFields()
		if len(labels) > 0 {
			sb.Write([]byte("| **Labels** | "))
			for i, label := range labels {
				sb.Write([]byte(fmt.Sprintf("`%s`", label)))
				if i < len(labels)-1 {
					sb.Write([]byte(", "))
				}
			}
			sb.Write([]byte(" |\n"))
		}
		sb.Write([]byte(fmt.Sprintf("| **Type** | `%s` |\n", metric.Type())))
		if metric.Type() == "latency" {
			sb.Write([]byte(fmt.Sprintf("| **Quantiles** | %s |\n", m.quantilesMap[*id])))
		}
		sb.Write([]byte(fmt.Sprintf("| **Fully Qualified Name** | `%s_%s_%s` |\n",
			m.config.Namespace, id.name, id.unit)))
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
