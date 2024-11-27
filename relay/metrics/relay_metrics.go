package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

type RelayMetrics struct {
	metricsServer    metrics.Metrics
	grpcServerOption grpc.ServerOption

	GetChunksLatency               metrics.LatencyMetric
	GetChunksAuthenticationLatency metrics.LatencyMetric
	GetChunksMetadataLatency       metrics.LatencyMetric
	GetChunksDataLatency           metrics.LatencyMetric
	GetChunksAuthFailures          metrics.CountMetric
	GetChunksRateLimited           metrics.CountMetric
	GetChunksKeyCountHistogram     metrics.HistogramMetric
	GetChunksDataSizeHistogram     metrics.HistogramMetric
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
		limiter.RateLimitLabel{})
	if err != nil {
		return nil, err
	}

	getChunksKeyCountHistogram, err := server.NewHistogramMetric(
		"get_chunks_key",
		"count",
		"Number of keys in a GetChunks request",
		1.1,
		nil)
	if err != nil {
		return nil, err
	}

	getChunksDataSizeHistogram, err := server.NewHistogramMetric(
		"get_chunks_data",
		"bytes",
		"Size of data in a GetChunks request, in bytes",
		1.1,
		nil)
	if err != nil {
		return nil, err
	}

	return &RelayMetrics{
		metricsServer:                  server,
		grpcServerOption:               grpcServerOption,
		GetChunksLatency:               getChunksLatencyMetric,
		GetChunksAuthenticationLatency: getChunksAuthenticationLatencyMetric,
		GetChunksMetadataLatency:       getChunksMetadataLatencyMetric,
		GetChunksDataLatency:           getChunksDataLatencyMetric,
		GetChunksAuthFailures:          getChunksAuthFailures,
		GetChunksRateLimited:           getChunksRateLimited,
		GetChunksKeyCountHistogram:     getChunksKeyCountHistogram,
		GetChunksDataSizeHistogram:     getChunksDataSizeHistogram,
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
