package churner

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FailReason string

// Note: failure reason constants must be maintained in sync with statusCodeMap.
const (
	FailReasonRateLimitExceeded           FailReason = "rate_limit_exceeded"            // Rate limited: per operator rate limiting
	FailReasonInsufficientStakeToRegister FailReason = "insufficient_stake_to_register" // Operator doesn't have enough stake to be registered
	FailReasonInsufficientStakeToChurn    FailReason = "insufficient_stake_to_churn"    // Operator doesn't have enough stake to be churned
	FailReasonQuorumIdOutOfRange          FailReason = "quorum_id_out_of_range"         // Quorum ID out of range: quorum is not in the range of [0, QuorumCount]
	FailReasonPrevApprovalNotExpired      FailReason = "prev_approval_not_expired"      // Expiry: previous approval hasn't expired
	FailReasonInvalidSignature            FailReason = "invalid_signature"              // Invalid signature: operator's signature is wrong
	FailReasonProcessChurnRequestFailed   FailReason = "failed_process_churn_request"   // Failed to process churn request
	FailReasonInvalidRequest              FailReason = "invalid_request"                // Invalid request: request is malformed
)

// Note: statusCodeMap must be maintained in sync with failure reason constants.
var statusCodeMap map[FailReason]string = map[FailReason]string{
	FailReasonRateLimitExceeded:           "429",
	FailReasonInsufficientStakeToRegister: "400",
	FailReasonInsufficientStakeToChurn:    "400",
	FailReasonQuorumIdOutOfRange:          "400",
	FailReasonPrevApprovalNotExpired:      "429",
	FailReasonInvalidSignature:            "400",
	FailReasonProcessChurnRequestFailed:   "500",
	FailReasonInvalidRequest:              "400",
}

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	registry *prometheus.Registry

	NumRequests *prometheus.CounterVec
	Latency     *prometheus.SummaryVec

	httpPort string
	logger   common.Logger
}

func NewMetrics(httpPort string, logger common.Logger) *Metrics {
	namespace := "eigenda_churner"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		NumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests",
				Help:      "the number of requests",
			},
			[]string{"status", "reason", "method"},
		),
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"method"},
		),
		registry: reg,
		httpPort: httpPort,
		logger:   logger,
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, latencyMs float64) {
	g.Latency.WithLabelValues(method).Observe(latencyMs)
}

// IncrementRequestNum increments the number of successful requests
func (g *Metrics) IncrementRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "total",
		"method": method,
		"reason": "",
	}).Inc()
}

// IncrementSuccessfulRequestNum increments the number of successful requests
func (g *Metrics) IncrementSuccessfulRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "success",
		"method": method,
		"reason": "",
	}).Inc()
}

// IncrementFailedRequestNum increments the number of failed requests
func (g *Metrics) IncrementFailedRequestNum(method string, reason FailReason) {
	code, ok := statusCodeMap[reason]
	if !ok {
		g.logger.Error("cannot map failure reason to status code", "failure reason", reason)
		// Treat this as an internal server error. This is a conservative approach to
		// handle a negligence of mapping from failure reason to status code.
		code = "500"
	}
	g.NumRequests.With(prometheus.Labels{
		"status": code,
		"reason": string(reason),
		"method": method,
	}).Inc()
}

// Start starts the metrics server
func (g *Metrics) Start(ctx context.Context) {
	g.logger.Info("Starting metrics server at ", "port", g.httpPort)
	addr := fmt.Sprintf(":%s", g.httpPort)
	go func() {
		log := g.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			g.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("Prometheus server failed", "err", err)
	}()
}
