package validator

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
)

// ValidatorClientMetrics encapsulates metrics for the validator client. If nil, then this object becomes a no-op.
// One ValidatorClientMetrics instance can be shared across multiple ValidatorClient instances.
type ValidatorClientMetrics struct {
	stageTimer *common.StageTimer
}

// NewValidatorClientMetrics creates a new ValidatorClientMetrics instance. If a nil registry is provided,
// then this object becomes a no-op.
func NewValidatorClientMetrics(registry *prometheus.Registry) *ValidatorClientMetrics {
	if registry == nil {
		return nil
	}

	stageTimer := common.NewStageTimer(registry, "RetrievalClient", "GetBlob", false)
	return &ValidatorClientMetrics{
		stageTimer: stageTimer,
	}
}

// NewGetBlobProbe creates a new probe for the GetBlob method.
func (m *ValidatorClientMetrics) newGetBlobProbe() *common.SequenceProbe {
	if m == nil {
		return nil
	}
	return m.stageTimer.NewSequence()
}
