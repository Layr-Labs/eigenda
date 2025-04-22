package clients

import (
	"crypto/tls"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// getGrpcDialOptions builds the gRPC dial options based on the useSecureGrpcFlag and maxMessageSize.
// When enableOpenTelemetry is true, OpenTelemetry instrumentation will be added to the dial options.
func getGrpcDialOptions(useSecureGrpcFlag bool, maxMessageSize uint, enableOpenTelemetry bool) []grpc.DialOption {
	options := []grpc.DialOption{}
	if useSecureGrpcFlag {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMessageSize))))

	// Add OpenTelemetry instrumentation if enabled
	if enableOpenTelemetry {
		options = append(options, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	return options
}
