package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// StageTimer encapsulates metrics to help track the time spent in each stage of the payload dispersal process.
type StageTimer struct {
	stageCount   *prometheus.GaugeVec
	stageLatency *prometheus.SummaryVec
}

// SequenceProbe can be used to track the amount of time that a particular operation spends doing particular
// sub-operations (i.e. stages). Multiple instances of a particular operation can be tracked concurrently by the same
// StageTimer. For each operation, the StageTimer builds a SequenceProbe. Each SequenceProbe the lifecycle of a single
// iteration of an operation.
type SequenceProbe struct {
	stageTimer        *StageTimer
	currentStage      string
	currentStageStart time.Time
}

// NewStageTimer creates a new stageTimer with the given prefix and name.
func NewStageTimer(registry *prometheus.Registry, prefix, name string) *StageTimer {
	if registry == nil {
		return nil
	}

	stageLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  prefix,
			Name:       name + "_stage_latency_ms",
			Help:       "the latency of each type of operation",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"stage"},
	)

	stageCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Name:      name + "_stage_count",
			Help:      "the number of operations with a specific stage",
		},
		[]string{"stage"},
	)

	return &StageTimer{
		stageLatency: stageLatency,
		stageCount:   stageCount,
	}
}

// NewSequence creates a new sequenceProbe with the given initial stage.
func (s *StageTimer) NewSequence(initialStage string) *SequenceProbe {
	if s == nil {
		return nil
	}

	s.stageCount.WithLabelValues(initialStage).Inc()
	return &SequenceProbe{
		stageTimer:        s,
		currentStage:      initialStage,
		currentStageStart: time.Now(),
	}
}

// SetStage updates the stage of the current sequence.
func (p *SequenceProbe) SetStage(stage string) {
	if p == nil {
		return
	}
	if p.currentStage == stage {
		return
	}

	now := time.Now()
	elapsed := now.Sub(p.currentStageStart)
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(ToMilliseconds(elapsed))
	p.currentStageStart = now

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()
	p.stageTimer.stageCount.WithLabelValues(stage).Inc()
	p.currentStage = stage
}

// End completes the current sequence. It is important to call this before discarding the sequenceProbe.
func (p *SequenceProbe) End() {
	if p == nil {
		return
	}

	now := time.Now()
	elapsed := now.Sub(p.currentStageStart)
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(ToMilliseconds(elapsed))

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()
}
