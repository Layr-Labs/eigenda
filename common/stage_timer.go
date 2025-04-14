package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// StageTimer encapsulates metrics to help track the time spent in each stage of the payload dispersal process.
type StageTimer struct {
	// counts the number of operations in each specific stage
	stageCount *prometheus.GaugeVec
	// tracks the latency for each stage
	stageLatency *prometheus.SummaryVec
}

// SequenceProbe can be used to track the amount of time that a particular operation spends doing particular
// sub-operations (i.e. stages). Multiple instances of a particular operation can be tracked concurrently by the same
// StageTimer. For each operation, the StageTimer builds a SequenceProbe. Each SequenceProbe is responsible for
// tracking the lifecycle of a single  iteration of an operation.
type SequenceProbe struct {
	// the parent StageTimer
	stageTimer *StageTimer
	// set to true when the SequenceProbe has entered its first stage
	initialized bool
	// the current stage of the operation
	currentStage string
	// the time when the current stage started
	currentStageStart time.Time
	// a history of operations, concatenate and log if a lifecycle description is needed
	history *strings.Builder
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
func (s *StageTimer) NewSequence() *SequenceProbe {
	if s == nil {
		return nil
	}
	history := &strings.Builder{}
	return &SequenceProbe{
		stageTimer: s,
		history:    history,
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
	p.stageTimer.stageCount.WithLabelValues(stage).Inc()

	if !p.initialized {
		// First stage setup
		p.currentStage = stage
		p.currentStageStart = now
		p.history.WriteString(p.currentStage)
		p.initialized = true
		return
	}

	elapsed := ToMilliseconds(now.Sub(p.currentStageStart))
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(elapsed)
	p.currentStageStart = now

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()
	p.currentStage = stage

	p.history.WriteString(fmt.Sprintf(":%0.1f,%s", elapsed, stage))
}

// End completes the current sequence. It is important to call this before discarding the sequenceProbe.
func (p *SequenceProbe) End() {
	if p == nil {
		return
	}

	now := time.Now()
	elapsed := ToMilliseconds(now.Sub(p.currentStageStart))
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(elapsed)

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()

	p.history.WriteString(fmt.Sprintf(":%0.1f", elapsed))
}

// History returns a string representation of the history of stages for this sequenceProbe. Useful for debugging
// specific executions of a sequence. Format is as follows. The elapsed time is in milliseconds.
//
//	<stage1>:<elapsed1>,<stage2>:<elapsed2>,...,<stageN>:<elapsedN>
func (p *SequenceProbe) History() string {
	if p == nil {
		return ""
	}
	return p.history.String()
}
