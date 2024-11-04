package relay

import (
	"context"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	config        *Config
	metadataStore *blobstore.BlobMetadataStore
	blobStore     *blobstore.BlobStore
	chunkReader   *chunkstore.ChunkReader
}

func NewServer(
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader *chunkstore.ChunkReader) *Server {

	return &Server{
		config:        config,
		metadataStore: metadataStore,
		blobStore:     blobStore,
		chunkReader:   chunkReader,
	}
}

// GetBlobs retrieves blobs stored by the relay.
func (s *Server) GetBlobs(context.Context, *pb.GetBlobsRequest) (*pb.GetBlobsReply, error) {
	return nil, nil // TODO
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	return nil, nil // TODO
}
