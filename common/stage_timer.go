package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// It may be helpful to generalize this utility to be used in other parts of the codebase. For a future PR, perhaps.

// StageTimer encapsulates metrics to help track the time spent in each stage of the payload dispersal process.
type StageTimer struct {
	stageCount   *prometheus.GaugeVec
	stageLatency *prometheus.SummaryVec
}

// SequenceProbe tracks the timing of a single sequence of operations. Multiple sequences can be tracked concurrently.
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

	statusLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  prefix,
			Name:       name + "_stage_latency_ms",
			Help:       "the latency of each type of operation",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"stage"},
	)

	statusCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Name:      name + "_stage_count",
			Help:      "the number of operations with a specific status",
		},
		[]string{"stage"},
	)

	return &StageTimer{
		stageLatency: statusLatency,
		stageCount:   statusCount,
	}
}

// NewSequence creates a new sequenceProbe with the given initial status.
func (s *StageTimer) NewSequence(initialStatus string) *SequenceProbe {
	if s == nil {
		return nil
	}

	s.stageCount.WithLabelValues(initialStatus).Inc()
	return &SequenceProbe{
		stageTimer:        s,
		currentStage:      initialStatus,
		currentStageStart: time.Now(),
	}
}

// SetStage updates the status of the current sequence.
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
