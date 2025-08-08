package metrics

import (
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	accountantSubsystem         = "accountant"
	disperserSubsystem          = "disperser_client"
	relayRetrieverSubsystem     = "relay_retriever"
	validatorRetrieverSubsystem = "validator_retriever"
)

type ClientMetricer interface {
	// accountant
	ReportPaymentUsed(accountID string, method string, symbols uint64)
	ReportOnDemandPayment(accountID string, amount *big.Int)
	ReportCumulativePayment(accountID string, amount *big.Int)
	ReportPaymentFailure(accountID string, reason string)
	// disperser_client
	ReportBlobDispersal(method string, blobSize uint64)
	ReportValidationError(errorType string)
	ReportBlobKeyVerificationError()
	// relay_retriever
	ReportPayloadRetrieval(status, method string, payloadSize uint64, relayAttempts int)
	ReportRelayRequest(relayKey string, status string)
	ReportCommitmentVerificationFailure(relayKey string)
	ReportBlobDeserializationFailure(relayKey string)
	ReportPayloadDecodeFailure(blobKey string)
	ReportRelayTimeout(relayKey string)
	ReportRelayConnectionError(relayKey string)
	// validator_retriever
	ReportValidatorPayloadRetrieval(status, method string, payloadSize uint64, quorumAttempts int)
	ReportValidatorQuorumRequest(quorumId string, status string)
	ReportValidatorCommitmentVerificationFailure(quorumId string)
	ReportValidatorBlobDeserializationFailure(quorumId string)
	ReportValidatorPayloadDecodeFailure(blobKey string)
	ReportValidatorTimeout(quorumId string)
	ReportValidatorConnectionError(quorumId string)
}

type ClientMetrics struct {
	// accountant
	PaymentMethodUsed *prometheus.CounterVec
	OnDemandPayment   *prometheus.CounterVec
	CumulativePayment *prometheus.GaugeVec
	PaymentFailures   *prometheus.CounterVec
	// disperser_client
	BlobDispersalRequests     *prometheus.CounterVec
	BlobDispersalDuration     *prometheus.HistogramVec
	BlobSize                  *prometheus.HistogramVec
	GrpcRequests              *prometheus.CounterVec
	GrpcRequestDuration       *prometheus.HistogramVec
	ValidationErrors          *prometheus.CounterVec
	BlobKeyVerificationErrors *prometheus.CounterVec
	ConnectionErrors          *prometheus.CounterVec
	// relay_retriever
	PayloadRetrievalRequests       *prometheus.CounterVec
	PayloadRetrievalDuration       *prometheus.HistogramVec
	PayloadSize                    *prometheus.HistogramVec
	RelayAttemptsPerRetrieval      *prometheus.HistogramVec
	RelayRequests                  *prometheus.CounterVec
	RelayRequestDuration           *prometheus.HistogramVec
	CommitmentVerificationFailures *prometheus.CounterVec
	BlobDeserializationFailures    *prometheus.CounterVec
	PayloadDecodeFailures          *prometheus.CounterVec
	RelayTimeouts                  *prometheus.CounterVec
	RelayConnectionErrors          *prometheus.CounterVec
	// validator_retriever
	ValidatorPayloadRetrievalRequests       *prometheus.CounterVec
	ValidatorPayloadRetrievalDuration       *prometheus.HistogramVec
	ValidatorPayloadSize                    *prometheus.HistogramVec
	ValidatorQuorumAttemptsPerRetrieval     *prometheus.HistogramVec
	ValidatorQuorumRequests                 *prometheus.CounterVec
	ValidatorQuorumRequestDuration          *prometheus.HistogramVec
	ValidatorCommitmentVerificationFailures *prometheus.CounterVec
	ValidatorBlobDeserializationFailures    *prometheus.CounterVec
	ValidatorPayloadDecodeFailures          *prometheus.CounterVec
	ValidatorTimeouts                       *prometheus.CounterVec
	ValidatorConnectionErrors               *prometheus.CounterVec
}

func NewClientMetrics(namespace string, factory metrics.Factory) ClientMetrics {
	return ClientMetrics{
		PaymentMethodUsed: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payment_method_used_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total number of payments by method (reservation or on-demand)",
		}, []string{
			"account_id",
			"method",
		}),
		OnDemandPayment: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "on_demand_payment_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total on-demand payments made",
		}, []string{
			"account_id",
		}),
		CumulativePayment: factory.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cumulative_payment",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Current cumulative payment balance",
		}, []string{
			"account_id",
		}),
		PaymentFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payment_failures_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total payment failures by reason",
		}, []string{
			"account_id",
			"reason",
		}),
		BlobDispersalRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "blob_dispersal_requests_total",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Total blob dispersal requests by status and method",
		}, []string{
			"status",
			"method",
		}),
		BlobDispersalDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "blob_dispersal_duration_seconds",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Time taken for blob dispersal operations",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"method",
		}),
		BlobSize: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "blob_size_bytes",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Size distribution of dispersed blobs",
			Buckets:   prometheus.ExponentialBuckets(1024, 2, 20), // 1KB to ~1GB
		}, []string{
			"method",
		}),
		GrpcRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "grpc_requests_total",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Total GRPC requests by method and status",
		}, []string{
			"method",
			"status",
		}),
		GrpcRequestDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "grpc_request_duration_seconds",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "GRPC request latency",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"method",
		}),
		ValidationErrors: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "validation_errors_total",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Total validation errors by type",
		}, []string{
			"error_type",
		}),
		BlobKeyVerificationErrors: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "blob_key_verification_errors_total",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Total blob key verification failures",
		}, []string{}),
		ConnectionErrors: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "connection_errors_total",
			Namespace: namespace,
			Subsystem: disperserSubsystem,
			Help:      "Total connection errors by reason",
		}, []string{
			"reason",
		}),
		PayloadRetrievalRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payload_retrieval_requests_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total payload retrieval requests by status and method",
		}, []string{
			"status",
			"method",
		}),
		PayloadRetrievalDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "payload_retrieval_duration_seconds",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Time taken for payload retrieval operations",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"method",
		}),
		PayloadSize: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "payload_size_bytes",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Size distribution of retrieved payloads",
			Buckets:   prometheus.ExponentialBuckets(1024, 2, 20), // 1KB to ~1GB
		}, []string{
			"method",
		}),
		RelayAttemptsPerRetrieval: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "relay_attempts_per_retrieval",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Number of relay attempts per retrieval operation",
			Buckets:   []float64{1, 2, 3, 4, 5, 8, 10, 15, 20}, // reasonable for relay count
		}, []string{
			"method",
		}),
		RelayRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "relay_requests_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total requests to individual relays by relay key and status",
		}, []string{
			"relay_key",
			"status",
		}),
		RelayRequestDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "relay_request_duration_seconds",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Per-relay request duration",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"relay_key",
		}),
		CommitmentVerificationFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "commitment_verification_failures_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total commitment verification failures by relay key",
		}, []string{
			"relay_key",
		}),
		BlobDeserializationFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "blob_deserialization_failures_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total blob deserialization failures by relay key",
		}, []string{
			"relay_key",
		}),
		PayloadDecodeFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payload_decode_failures_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total payload decode failures by blob key",
		}, []string{
			"blob_key",
		}),
		RelayTimeouts: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "relay_timeouts_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total relay timeout errors by relay key",
		}, []string{
			"relay_key",
		}),
		RelayConnectionErrors: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "relay_connection_errors_total",
			Namespace: namespace,
			Subsystem: relayRetrieverSubsystem,
			Help:      "Total relay connection errors by relay key",
		}, []string{
			"relay_key",
		}),
		ValidatorPayloadRetrievalRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payload_retrieval_requests_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total validator payload retrieval requests by status and method",
		}, []string{
			"status",
			"method",
		}),
		ValidatorPayloadRetrievalDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "payload_retrieval_duration_seconds",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Time taken for validator payload retrieval operations",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"method",
		}),
		ValidatorPayloadSize: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "payload_size_bytes",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Size distribution of validator retrieved payloads",
			Buckets:   prometheus.ExponentialBuckets(1024, 2, 20), // 1KB to ~1GB
		}, []string{
			"method",
		}),
		ValidatorQuorumAttemptsPerRetrieval: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "quorum_attempts_per_retrieval",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Number of quorum attempts per validator retrieval operation",
			Buckets:   []float64{1, 2, 3, 4, 5, 8, 10}, // reasonable for quorum count
		}, []string{
			"method",
		}),
		ValidatorQuorumRequests: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "quorum_requests_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total requests to individual quorums by quorum ID and status",
		}, []string{
			"quorum_id",
			"status",
		}),
		ValidatorQuorumRequestDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "quorum_request_duration_seconds",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Per-quorum request duration in validator retrieval",
			Buckets:   prometheus.DefBuckets,
		}, []string{
			"quorum_id",
		}),
		ValidatorCommitmentVerificationFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "commitment_verification_failures_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total commitment verification failures by quorum ID",
		}, []string{
			"quorum_id",
		}),
		ValidatorBlobDeserializationFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "blob_deserialization_failures_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total blob deserialization failures by quorum ID",
		}, []string{
			"quorum_id",
		}),
		ValidatorPayloadDecodeFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payload_decode_failures_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total payload decode failures by blob key",
		}, []string{
			"blob_key",
		}),
		ValidatorTimeouts: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "timeouts_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total validator timeout errors by quorum ID",
		}, []string{
			"quorum_id",
		}),
		ValidatorConnectionErrors: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "connection_errors_total",
			Namespace: namespace,
			Subsystem: validatorRetrieverSubsystem,
			Help:      "Total validator connection errors by quorum ID",
		}, []string{
			"quorum_id",
		}),
	}
}

func (m *ClientMetrics) ReportPaymentUsed(accountID string, method string, symbols uint64) {
	m.PaymentMethodUsed.WithLabelValues(accountID, method).Add(float64(symbols))
}

// ReportOnDemandPayment records an on-demand payment amount in wei
func (m *ClientMetrics) ReportOnDemandPayment(accountID string, amount *big.Int) {
	// Convert big.Int to float64 for prometheus.
	// prometheus/client_golang expects a float64, so we lose precision.
	weiFloat, _ := amount.Float64()
	m.OnDemandPayment.WithLabelValues(accountID).Add(weiFloat)
}

func (m *ClientMetrics) ReportCumulativePayment(accountID string, amount *big.Int) {
	// Convert big.Int to float64 for prometheus
	// prometheus/client_golang expects a float64, so we lose precision.
	weiFloat, _ := amount.Float64()
	m.CumulativePayment.WithLabelValues(accountID).Set(weiFloat)
}

func (m *ClientMetrics) ReportPaymentFailure(accountID string, reason string) {
	m.PaymentFailures.WithLabelValues(accountID, reason).Inc()
}

func (m *ClientMetrics) ReportBlobDispersal(method string, blobSize uint64) {
	m.BlobSize.WithLabelValues(method).Observe(float64(blobSize))
}

func (m *ClientMetrics) ReportValidationError(errorType string) {
	m.ValidationErrors.WithLabelValues(errorType).Inc()
}

func (m *ClientMetrics) ReportBlobKeyVerificationError() {
	m.BlobKeyVerificationErrors.WithLabelValues().Inc()
}

func (m *ClientMetrics) ReportPayloadRetrieval(status, method string, payloadSize uint64, relayAttempts int) {
	m.PayloadRetrievalRequests.WithLabelValues(status, method).Inc()
	m.PayloadSize.WithLabelValues(method).Observe(float64(payloadSize))
	m.RelayAttemptsPerRetrieval.WithLabelValues(method).Observe(float64(relayAttempts))
}

func (m *ClientMetrics) ReportRelayRequest(relayKey string, status string) {
	m.RelayRequests.WithLabelValues(relayKey, status).Inc()
}

func (m *ClientMetrics) ReportCommitmentVerificationFailure(relayKey string) {
	m.CommitmentVerificationFailures.WithLabelValues(relayKey).Inc()
}

func (m *ClientMetrics) ReportBlobDeserializationFailure(relayKey string) {
	m.BlobDeserializationFailures.WithLabelValues(relayKey).Inc()
}

func (m *ClientMetrics) ReportPayloadDecodeFailure(blobKey string) {
	m.PayloadDecodeFailures.WithLabelValues(blobKey).Inc()
}

func (m *ClientMetrics) ReportRelayTimeout(relayKey string) {
	m.RelayTimeouts.WithLabelValues(relayKey).Inc()
}

func (m *ClientMetrics) ReportRelayConnectionError(relayKey string) {
	m.RelayConnectionErrors.WithLabelValues(relayKey).Inc()
}

func (m *ClientMetrics) ReportValidatorPayloadRetrieval(status, method string, payloadSize uint64, quorumAttempts int) {
	m.ValidatorPayloadRetrievalRequests.WithLabelValues(status, method).Inc()
	m.ValidatorPayloadSize.WithLabelValues(method).Observe(float64(payloadSize))
	m.ValidatorQuorumAttemptsPerRetrieval.WithLabelValues(method).Observe(float64(quorumAttempts))
}

func (m *ClientMetrics) ReportValidatorQuorumRequest(quorumId string, status string) {
	m.ValidatorQuorumRequests.WithLabelValues(quorumId, status).Inc()
}

func (m *ClientMetrics) ReportValidatorCommitmentVerificationFailure(quorumId string) {
	m.ValidatorCommitmentVerificationFailures.WithLabelValues(quorumId).Inc()
}

func (m *ClientMetrics) ReportValidatorBlobDeserializationFailure(quorumId string) {
	m.ValidatorBlobDeserializationFailures.WithLabelValues(quorumId).Inc()
}

func (m *ClientMetrics) ReportValidatorPayloadDecodeFailure(blobKey string) {
	m.ValidatorPayloadDecodeFailures.WithLabelValues(blobKey).Inc()
}

func (m *ClientMetrics) ReportValidatorTimeout(quorumId string) {
	m.ValidatorTimeouts.WithLabelValues(quorumId).Inc()
}

func (m *ClientMetrics) ReportValidatorConnectionError(quorumId string) {
	m.ValidatorConnectionErrors.WithLabelValues(quorumId).Inc()
}

type NoopAccountantMetricer struct {
}

var NoopAccountantMetrics ClientMetricer = new(NoopAccountantMetricer)

func (m *NoopAccountantMetricer) ReportPaymentUsed(_ string, _ string, _ uint64) {
}

func (m *NoopAccountantMetricer) ReportOnDemandPayment(_ string, _ *big.Int) {
}

func (m *NoopAccountantMetricer) ReportCumulativePayment(_ string, _ *big.Int) {
}

func (m *NoopAccountantMetricer) ReportPaymentFailure(_ string, _ string) {
}

func (m *NoopAccountantMetricer) ReportBlobDispersal(_ string, _ uint64) {
}

func (m *NoopAccountantMetricer) ReportValidationError(_ string) {
}

func (m *NoopAccountantMetricer) ReportBlobKeyVerificationError() {
}

func (m *NoopAccountantMetricer) ReportPayloadRetrieval(_ string, _ string, _ uint64, _ int) {
}

func (m *NoopAccountantMetricer) ReportRelayRequest(_ string, _ string) {
}

func (m *NoopAccountantMetricer) ReportCommitmentVerificationFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportBlobDeserializationFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportPayloadDecodeFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportRelayTimeout(_ string) {
}

func (m *NoopAccountantMetricer) ReportRelayConnectionError(_ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorPayloadRetrieval(_ string, _ string, _ uint64, _ int) {
}

func (m *NoopAccountantMetricer) ReportValidatorQuorumRequest(_ string, _ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorCommitmentVerificationFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorBlobDeserializationFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorPayloadDecodeFailure(_ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorTimeout(_ string) {
}

func (m *NoopAccountantMetricer) ReportValidatorConnectionError(_ string) {
}
