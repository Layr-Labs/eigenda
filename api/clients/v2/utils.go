package clients

import (
	"crypto/tls"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GetGrpcDialOptions builds the gRPC dial options based on the useSecureGrpcFlag and maxMessageSize.
func GetGrpcDialOptions(useSecureGrpcFlag bool, maxMessageSize uint) []grpc.DialOption {
	options := []grpc.DialOption{
		// Automatic OpenTelemetry tracing for all gRPC calls
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	if useSecureGrpcFlag {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMessageSize))))

	return options
}
