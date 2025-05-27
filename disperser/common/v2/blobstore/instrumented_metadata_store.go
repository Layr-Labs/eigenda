package blobstore

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "eigenda"
	subsystem = "metadata_store"
)

// BackendType represents the type of backend storage
type BackendType string

const (
	BackendDynamoDB   BackendType = "dynamodb"
	BackendPostgreSQL BackendType = "postgresql"
	BackendUnknown    BackendType = "unknown"
)

type Config struct {
	ServiceName string
	Backend     BackendType
	Registry    *prometheus.Registry
}

var _ MetadataStore = (*InstrumentedMetadataStore)(nil)

type InstrumentedMetadataStore struct {
	metadataStore MetadataStore
	metrics       *metadataStoreMetricsCollector
	config        Config
}

type metadataStoreMetricsCollector struct {
	// Request latency summary
	requestLatency *prometheus.SummaryVec
	// Request counter
	requestTotal *prometheus.CounterVec
	// Errors counter
	errorTotal *prometheus.CounterVec
	// Concurrent requests gauge
	requestsInFlight *prometheus.GaugeVec
}

func NewInstrumentedMetadataStore(metadataStore MetadataStore, config Config) *InstrumentedMetadataStore {
	if config.Registry == nil {
		config.Registry = prometheus.NewRegistry()
	}

	metrics := &metadataStoreMetricsCollector{
		requestLatency: promauto.With(config.Registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  subsystem,
				Name:       "request_duration_seconds",
				Help:       "Duration of metadata store requests",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"method", "status", "service", "backend"},
		),
		requestTotal: promauto.With(config.Registry).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_total",
				Help:      "Total number of metadata store requests",
			},
			[]string{"method", "status", "service", "backend"},
		),
		errorTotal: promauto.With(config.Registry).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "error_total",
				Help:      "Total number of metadata store errors",
			},
			[]string{"method", "status", "service", "backend"},
		),
		requestsInFlight: promauto.With(config.Registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_in_flight",
				Help:      "Number of metadata store requests currently being processed",
			},
			[]string{"method", "service", "backend"},
		),
	}

	return &InstrumentedMetadataStore{
		metadataStore: metadataStore,
		metrics:       metrics,
	}
}

// Helper function to record metrics
func (m *InstrumentedMetadataStore) recordMetrics(method string, start time.Time, err error) {
	duration := time.Since(start).Seconds()
	status := "success"
	backend := string(m.config.Backend)

	if err != nil {
		status = "error"
		errorType := getErrorType(err)
		m.metrics.errorTotal.WithLabelValues(method, m.config.ServiceName, errorType, backend).Inc()
	}

	m.metrics.requestLatency.WithLabelValues(method, m.config.ServiceName, status, backend).Observe(duration)
	m.metrics.requestTotal.WithLabelValues(method, m.config.ServiceName, status, backend).Inc()
}

// Helper function to track in-flight requests
func (m *InstrumentedMetadataStore) trackInFlight(method string) func() {
	backend := string(m.config.Backend)
	m.metrics.requestsInFlight.WithLabelValues(method, m.config.ServiceName, backend).Inc()
	return func() {
		m.metrics.requestsInFlight.WithLabelValues(method, m.config.ServiceName, backend).Dec()
	}
}

// Helper function to categorize errors
func getErrorType(err error) string {
	if err == nil {
		return "none"
	}
	// Add more specific error type detection based on your error types
	switch err {
	case ErrAlreadyExists:
		return "already_exists"
	case ErrMetadataNotFound:
		return "not_found"
	case ErrInvalidStateTransition:
		return "invalid_state_transition"
	default:
		return "unknown"
	}
}

func (m *InstrumentedMetadataStore) CheckBlobExists(ctx context.Context, blobKey corev2.BlobKey) (bool, error) {
	defer m.trackInFlight("CheckBlobExists")()
	start := time.Now()
	exists, err := m.metadataStore.CheckBlobExists(ctx, blobKey)
	m.recordMetrics("CheckBlobExists", start, err)
	return exists, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobMetadata, error) {
	defer m.trackInFlight("GetBlobMetadata")()
	start := time.Now()
	metadata, err := m.metadataStore.GetBlobMetadata(ctx, blobKey)
	m.recordMetrics("GetBlobMetadata", start, err)
	return metadata, err
}

func (m *InstrumentedMetadataStore) PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	defer m.trackInFlight("PutBlobMetadata")()
	start := time.Now()
	err := m.metadataStore.PutBlobMetadata(ctx, blobMetadata)
	m.recordMetrics("PutBlobMetadata", start, err)
	return err
}

func (m *InstrumentedMetadataStore) UpdateBlobStatus(ctx context.Context, key corev2.BlobKey, status v2.BlobStatus) error {
	defer m.trackInFlight("UpdateBlobStatus")()
	start := time.Now()
	err := m.metadataStore.UpdateBlobStatus(ctx, key, status)
	m.recordMetrics("UpdateBlobStatus", start, err)
	return err
}

func (m *InstrumentedMetadataStore) DeleteBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) error {
	defer m.trackInFlight("DeleteBlobMetadata")()
	start := time.Now()
	err := m.metadataStore.DeleteBlobMetadata(ctx, blobKey)
	m.recordMetrics("DeleteBlobMetadata", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataByAccountID(
	ctx context.Context,
	accountId gethcommon.Address,
	start uint64,
	end uint64,
	limit int,
	ascending bool,
) ([]*v2.BlobMetadata, error) {
	defer m.trackInFlight("GetBlobMetadataByAccountID")()
	startTime := time.Now()
	metadata, err := m.metadataStore.GetBlobMetadataByAccountID(ctx, accountId, start, end, limit, ascending)
	m.recordMetrics("GetBlobMetadataByAccountID", startTime, err)
	return metadata, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error) {
	defer m.trackInFlight("GetBlobMetadataByStatus")()
	start := time.Now()
	metadata, err := m.metadataStore.GetBlobMetadataByStatus(ctx, status, lastUpdatedAt)
	m.recordMetrics("GetBlobMetadataByStatus", start, err)
	return metadata, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataByStatusPaginated(
	ctx context.Context,
	status v2.BlobStatus,
	exclusiveStartKey *StatusIndexCursor,
	limit int32,
) ([]*v2.BlobMetadata, *StatusIndexCursor, error) {
	defer m.trackInFlight("GetBlobMetadataByStatusPaginated")()
	start := time.Now()
	metadata, cursor, err := m.metadataStore.GetBlobMetadataByStatusPaginated(ctx, status, exclusiveStartKey, limit)
	m.recordMetrics("GetBlobMetadataByStatusPaginated", start, err)
	return metadata, cursor, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error) {
	defer m.trackInFlight("GetBlobMetadataCountByStatus")()
	start := time.Now()
	count, err := m.metadataStore.GetBlobMetadataCountByStatus(ctx, status)
	m.recordMetrics("GetBlobMetadataCountByStatus", start, err)
	return count, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataByRequestedAtForward(
	ctx context.Context,
	after BlobFeedCursor,
	before BlobFeedCursor,
	limit int,
) ([]*v2.BlobMetadata, *BlobFeedCursor, error) {
	defer m.trackInFlight("GetBlobMetadataByRequestedAtForward")()
	start := time.Now()
	metadata, cursor, err := m.metadataStore.GetBlobMetadataByRequestedAtForward(ctx, after, before, limit)
	m.recordMetrics("GetBlobMetadataByRequestedAtForward", start, err)
	return metadata, cursor, err
}

func (m *InstrumentedMetadataStore) GetBlobMetadataByRequestedAtBackward(
	ctx context.Context,
	before BlobFeedCursor,
	after BlobFeedCursor,
	limit int,
) ([]*v2.BlobMetadata, *BlobFeedCursor, error) {
	defer m.trackInFlight("GetBlobMetadataByRequestedAtBackward")()
	start := time.Now()
	metadata, cursor, err := m.metadataStore.GetBlobMetadataByRequestedAtBackward(ctx, before, after, limit)
	m.recordMetrics("GetBlobMetadataByRequestedAtBackward", start, err)
	return metadata, cursor, err
}

func (m *InstrumentedMetadataStore) PutBlobCertificate(ctx context.Context, blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) error {
	defer m.trackInFlight("PutBlobCertificate")()
	start := time.Now()
	err := m.metadataStore.PutBlobCertificate(ctx, blobCert, fragmentInfo)
	m.recordMetrics("PutBlobCertificate", start, err)
	return err
}

func (m *InstrumentedMetadataStore) DeleteBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) error {
	defer m.trackInFlight("DeleteBlobCertificate")()
	start := time.Now()
	err := m.metadataStore.DeleteBlobCertificate(ctx, blobKey)
	m.recordMetrics("DeleteBlobCertificate", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	defer m.trackInFlight("GetBlobCertificate")()
	start := time.Now()
	cert, info, err := m.metadataStore.GetBlobCertificate(ctx, blobKey)
	m.recordMetrics("GetBlobCertificate", start, err)
	return cert, info, err
}

func (m *InstrumentedMetadataStore) GetBlobCertificates(ctx context.Context, blobKeys []corev2.BlobKey) ([]*corev2.BlobCertificate, []*encoding.FragmentInfo, error) {
	defer m.trackInFlight("GetBlobCertificates")()
	start := time.Now()
	certs, infos, err := m.metadataStore.GetBlobCertificates(ctx, blobKeys)
	m.recordMetrics("GetBlobCertificates", start, err)
	return certs, infos, err
}

func (m *InstrumentedMetadataStore) PutBatch(ctx context.Context, batch *corev2.Batch) error {
	defer m.trackInFlight("PutBatch")()
	start := time.Now()
	err := m.metadataStore.PutBatch(ctx, batch)
	m.recordMetrics("PutBatch", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Batch, error) {
	defer m.trackInFlight("GetBatch")()
	start := time.Now()
	batch, err := m.metadataStore.GetBatch(ctx, batchHeaderHash)
	m.recordMetrics("GetBatch", start, err)
	return batch, err
}

func (m *InstrumentedMetadataStore) PutBatchHeader(ctx context.Context, batchHeader *corev2.BatchHeader) error {
	defer m.trackInFlight("PutBatchHeader")()
	start := time.Now()
	err := m.metadataStore.PutBatchHeader(ctx, batchHeader)
	m.recordMetrics("PutBatchHeader", start, err)
	return err
}

func (m *InstrumentedMetadataStore) DeleteBatchHeader(ctx context.Context, batchHeaderHash [32]byte) error {
	defer m.trackInFlight("DeleteBatchHeader")()
	start := time.Now()
	err := m.metadataStore.DeleteBatchHeader(ctx, batchHeaderHash)
	m.recordMetrics("DeleteBatchHeader", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetBatchHeader(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, error) {
	defer m.trackInFlight("GetBatchHeader")()
	start := time.Now()
	header, err := m.metadataStore.GetBatchHeader(ctx, batchHeaderHash)
	m.recordMetrics("GetBatchHeader", start, err)
	return header, err
}

func (m *InstrumentedMetadataStore) PutDispersalRequest(ctx context.Context, req *corev2.DispersalRequest) error {
	defer m.trackInFlight("PutDispersalRequest")()
	start := time.Now()
	err := m.metadataStore.PutDispersalRequest(ctx, req)
	m.recordMetrics("PutDispersalRequest", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetDispersalRequest(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalRequest, error) {
	defer m.trackInFlight("GetDispersalRequest")()
	start := time.Now()
	req, err := m.metadataStore.GetDispersalRequest(ctx, batchHeaderHash, operatorID)
	m.recordMetrics("GetDispersalRequest", start, err)
	return req, err
}

func (m *InstrumentedMetadataStore) PutDispersalResponse(ctx context.Context, res *corev2.DispersalResponse) error {
	defer m.trackInFlight("PutDispersalResponse")()
	start := time.Now()
	err := m.metadataStore.PutDispersalResponse(ctx, res)
	m.recordMetrics("PutDispersalResponse", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetDispersalResponse(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalResponse, error) {
	defer m.trackInFlight("GetDispersalResponse")()
	start := time.Now()
	res, err := m.metadataStore.GetDispersalResponse(ctx, batchHeaderHash, operatorID)
	m.recordMetrics("GetDispersalResponse", start, err)
	return res, err
}

func (m *InstrumentedMetadataStore) GetDispersalResponses(ctx context.Context, batchHeaderHash [32]byte) ([]*corev2.DispersalResponse, error) {
	defer m.trackInFlight("GetDispersalResponses")()
	start := time.Now()
	responses, err := m.metadataStore.GetDispersalResponses(ctx, batchHeaderHash)
	m.recordMetrics("GetDispersalResponses", start, err)
	return responses, err
}

func (m *InstrumentedMetadataStore) GetDispersalsByRespondedAt(
	ctx context.Context,
	operatorId core.OperatorID,
	start uint64,
	end uint64,
	limit int,
	ascending bool,
) ([]*corev2.DispersalResponse, error) {
	defer m.trackInFlight("GetDispersalsByRespondedAt")()
	startTime := time.Now()
	responses, err := m.metadataStore.GetDispersalsByRespondedAt(ctx, operatorId, start, end, limit, ascending)
	m.recordMetrics("GetDispersalsByRespondedAt", startTime, err)
	return responses, err
}

func (m *InstrumentedMetadataStore) PutAttestation(ctx context.Context, attestation *corev2.Attestation) error {
	defer m.trackInFlight("PutAttestation")()
	start := time.Now()
	err := m.metadataStore.PutAttestation(ctx, attestation)
	m.recordMetrics("PutAttestation", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetAttestation(ctx context.Context, batchHeaderHash [32]byte) (*corev2.Attestation, error) {
	defer m.trackInFlight("GetAttestation")()
	start := time.Now()
	attestation, err := m.metadataStore.GetAttestation(ctx, batchHeaderHash)
	m.recordMetrics("GetAttestation", start, err)
	return attestation, err
}

func (m *InstrumentedMetadataStore) GetAttestationByAttestedAtForward(
	ctx context.Context,
	after uint64,
	before uint64,
	limit int,
) ([]*corev2.Attestation, error) {
	defer m.trackInFlight("GetAttestationByAttestedAtForward")()
	start := time.Now()
	attestations, err := m.metadataStore.GetAttestationByAttestedAtForward(ctx, after, before, limit)
	m.recordMetrics("GetAttestationByAttestedAtForward", start, err)
	return attestations, err
}

func (m *InstrumentedMetadataStore) GetAttestationByAttestedAtBackward(
	ctx context.Context,
	before uint64,
	after uint64,
	limit int,
) ([]*corev2.Attestation, error) {
	defer m.trackInFlight("GetAttestationByAttestedAtBackward")()
	start := time.Now()
	attestations, err := m.metadataStore.GetAttestationByAttestedAtBackward(ctx, before, after, limit)
	m.recordMetrics("GetAttestationByAttestedAtBackward", start, err)
	return attestations, err
}

func (m *InstrumentedMetadataStore) PutBlobInclusionInfo(ctx context.Context, inclusionInfo *corev2.BlobInclusionInfo) error {
	defer m.trackInFlight("PutBlobInclusionInfo")()
	start := time.Now()
	err := m.metadataStore.PutBlobInclusionInfo(ctx, inclusionInfo)
	m.recordMetrics("PutBlobInclusionInfo", start, err)
	return err
}

func (m *InstrumentedMetadataStore) PutBlobInclusionInfos(ctx context.Context, inclusionInfos []*corev2.BlobInclusionInfo) error {
	defer m.trackInFlight("PutBlobInclusionInfos")()
	start := time.Now()
	err := m.metadataStore.PutBlobInclusionInfos(ctx, inclusionInfos)
	m.recordMetrics("PutBlobInclusionInfos", start, err)
	return err
}

func (m *InstrumentedMetadataStore) GetBlobInclusionInfo(ctx context.Context, blobKey corev2.BlobKey, batchHeaderHash [32]byte) (*corev2.BlobInclusionInfo, error) {
	defer m.trackInFlight("GetBlobInclusionInfo")()
	start := time.Now()
	info, err := m.metadataStore.GetBlobInclusionInfo(ctx, blobKey, batchHeaderHash)
	m.recordMetrics("GetBlobInclusionInfo", start, err)
	return info, err
}

func (m *InstrumentedMetadataStore) GetBlobInclusionInfos(ctx context.Context, blobKey corev2.BlobKey) ([]*corev2.BlobInclusionInfo, error) {
	defer m.trackInFlight("GetBlobInclusionInfos")()
	start := time.Now()
	infos, err := m.metadataStore.GetBlobInclusionInfos(ctx, blobKey)
	m.recordMetrics("GetBlobInclusionInfos", start, err)
	return infos, err
}

func (m *InstrumentedMetadataStore) GetBlobAttestationInfo(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobAttestationInfo, error) {
	defer m.trackInFlight("GetBlobAttestationInfo")()
	start := time.Now()
	info, err := m.metadataStore.GetBlobAttestationInfo(ctx, blobKey)
	m.recordMetrics("GetBlobAttestationInfo", start, err)
	return info, err
}

func (m *InstrumentedMetadataStore) GetSignedBatch(ctx context.Context, batchHeaderHash [32]byte) (*corev2.BatchHeader, *corev2.Attestation, error) {
	defer m.trackInFlight("GetSignedBatch")()
	start := time.Now()
	header, attestation, err := m.metadataStore.GetSignedBatch(ctx, batchHeaderHash)
	m.recordMetrics("GetSignedBatch", start, err)
	return header, attestation, err
}
