package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// StageTimer encapsulates metrics to help track the time spent in each stage of the payload dispersal process.
//
// This object is thread safe.
type StageTimer struct {
	// counts the number of operations in each specific stage
	stageCount *prometheus.GaugeVec
	// tracks the latency for each stage
	stageLatency *prometheus.SummaryVec
	// if true, then history is captured for debugging purposes
	historyEnabled bool
}

// SequenceProbe can be used to track the amount of time that a particular operation spends doing particular
// sub-operations (i.e. stages). Multiple instances of a particular operation can be tracked concurrently by the same
// StageTimer. For each operation, the StageTimer builds a SequenceProbe. Each SequenceProbe is responsible for
// tracking the lifecycle of a single  iteration of an operation.
//
// A SequenceProbe is not thread safe. It is intended for use in measuring a linear sequence of operations. Do not call
// SetStage or End from multiple goroutines at the same time.
type SequenceProbe struct {
	// the parent StageTimer
	stageTimer *StageTimer
	// set to true when the SequenceProbe has entered its first stage
	initialized bool
	// the current stage of the operation
	currentStage string
	// the time when the current stage started
	currentStageStart time.Time
	// true after End() is called
	ended bool
	// a history of operations, concatenate and log if a lifecycle description is needed.
	// If nil, then no history is captured.
	history *strings.Builder
}

// NewStageTimer creates a new stageTimer with the given prefix and name.
func NewStageTimer(registry *prometheus.Registry, prefix, name string, historyEnabled bool) *StageTimer {
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
		stageLatency:   stageLatency,
		stageCount:     stageCount,
		historyEnabled: historyEnabled,
	}
}

// NewSequence creates a new sequenceProbe with the given initial stage.
func (s *StageTimer) NewSequence() *SequenceProbe {
	if s == nil {
		return nil
	}
	var history *strings.Builder
	if s.historyEnabled {
		history = &strings.Builder{}
	}
	return &SequenceProbe{
		stageTimer: s,
		history:    history,
	}
}

// SetStage updates the stage of the current sequence. This method is a no-op if the new stage is the same as
// the current stage or if the sequenceProbe has already ended.
func (p *SequenceProbe) SetStage(stage string) {
	if p == nil || p.ended || p.currentStage == stage {
		return
	}

	now := time.Now()
	p.stageTimer.stageCount.WithLabelValues(stage).Inc()

	if !p.initialized {
		// First stage setup
		p.currentStage = stage
		p.currentStageStart = now
		if p.history != nil {
			p.history.WriteString(p.currentStage)
		}
		p.initialized = true
		return
	}

	elapsed := ToMilliseconds(now.Sub(p.currentStageStart))
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(elapsed)
	p.currentStageStart = now

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()
	p.currentStage = stage

	if p.history != nil {
		fmt.Fprintf(p.history, ":%0.1f,%s", elapsed, stage)
	}
}

// End completes the current sequence. It is important to call this before discarding the sequenceProbe.
// This method is a no-op if called more than once.
func (p *SequenceProbe) End() {
	if p == nil || p.ended {
		return
	}
	p.ended = true

	now := time.Now()
	elapsed := ToMilliseconds(now.Sub(p.currentStageStart))
	p.stageTimer.stageLatency.WithLabelValues(p.currentStage).Observe(elapsed)

	p.stageTimer.stageCount.WithLabelValues(p.currentStage).Dec()

	if p.history != nil {
		fmt.Fprintf(p.history, ":%0.1f", elapsed)
	}
}

// History returns a string representation of the history of stages for this sequenceProbe. Useful for debugging
// specific executions of a sequence. Format is as follows. The elapsed time is in milliseconds.
//
//	<stage1>:<elapsed1>,<stage2>:<elapsed2>,...,<stageN>:<elapsedN>
func (p *SequenceProbe) History() string {
	if p == nil {
		return ""
	}
	if p.history == nil {
		return ""
	}
	return p.history.String()
}
