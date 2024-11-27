package churner

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"google.golang.org/grpc/codes"
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
var statusCodeMap = map[FailReason]string{
	FailReasonRateLimitExceeded:           codes.ResourceExhausted.String(),
	FailReasonInsufficientStakeToRegister: codes.InvalidArgument.String(),
	FailReasonInsufficientStakeToChurn:    codes.InvalidArgument.String(),
	FailReasonQuorumIdOutOfRange:          codes.InvalidArgument.String(),
	FailReasonPrevApprovalNotExpired:      codes.ResourceExhausted.String(),
	FailReasonInvalidSignature:            codes.InvalidArgument.String(),
	FailReasonProcessChurnRequestFailed:   codes.Internal.String(),
	FailReasonInvalidRequest:              codes.InvalidArgument.String(),
}

type MetricsConfig struct {
	HTTPPort      int
	EnableMetrics bool
}

type Metrics struct {
	metricsServer metrics.Metrics

	numRequests metrics.CountMetric
	latency     metrics.LatencyMetric

	logger logging.Logger
}

type latencyLabel struct {
	method string
}

type numRequestsLabel struct {
	status string
	method string
	reason string
}

func NewMetrics(httpPort int, logger logging.Logger) (*Metrics, error) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metricsServer := metrics.NewMetrics(logger, "eigenda_churner", httpPort)

	numRequests, err := metricsServer.NewCountMetric(
		"request",
		"the number of requests",
		numRequestsLabel{})
	if err != nil {
		return nil, err
	}

	latency, err := metricsServer.NewLatencyMetric(
		"latency",
		"latency summary in milliseconds",
		latencyLabel{},
		&metrics.Quantile{Quantile: 0.5, Error: 0.05},
		&metrics.Quantile{Quantile: 0.9, Error: 0.01},
		&metrics.Quantile{Quantile: 0.95, Error: 0.01},
		&metrics.Quantile{Quantile: 0.99, Error: 0.001})
	if err != nil {
		return nil, err
	}

	return &Metrics{
		metricsServer: metricsServer,
		numRequests:   numRequests,
		latency:       latency,
		logger:        logger.With("component", "ChurnerMetrics"),
	}, nil
}

// WriteMetricsDocumentation writes the metrics for the churner to a markdown file.
func (g *Metrics) WriteMetricsDocumentation() error {
	return g.metricsServer.WriteMetricsDocumentation("operators/churner/mdoc/churner-metrics.md")
}

// ObserveLatency observes the latency of a stage
func (g *Metrics) ObserveLatency(method string, latency time.Duration) {
	g.latency.ReportLatency(latency, latencyLabel{method: method})
}

// IncrementSuccessfulRequestNum increments the number of successful requests
func (g *Metrics) IncrementSuccessfulRequestNum(method string) {
	g.numRequests.Increment(numRequestsLabel{status: "success", method: method})
}

// IncrementFailedRequestNum increments the number of failed requests
func (g *Metrics) IncrementFailedRequestNum(method string, reason FailReason) {
	code, ok := statusCodeMap[reason]
	if !ok {
		g.logger.Error("cannot map failure reason to status code", "failure reason", reason)
		// Treat this as an internal server error. This is a conservative approach to
		// handle a negligence of mapping from failure reason to status code.
		code = codes.Internal.String()
	}

	g.numRequests.Increment(numRequestsLabel{status: code, reason: string(reason), method: method})
}

// Start starts the metrics server
func (g *Metrics) Start() error {
	return g.metricsServer.Start()
}
