package signingrate

import (
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/prometheus/client_golang/prometheus"
)

// Encapsulates metrics about the signing rate of validators.
type SigningRateMetrics struct {
	// TODO
}

// NewSigningRateMetrics creates a new SigningRateMetrics instance.
func NewSigningRateMetrics(registry *prometheus.Registry) *SigningRateMetrics {
	if registry == nil {
		return nil
	}

	return &SigningRateMetrics{}
}

// Report a successful signing event for a validator.
func (s *SigningRateMetrics) ReportSuccess(
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration) {

	if s == nil {
		return
	}

	// TODO

}

// Report a failed signing event for a validator.
func (s *SigningRateMetrics) ReportFailure(
	id core.OperatorID,
	batchSize uint64) {

	if s == nil {
		return
	}

	// TODO

}
