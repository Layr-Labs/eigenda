package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
)

const namespace := "eigenda_relay"

type RelayMetrics struct {
	grpcServerOption grpc.ServerOption

	// Cache metrics
	MetadataCacheMetrics *cache.CacheAccessorMetrics
	ChunkCacheMetrics    *cache.CacheAccessorMetrics
	BlobCacheMetrics     *cache.CacheAccessorMetrics

	// GetChunks metrics
	GetChunksLatency               *prometheus.SummaryVec
	GetChunksAuthenticationLatency *prometheus.SummaryVec
	GetChunksMetadataLatency       *prometheus.SummaryVec
	GetChunksDataLatency           *prometheus.SummaryVec
	GetChunksAuthFailures          *prometheus.CounterVec
	GetChunksRateLimited           *prometheus.CounterVec
	GetChunksKeyCount              *prometheus.GaugeVec
	GetChunksDataSize              *prometheus.GaugeVec

	// GetBlob metrics
	GetBlobLatency         *prometheus.SummaryVec
	GetBlobMetadataLatency *prometheus.SummaryVec
	GetBlobDataLatency     *prometheus.SummaryVec
	GetBlobRateLimited     *prometheus.CounterVec
	GetBlobDataSize        *prometheus.GaugeVec
}

// NewRelayMetrics creates a new RelayMetrics instance, which encapsulates all metrics related to the relay.
func NewRelayMetrics(logger logging.Logger, port int) (*RelayMetrics, error) {

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", port)
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{},
	))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	grpcMetrics := grpcprom.NewServerMetrics()
	reg.MustRegister(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	metadataCacheMetrics, err := cache.NewCacheAccessorMetrics(server, "metadata")
	if err != nil {
		return nil, err
	}

	chunkCacheMetrics, err := cache.NewCacheAccessorMetrics(server, "chunk")
	if err != nil {
		return nil, err
	}

	blobCacheMetrics, err := cache.NewCacheAccessorMetrics(server, "blob")
	if err != nil {
		return nil, err
	}

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	getChunksLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_latency_ms",
			Help:       "Latency of the GetChunks RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksAuthenticationLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_authentication_latency_ms",
			Help:       "Latency of the GetChunks RPC client authentication",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksMetadataLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_metadata_latency_ms",
			Help:       "Latency of the GetChunks RPC metadata retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksDataLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_data_latency_ms",
			Help:       "Latency of the GetChunks RPC data retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getChunksAuthFailures := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_auth_failure_count",
			Help:      "Number of GetChunks RPC authentication failures",
		},
		[]string{},
	)

	getChunksRateLimited := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_chunks_rate_limited_count",
			Help:      "Number of GetChunks RPC rate limited",
		},
		[]string{"reason"},
	)

	getChunksKeyCount := promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "get_chunks_key_count",
			Help:      "Number of keys in a GetChunks request.",
		},
		[]string{},
	)

	getChunksDataSize := promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "get_chunks_data_size_bytes",
			Help:      "Data size in a GetChunks request.",
		},
		[]string{},
	)

	getBlobLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_latency_ms",
			Help:       "Latency of the GetBlob RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobMetadataLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_metadata_latency_ms",
			Help:       "Latency of the GetBlob RPC metadata retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobDataLatency := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_data_latency_ms",
			Help:       "Latency of the GetBlob RPC data retrieval",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobRateLimited := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "get_blob_rate_limited_count",
			Help:      "Number of GetBlob RPC rate limited",
		},
		[]string{"reason"},
	)

	getBlobDataSize := promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "get_blob_data_size_bytes",
			Help:      "Data size of requested blobs.",
		},
		[]string{},
	)

	return &RelayMetrics{
		grpcServerOption: grpcServerOption,
		MetadataCacheMetrics: metadataCacheMetrics,
		ChunkCacheMetrics:    chunkCacheMetrics,
		BlobCacheMetrics:     blobCacheMetrics,
		GetChunksLatency:               getChunksLatency,
		GetChunksAuthenticationLatency: getChunksAuthenticationLatency,
		GetChunksMetadataLatency:       getChunksMetadataLatency,
		GetChunksDataLatency:           getChunksDataLatency,
		GetChunksAuthFailures:          getChunksAuthFailures,
		GetChunksRateLimited:           getChunksRateLimited,
		GetChunksKeyCount:              getChunksKeyCount,
		GetChunksDataSize:              getChunksDataSize,
		GetBlobLatency:                 getBlobLatency,
		GetBlobMetadataLatency:         getBlobMetadataLatency,
		GetBlobDataLatency:             getBlobDataLatency,
		GetBlobRateLimited:             getBlobRateLimited,
		GetBlobDataSize:                getBlobDataSize,
	}, nil
}

// Start starts the metrics server.
func (m *RelayMetrics) Start() error {
	return m.metricsServer.Start()
}

// Stop stops the metrics server.
func (m *RelayMetrics) Stop() error {
	return m.metricsServer.Stop()
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *RelayMetrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}

// WriteMetricsDocumentation writes the metrics for the churner to a markdown file.
func (m *RelayMetrics) WriteMetricsDocumentation() error {
	return m.metricsServer.WriteMetricsDocumentation("relay/mdoc/relay-metrics.md")
}
