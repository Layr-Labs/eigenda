package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "eigenda_test_client"

// testClientMetrics encapsulates the metrics for the test client.
type testClientMetrics struct {
	logger   logging.Logger
	server   *http.Server
	registry *prometheus.Registry

	dispersalTime     *prometheus.SummaryVec
	relayReadTime     *prometheus.SummaryVec
	validatorReadTime *prometheus.SummaryVec
	proxyReadTime     *prometheus.SummaryVec

	operationsInFlight     *prometheus.GaugeVec
	dispersalSuccesses     *prometheus.CounterVec
	dispersalFailures      *prometheus.CounterVec
	relayReadSuccesses     *prometheus.CounterVec
	relayReadFailures      *prometheus.CounterVec
	validatorReadSuccesses *prometheus.CounterVec
	validatorReadFailures  *prometheus.CounterVec
	proxyReadSuccesses     *prometheus.CounterVec
	proxyReadFailures      *prometheus.CounterVec
}

// newTestClientMetrics creates a new testClientMetrics.
func newTestClientMetrics(logger logging.Logger, port int) *testClientMetrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", port)
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	dispersalTime := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "dispersal_time_ms",
			Help:      "Time taken to disperse a blob, in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{},
	)

	relayReadTime := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "relay_read_time_ms",
			Help:      "Time taken to read a blob from a relay, in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"relay_id"},
	)

	validatorReadTime := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "validator_read_time_ms",
			Help:      "Time taken to read a blob from a validator, in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{},
	)

	proxyReadTime := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "proxy_read_time_ms",
			Help:      "Time taken to read a blob from a proxy, in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{},
	)

	operationsInFlight := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "operations_in_flight",
			Help:      "Number of operations in flight",
		},
		[]string{"operation"},
	)

	dispersalSuccesses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "dispersal_successes",
			Help:      "Number of successful dispersal operations",
		},
		[]string{},
	)

	dispersalFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,

			Name: "dispersal_failures",
			Help: "Number of failed dispersals",
		},
		[]string{},
	)

	relayReadSuccesses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "relay_read_successes",
			Help:      "Number of relay read successes",
		},
		[]string{},
	)

	relayReadFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "relay_read_failures",
			Help:      "Number of relay read failures",
		},
		[]string{},
	)

	validatorReadSuccesses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "validator_read_successes",
			Help:      "Number of validator read successes",
		},
		[]string{},
	)

	validatorReadFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "validator_read_failures",
			Help:      "Number of validator read failures",
		},
		[]string{},
	)

	proxyReadSuccesses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "proxy_read_successes",
			Help:      "Number of proxy read successes",
		},
		[]string{},
	)

	proxyReadFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "proxy_read_failures",
			Help:      "Number of proxy read failures",
		},
		[]string{},
	)

	return &testClientMetrics{
		logger:                 logger,
		server:                 server,
		registry:               registry,
		dispersalTime:          dispersalTime,
		relayReadTime:          relayReadTime,
		validatorReadTime:      validatorReadTime,
		proxyReadTime:          proxyReadTime,
		operationsInFlight:     operationsInFlight,
		dispersalSuccesses:     dispersalSuccesses,
		dispersalFailures:      dispersalFailures,
		relayReadSuccesses:     relayReadSuccesses,
		relayReadFailures:      relayReadFailures,
		validatorReadSuccesses: validatorReadSuccesses,
		validatorReadFailures:  validatorReadFailures,
		proxyReadSuccesses:     proxyReadSuccesses,
		proxyReadFailures:      proxyReadFailures,
	}
}

// start starts the metrics server.
func (m *testClientMetrics) start() {
	if m == nil {
		return
	}
	go func() {
		err := m.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.Errorf("failed to start metrics server: %v", err)
		}
	}()
}

// stop stops the metrics server.
func (m *testClientMetrics) stop() {
	if m == nil {
		return
	}
	err := m.server.Close()
	if err != nil {
		m.logger.Errorf("failed to close metrics server: %v", err)
	}
}

func (m *testClientMetrics) reportDispersalTime(duration time.Duration) {
	if m == nil {
		return
	}
	m.dispersalTime.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *testClientMetrics) reportRelayReadTime(duration time.Duration, relayID uint32) {
	if m == nil {
		return
	}
	m.relayReadTime.WithLabelValues(fmt.Sprintf("%d", relayID)).Observe(common.ToMilliseconds(duration))
}

func (m *testClientMetrics) reportValidatorReadTime(duration time.Duration) {
	if m == nil {
		return
	}
	m.validatorReadTime.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *testClientMetrics) reportProxyReadTime(duration time.Duration) {
	if m == nil {
		return
	}
	m.proxyReadTime.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// startOperation should be called when starting the process of dispersing + verifying a blob
func (m *testClientMetrics) startOperation(operation string) {
	if m == nil {
		return
	}
	m.operationsInFlight.WithLabelValues(operation).Inc()
}

// endOperation should be called when finishing the process of dispersing + verifying a blob
func (m *testClientMetrics) endOperation(operation string) {
	if m == nil {
		return
	}
	m.operationsInFlight.WithLabelValues(operation).Dec()
}

func (m *testClientMetrics) reportDispersalSuccess() {
	if m == nil {
		return
	}
	m.dispersalSuccesses.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportDispersalFailure() {
	if m == nil {
		return
	}
	m.dispersalFailures.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportRelayReadSuccess() {
	if m == nil {
		return
	}
	m.relayReadSuccesses.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportRelayReadFailure() {
	if m == nil {
		return
	}
	m.relayReadFailures.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportValidatorReadSuccess() {
	if m == nil {
		return
	}
	m.validatorReadSuccesses.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportValidatorReadFailure() {
	if m == nil {
		return
	}
	m.validatorReadFailures.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportProxyReadSuccess() {
	if m == nil {
		return
	}
	m.proxyReadSuccesses.WithLabelValues().Inc()
}

func (m *testClientMetrics) reportProxyReadFailure() {
	if m == nil {
		return
	}
	m.proxyReadFailures.WithLabelValues().Inc()
}
