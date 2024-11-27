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

// NewRelayMetrics creates a new RelayMetrics instance, which encapsulates all metrics related to the relay.
func NewRelayMetrics(logger logging.Logger, port int) (*RelayMetrics, error) {

	server := metrics.NewMetrics(logger, "relay", port)

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

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *RelayMetrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}
