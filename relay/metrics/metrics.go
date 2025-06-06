package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const namespace = "eigenda_relay"

type RelayMetrics struct {
	logger           logging.Logger
	grpcServerOption grpc.ServerOption
	server           *http.Server

	// Cache metrics
	MetadataCacheMetrics *cache.CacheAccessorMetrics
	ChunkCacheMetrics    *cache.CacheAccessorMetrics
	BlobCacheMetrics     *cache.CacheAccessorMetrics

	// GetChunks metrics
	getChunksLatency               *prometheus.SummaryVec
	getChunksAuthenticationLatency *prometheus.SummaryVec
	getChunksMetadataLatency       *prometheus.SummaryVec
	getChunksDataLatency           *prometheus.SummaryVec
	getChunksAuthFailures          *prometheus.CounterVec
	getChunksRateLimited           *prometheus.CounterVec
	getChunksKeyCount              *prometheus.GaugeVec
	getChunksBandwidth             *prometheus.CounterVec
	getChunksRequestedBandwidth    *prometheus.CounterVec

	// GetBlob metrics
	getBlobLatency            *prometheus.SummaryVec
	getBlobMetadataLatency    *prometheus.SummaryVec
	getBlobDataLatency        *prometheus.SummaryVec
	getBlobRateLimited        *prometheus.CounterVec
	getBlobBandwidth          *prometheus.CounterVec
	getBlobRequestedBandwidth *prometheus.CounterVec
}

// NewRelayMetrics creates a new RelayMetrics instance, which encapsulates all metrics related to the relay.
func NewRelayMetrics(registry *prometheus.Registry, logger logging.Logger, port int) *RelayMetrics {
	if registry == nil {
		registry = prometheus.NewRegistry()
	}
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", port)
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	metadataCacheMetrics := cache.NewCacheAccessorMetrics(registry, "metadata")
	chunkCacheMetrics := cache.NewCacheAccessorMetrics(registry, "chunk")
	blobCacheMetrics := cache.NewCacheAccessorMetrics(registry, "blob")

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	getChunksLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_latency_ms",
			Help:       "Latency of the GetChunks RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksAuthenticationLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_authentication_latency_ms",
			Help:       "Latency of the GetChunks RPC client authentication",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksMetadataLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_metadata_latency_ms",
			Help:       "Latency of the GetChunks RPC metadata retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksDataLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_data_latency_ms",
			Help:       "Latency of the GetChunks RPC data retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksAuthFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_auth_failure_count",
			Help:      "Number of GetChunks RPC authentication failures",
		},
		[]string{},
	)

	getChunksRateLimited := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_rate_limited_count",
			Help:      "Number of GetChunks RPC rate limited",
		},
		[]string{"reason"},
	)

	getChunksKeyCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "get_chunks_key_count",
			Help:      "Number of keys in a GetChunks request.",
		},
		[]string{},
	)

	getChunksBandwidth := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_bandwidth_bytes",
			Help:      "Running total bandwidth used in GetChunks requests.",
		},
		[]string{},
	)

	getChunksRequestedBandwidth := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_requested_bandwidth_bytes",
			Help:      "Running total requested bandwidth in GetChunks requests (prior to throttling).",
		},
		[]string{},
	)

	getBlobLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_latency_ms",
			Help:       "Latency of the GetBlob RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobMetadataLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_metadata_latency_ms",
			Help:       "Latency of the GetBlob RPC metadata retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobDataLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_data_latency_ms",
			Help:       "Latency of the GetBlob RPC data retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobRateLimited := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_blob_rate_limited_count",
			Help:      "Number of GetBlob RPC rate limited",
		},
		[]string{"reason"},
	)

	getBlobBandwidth := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_blob_bandwidth_bytes",
			Help:      "Running total bandwidth used in GetBlob requests.",
		},
		[]string{},
	)

	getBlobRequestedBandwidth := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_blob_requested_bandwidth_bytes",
			Help:      "Running total requested bandwidth in GetBlob requests (prior to throttling).",
		},
		[]string{},
	)

	return &RelayMetrics{
		logger:                         logger,
		grpcServerOption:               grpcServerOption,
		server:                         server,
		MetadataCacheMetrics:           metadataCacheMetrics,
		ChunkCacheMetrics:              chunkCacheMetrics,
		BlobCacheMetrics:               blobCacheMetrics,
		getChunksLatency:               getChunksLatency,
		getChunksAuthenticationLatency: getChunksAuthenticationLatency,
		getChunksMetadataLatency:       getChunksMetadataLatency,
		getChunksDataLatency:           getChunksDataLatency,
		getChunksAuthFailures:          getChunksAuthFailures,
		getChunksRateLimited:           getChunksRateLimited,
		getChunksKeyCount:              getChunksKeyCount,
		getChunksBandwidth:             getChunksBandwidth,
		getChunksRequestedBandwidth:    getChunksRequestedBandwidth,
		getBlobLatency:                 getBlobLatency,
		getBlobMetadataLatency:         getBlobMetadataLatency,
		getBlobDataLatency:             getBlobDataLatency,
		getBlobRateLimited:             getBlobRateLimited,
		getBlobBandwidth:               getBlobBandwidth,
		getBlobRequestedBandwidth:      getBlobRequestedBandwidth,
	}
}

// Start starts the metrics server.
func (m *RelayMetrics) Start() {
	go func() {
		err := m.server.ListenAndServe()
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			m.logger.Errorf("metrics server error: %v", err)
		}
	}()
}

// Stop stops the metrics server.
func (m *RelayMetrics) Stop() error {
	return m.server.Close()
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *RelayMetrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}

func (m *RelayMetrics) ReportChunkLatency(duration time.Duration) {
	m.getChunksLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportChunkAuthenticationLatency(duration time.Duration) {
	m.getChunksAuthenticationLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportChunkMetadataLatency(duration time.Duration) {
	m.getChunksMetadataLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportChunkDataLatency(duration time.Duration) {
	m.getChunksDataLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportChunkAuthFailure() {
	m.getChunksAuthFailures.WithLabelValues().Inc()
}

func (m *RelayMetrics) ReportChunkRateLimited(reason string) {
	m.getChunksRateLimited.WithLabelValues(reason).Inc()
}

func (m *RelayMetrics) ReportChunkKeyCount(count int) {
	m.getChunksKeyCount.WithLabelValues().Set(float64(count))
}

func (m *RelayMetrics) ReportGetChunksBandwidthUsage(size uint32) {
	m.getChunksBandwidth.WithLabelValues().Add(float64(size))
}

func (m *RelayMetrics) ReportGetChunksRequestedBandwidthUsage(size uint32) {
	m.getChunksRequestedBandwidth.WithLabelValues().Add(float64(size))
}

func (m *RelayMetrics) ReportBlobLatency(duration time.Duration) {
	m.getBlobLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportBlobMetadataLatency(duration time.Duration) {
	m.getBlobMetadataLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportBlobDataLatency(duration time.Duration) {
	m.getBlobDataLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *RelayMetrics) ReportBlobRateLimited(reason string) {
	m.getBlobRateLimited.WithLabelValues(reason).Inc()
}

func (m *RelayMetrics) ReportBlobBandwidthUsage(size int) {
	m.getBlobBandwidth.WithLabelValues().Add(float64(size))
}

func (m *RelayMetrics) ReportBlobRequestedBandwidthUsage(size int) {
	m.getBlobRequestedBandwidth.WithLabelValues().Add(float64(size))
}
