package encodingload

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "eigenda"
const subsystem = "encoding_load_generator"

// encodingLoadGeneratorMetrics encapsulates the metrics for the encoding load generator.
type encodingLoadGeneratorMetrics struct {
	operationsInFlight *prometheus.GaugeVec
	// TODO: Add more metrics specific to encoding operations
}

// newEncodingLoadGeneratorMetrics creates a new encodingLoadGeneratorMetrics.
func newEncodingLoadGeneratorMetrics(registry *prometheus.Registry) *encodingLoadGeneratorMetrics {
	operationsInFlight := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "operations_in_flight",
			Help:      "Number of encoding operations in flight",
		},
		[]string{},
	)

	return &encodingLoadGeneratorMetrics{
		operationsInFlight: operationsInFlight,
	}
}

// startOperation should be called when starting the process of encoding a blob
func (m *encodingLoadGeneratorMetrics) startOperation() {
	m.operationsInFlight.WithLabelValues().Inc()
}

// endOperation should be called when finishing the process of encoding a blob
func (m *encodingLoadGeneratorMetrics) endOperation() {
	m.operationsInFlight.WithLabelValues().Dec()
}
