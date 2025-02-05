package client

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

const namespace = "eigenda_test_client"

// testClientMetrics encapsulates the metrics for the test client.
type testClientMetrics struct {
	logger   logging.Logger
	server   *http.Server
	registry *prometheus.Registry

	dispersalTime     *prometheus.SummaryVec
	certificationTime *prometheus.SummaryVec
	relayReadTime     *prometheus.SummaryVec
	validatorReadTime *prometheus.SummaryVec
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

	certificationTime := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "certification_time_ms",
			Help:      "Time taken to certify a blob, in milliseconds",
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
		[]string{"quorum"},
	)

	return &testClientMetrics{
		logger:            logger,
		server:            server,
		registry:          registry,
		dispersalTime:     dispersalTime,
		certificationTime: certificationTime,
		relayReadTime:     relayReadTime,
		validatorReadTime: validatorReadTime,
	}
}

// start starts the metrics server.
func (m *testClientMetrics) start() {
	go func() {
		err := m.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.Errorf("failed to start metrics server: %v", err)
		}
	}()
}

// stop stops the metrics server.
func (m *testClientMetrics) stop() {
	err := m.server.Close()
	if err != nil {
		m.logger.Errorf("failed to close metrics server: %v", err)
	}
}

// reportDispersalTime reports the time taken to disperse a blob.
func (m *testClientMetrics) reportDispersalTime(duration time.Duration) {
	m.dispersalTime.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// reportCertificationTime reports the time taken to certify a blob.
func (m *testClientMetrics) reportCertificationTime(duration time.Duration) {
	m.certificationTime.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// reportRelayReadTime reports the time taken to read a blob from a relay.
func (m *testClientMetrics) reportRelayReadTime(duration time.Duration, relayID uint32) {
	m.relayReadTime.WithLabelValues(fmt.Sprintf("%d", relayID)).Observe(common.ToMilliseconds(duration))
}

// reportValidatorReadTime reports the time taken to read a blob from a validator.
func (m *testClientMetrics) reportValidatorReadTime(duration time.Duration, quorum core.QuorumID) {
	m.validatorReadTime.WithLabelValues(fmt.Sprintf("%d", quorum)).Observe(common.ToMilliseconds(duration))
}
