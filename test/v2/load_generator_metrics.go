package v2

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// loadGeneratorMetrics encapsulates the metrics for the load generator.
type loadGeneratorMetrics struct {
	operationsInFlight *prometheus.GaugeVec
	// TODO (cody-littley) count successes, failures, and timeouts
}

// newLoadGeneratorMetrics creates a new loadGeneratorMetrics.0
func newLoadGeneratorMetrics(registry *prometheus.Registry) *loadGeneratorMetrics {
	operationsInFlight := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "operations_in_flight",
			Help:      "Number of operations in flight",
		},
		[]string{},
	)

	return &loadGeneratorMetrics{
		operationsInFlight: operationsInFlight,
	}
}

// startOperation should be called when starting the process of dispersing + verifying a blob
func (m *loadGeneratorMetrics) startOperation() {
	m.operationsInFlight.WithLabelValues().Inc()
}

// endOperation should be called when finishing the process of dispersing + verifying a blob
func (m *loadGeneratorMetrics) endOperation() {
	m.operationsInFlight.WithLabelValues().Dec()
}
