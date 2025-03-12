package clients

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// getGrpcDialOptions builds the gRPC dial options based on the useSecureGrpcFlag and maxMessageSize.
func getGrpcDialOptions(useSecureGrpcFlag bool, maxMessageSize uint) []grpc.DialOption {
	options := []grpc.DialOption{}
	if useSecureGrpcFlag {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMessageSize))))

	return options
}
