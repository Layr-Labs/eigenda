package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

type RelayMetrics struct {
	metricsServer    metrics.Metrics
	grpcServerOption grpc.ServerOption

	// Cache metrics
	MetadataCacheMetrics *cache.CacheAccessorMetrics
	ChunkCacheMetrics    *cache.CacheAccessorMetrics
	BlobCacheMetrics     *cache.CacheAccessorMetrics

	// GetChunks metrics
	GetChunksLatency               metrics.LatencyMetric
	GetChunksAuthenticationLatency metrics.LatencyMetric
	GetChunksMetadataLatency       metrics.LatencyMetric
	GetChunksDataLatency           metrics.LatencyMetric
	GetChunksAuthFailures          metrics.CountMetric
	GetChunksRateLimited           metrics.CountMetric
	GetChunksKeyCount              metrics.GaugeMetric
	GetChunksDataSize              metrics.GaugeMetric

	// GetBlob metrics
	GetBlobLatency         metrics.LatencyMetric
	GetBlobMetadataLatency metrics.LatencyMetric
	GetBlobDataLatency     metrics.LatencyMetric
	GetBlobRateLimited     metrics.CountMetric
	GetBlobDataSize        metrics.GaugeMetric
}

type RateLimitLabel struct {
	Reason string
}

// NewRelayMetrics creates a new RelayMetrics instance, which encapsulates all metrics related to the relay.
func NewRelayMetrics(logger logging.Logger, port int) (*RelayMetrics, error) {

	server := metrics.NewMetrics(logger, "relay", port)

	grpcMetrics := grpcprom.NewServerMetrics()
	server.RegisterExternalMetrics(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	standardQuantiles := []*metrics.Quantile{
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99),
	}

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

	getChunksLatencyMetric, err := server.NewLatencyMetric(
		"get_chunks_latency",
		"Latency of the GetChunks RPC",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getChunksAuthenticationLatencyMetric, err := server.NewLatencyMetric(
		"get_chunks_authentication_latency",
		"Latency of the GetChunks RPC client authentication",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getChunksMetadataLatencyMetric, err := server.NewLatencyMetric(
		"get_chunks_metadata_latency",
		"Latency of the GetChunks RPC metadata retrieval",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getChunksDataLatencyMetric, err := server.NewLatencyMetric(
		"get_chunks_data_latency",
		"Latency of the GetChunks RPC data retrieval",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getChunksAuthFailures, err := server.NewCountMetric(
		"get_chunks_auth_failure",
		"Number of GetChunks RPC authentication failures",
		nil)
	if err != nil {
		return nil, err
	}

	getChunksRateLimited, err := server.NewCountMetric(
		"get_chunks_rate_limited",
		"Number of GetChunks RPC rate limited",
		RateLimitLabel{})
	if err != nil {
		return nil, err
	}

	getChunksKeyCount, err := server.NewGaugeMetric(
		"get_chunks_key",
		"count",
		"Number of keys in a GetChunks request.",
		nil)
	if err != nil {
		return nil, err
	}

	getChunksDataSize, err := server.NewGaugeMetric(
		"get_chunks_data_size",
		"bytes",
		"Data size in a GetChunks request.",
		nil)
	if err != nil {
		return nil, err
	}

	getBlobLatencyMetric, err := server.NewLatencyMetric(
		"get_blob_latency",
		"Latency of the GetBlob RPC",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getBlobMetadataLatencyMetric, err := server.NewLatencyMetric(
		"get_blob_metadata_latency",
		"Latency of the GetBlob RPC metadata retrieval",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getBlobDataLatencyMetric, err := server.NewLatencyMetric(
		"get_blob_data_latency",
		"Latency of the GetBlob RPC data retrieval",
		nil,
		standardQuantiles...)
	if err != nil {
		return nil, err
	}

	getBlobRateLimited, err := server.NewCountMetric(
		"get_blob_rate_limited",
		"Number of GetBlob RPC rate limited",
		RateLimitLabel{})
	if err != nil {
		return nil, err
	}

	getBlobDataSize, err := server.NewGaugeMetric(
		"get_blob_data_size",
		"bytes",
		"Data size of requested blobs.",
		nil)
	if err != nil {
		return nil, err
	}

	return &RelayMetrics{
		metricsServer:                  server,
		MetadataCacheMetrics:           metadataCacheMetrics,
		ChunkCacheMetrics:              chunkCacheMetrics,
		BlobCacheMetrics:               blobCacheMetrics,
		grpcServerOption:               grpcServerOption,
		GetChunksLatency:               getChunksLatencyMetric,
		GetChunksAuthenticationLatency: getChunksAuthenticationLatencyMetric,
		GetChunksMetadataLatency:       getChunksMetadataLatencyMetric,
		GetChunksDataLatency:           getChunksDataLatencyMetric,
		GetChunksAuthFailures:          getChunksAuthFailures,
		GetChunksRateLimited:           getChunksRateLimited,
		GetChunksKeyCount:              getChunksKeyCount,
		GetChunksDataSize:              getChunksDataSize,
		GetBlobLatency:                 getBlobLatencyMetric,
		GetBlobMetadataLatency:         getBlobMetadataLatencyMetric,
		GetBlobDataLatency:             getBlobDataLatencyMetric,
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
