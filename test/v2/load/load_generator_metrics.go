package load

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "eigenda_test_client"

// loadGeneratorMetrics encapsulates the metrics for the load generator.
type loadGeneratorMetrics struct {
	operationsInFlight     *prometheus.GaugeVec
	dispersalSuccesses     *prometheus.CounterVec
	dispersalFailures      *prometheus.CounterVec
	relayReadSuccesses     *prometheus.CounterVec
	relayReadFailures      *prometheus.CounterVec
	validatorReadSuccesses *prometheus.CounterVec
	validatorReadFailures  *prometheus.CounterVec
}

// newLoadGeneratorMetrics creates a new loadGeneratorMetrics.
func newLoadGeneratorMetrics(registry *prometheus.Registry) *loadGeneratorMetrics {
	// This workaround is needed because of the bug-prone API promauto provides.
	// See https://github.com/prometheus/client_golang/issues/1830 for more details.
	var registerer prometheus.Registerer
	if registry != nil {
		registerer = registry
	}
	operationsInFlight := promauto.With(registerer).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "operations_in_flight",
			Help:      "Number of operations in flight",
		},
		[]string{"operation"},
	)

	dispersalSuccesses := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "dispersal_successes",
			Help:      "Number of successful dispersal operations",
		},
		[]string{},
	)

	dispersalFailures := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,

			Name: "dispersal_failures",
			Help: "Number of failed dispersals",
		},
		[]string{},
	)

	relayReadSuccesses := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "relay_read_successes",
			Help:      "Number of relay read successes",
		},
		[]string{},
	)

	relayReadFailures := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "relay_read_failures",
			Help:      "Number of relay read failures",
		},
		[]string{},
	)

	validatorReadSuccesses := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "validator_read_successes",
			Help:      "Number of validator read successes",
		},
		[]string{},
	)

	validatorReadFailures := promauto.With(registerer).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "validator_read_failures",
			Help:      "Number of validator read failures",
		},
		[]string{},
	)

	return &loadGeneratorMetrics{
		operationsInFlight:     operationsInFlight,
		dispersalSuccesses:     dispersalSuccesses,
		dispersalFailures:      dispersalFailures,
		relayReadSuccesses:     relayReadSuccesses,
		relayReadFailures:      relayReadFailures,
		validatorReadSuccesses: validatorReadSuccesses,
		validatorReadFailures:  validatorReadFailures,
	}
}

// startOperation should be called when starting the process of dispersing + verifying a blob
func (m *loadGeneratorMetrics) startOperation(operation string) {
	m.operationsInFlight.WithLabelValues(operation).Inc()
}

// endOperation should be called when finishing the process of dispersing + verifying a blob
func (m *loadGeneratorMetrics) endOperation(operation string) {
	m.operationsInFlight.WithLabelValues(operation).Dec()
}

func (m *loadGeneratorMetrics) reportDispersalSuccess() {
	m.dispersalSuccesses.WithLabelValues().Inc()
}

func (m *loadGeneratorMetrics) reportDispersalFailure() {
	m.dispersalFailures.WithLabelValues().Inc()
}

func (m *loadGeneratorMetrics) reportRelayReadSuccess() {
	m.relayReadSuccesses.WithLabelValues().Inc()
}

func (m *loadGeneratorMetrics) reportRelayReadFailure() {
	m.relayReadFailures.WithLabelValues().Inc()
}

func (m *loadGeneratorMetrics) reportValidatorReadSuccess() {
	m.validatorReadSuccesses.WithLabelValues().Inc()
}

func (m *loadGeneratorMetrics) reportValidatorReadFailure() {
	m.validatorReadFailures.WithLabelValues().Inc()
}
