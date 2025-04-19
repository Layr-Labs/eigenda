package node

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import for channelz service
	"google.golang.org/grpc/channelz/service"
)

// RegisterChannelzService registers the channelz service with a gRPC server
func RegisterChannelzService(s *grpc.Server) {
	service.RegisterChannelzServiceToServer(s)
}

// GetDialOptions returns the gRPC dial options configured with channelz
func GetDialOptions(metrics *Metrics, logger logging.Logger, useSecure bool, maxMsgSize uint, connectionID string) []grpc.DialOption {
	// Create a unique connection ID if none provided
	if connectionID == "" {
		connectionID = fmt.Sprintf("conn-%d", time.Now().UnixNano())
	}

	// Base options
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMsgSize))),
		// Enable channelz data collection with nil parent
		grpc.WithChannelzParentID(nil),
	}

	// Add secure or insecure transport credentials
	if useSecure {
		// Not implemented - would add TLS credentials
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return opts
}
