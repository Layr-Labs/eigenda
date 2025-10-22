package clients

import (
	"crypto/tls"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GetGrpcDialOptions builds the gRPC dial options based on the useSecureGrpcFlag, maxMessageSize, and tracingEnabled.
func GetGrpcDialOptions(useSecureGrpcFlag bool, maxMessageSize uint, tracingEnabled bool) []grpc.DialOption {
	options := []grpc.DialOption{}

	// Only add OpenTelemetry tracing interceptor if tracing is enabled
	if tracingEnabled {
		options = append(options, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}
	if useSecureGrpcFlag {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMessageSize))))

	return options
}
