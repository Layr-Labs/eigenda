package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

type RelayMetrics struct {
	metricsServer    metrics.Metrics
	grpcServerOption grpc.ServerOption
}

func NewRelayMetrics(
	logger logging.Logger,
	config *metrics.Config) (*RelayMetrics, error) {

	server := metrics.NewMetrics(logger, config)

	grpcMetrics := grpcprom.NewServerMetrics()
	server.RegisterExternalMetrics(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	return &RelayMetrics{
		metricsServer:    server,
		grpcServerOption: grpcServerOption,
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

// WriteMetricsDocumentation writes documentation for all currently registered metrics to a file.
func (m *RelayMetrics) WriteMetricsDocumentation() error {
	return m.metricsServer.WriteMetricsDocumentation("relay/metrics/relay-metrics.md")
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *RelayMetrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}
