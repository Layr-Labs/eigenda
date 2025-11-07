package reservationvalidation

import (
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Tracks metrics for the [ReservationPaymentValidator]
type ReservationValidatorMetrics struct {
	reservationBytes                 prometheus.Histogram
	reservationSymbolsTotal          prometheus.Counter
	reservationInsufficientBandwidth prometheus.Counter
	reservationQuorumNotPermitted    prometheus.Counter
	reservationTimeOutOfRange        prometheus.Counter
	reservationTimeMovedBackward     prometheus.Counter
	reservationUnexpectedErrors      prometheus.Counter
}

func NewReservationValidatorMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
) *ReservationValidatorMetrics {
	if registry == nil {
		return nil
	}

	bytes := promauto.With(registry).NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "reservation_bytes",
			Subsystem: subsystem,
			Help: "Distribution of byte counts for successful reservation payments. " +
				"Counts reflect actual dispersed bytes, not billed bytes (which may be higher due to min size).",
			// Buckets chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(128*units.KiB, 2, 8),
		},
	)

	symbolsTotal := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_symbols_total",
			Subsystem: subsystem,
			Help: "Total number of symbols validated for successful reservation payments. " +
				"Counts reflect actual dispersed symbols, not billed symbols (which may be higher due to min size).",
		},
	)

	insufficientBandwidth := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_insufficient_bandwidth_count",
			Subsystem: subsystem,
			Help:      "Total number of reservation payments rejected due to insufficient bandwidth",
		},
	)

	quorumNotPermitted := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_quorum_not_permitted_count",
			Subsystem: subsystem,
			Help:      "Total number of reservation payments rejected due to unpermitted quorums",
		},
	)

	timeOutOfRange := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_time_out_of_range_count",
			Subsystem: subsystem,
			Help:      "Total number of reservation payments rejected due to time out of range",
		},
	)

	timeMovedBackward := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_time_moved_backward_count",
			Subsystem: subsystem,
			Help:      "Total number of reservation payments rejected due to time moving backwards",
		},
	)

	unexpectedErrors := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_unexpected_errors_count",
			Subsystem: subsystem,
			Help:      "Total number of unexpected errors during reservation payment authorization",
		},
	)

	return &ReservationValidatorMetrics{
		reservationBytes:                 bytes,
		reservationSymbolsTotal:          symbolsTotal,
		reservationInsufficientBandwidth: insufficientBandwidth,
		reservationQuorumNotPermitted:    quorumNotPermitted,
		reservationTimeOutOfRange:        timeOutOfRange,
		reservationTimeMovedBackward:     timeMovedBackward,
		reservationUnexpectedErrors:      unexpectedErrors,
	}
}

// Records a successful reservation payment
func (m *ReservationValidatorMetrics) RecordSuccess(symbolCount uint32) {
	if m == nil {
		return
	}
	m.reservationBytes.Observe(float64(symbolCount) * encoding.BYTES_PER_SYMBOL)
	m.reservationSymbolsTotal.Add(float64(symbolCount))
}

// Increments the counter for when the holder of a reservation lacks bandwidth to perform the dispersal
func (m *ReservationValidatorMetrics) IncrementInsufficientBandwidth() {
	if m == nil {
		return
	}
	m.reservationInsufficientBandwidth.Inc()
}

// Increments the counter for quorum not permitted errors
func (m *ReservationValidatorMetrics) IncrementQuorumNotPermitted() {
	if m == nil {
		return
	}
	m.reservationQuorumNotPermitted.Inc()
}

// Increments the counter for time out of range errors
func (m *ReservationValidatorMetrics) IncrementTimeOutOfRange() {
	if m == nil {
		return
	}
	m.reservationTimeOutOfRange.Inc()
}

// Increments the counter for time moved backward errors
func (m *ReservationValidatorMetrics) IncrementTimeMovedBackward() {
	if m == nil {
		return
	}
	m.reservationTimeMovedBackward.Inc()
}

// Increments the counter for unexpected errors
func (m *ReservationValidatorMetrics) IncrementUnexpectedErrors() {
	if m == nil {
		return
	}
	m.reservationUnexpectedErrors.Inc()
}
