package metrics

import (
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	accountantSubsystem = "accountant"
	disperserSubsystem  = "disperser_client"
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
