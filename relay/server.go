package relay

import (
	"context"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
)

var _ pb.RelayServer = &Server{}

func NewServer(config *Config) *Server {
	return &Server{}
}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer
}

// GetBlobs retrieves blobs stored by the relay.
func (s *Server) GetBlobs(context.Context, *pb.GetBlobsRequest) (*pb.GetBlobsReply, error) {
	return nil, nil // TODO
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	return nil, nil // TODO
}
