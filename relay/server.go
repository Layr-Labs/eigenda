package relay

import (
	"context"
	"fmt"
	v2pb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	config      *Config
	blobStore   *blobstore.BlobStore
	chunkReader *chunkstore.ChunkReader

	// metadataServer encapsulates logic for fetching metadata for blobs.
	metadataServer *metadataServer

	// blobServer encapsulates logic for fetching blobs.
	blobServer *blobServer

	// chunkServer encapsulates logic for fetching chunks.
	chunkServer *chunkServer
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader *chunkstore.ChunkReader) (*Server, error) {

	ms, err := newMetadataServer(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataWorkPoolSize,
		config.Shards)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata server: %w", err)
	}

	bs, err := newBlobServer(
		ctx,
		logger,
		blobStore,
		config.BlobCacheSize,
		config.BlobWorkPoolSize)
	if err != nil {
		return nil, fmt.Errorf("error creating blob server: %w", err)
	}

	return &Server{
		config:         config,
		metadataServer: ms,
		blobServer:     bs,
		blobStore:      blobStore,
		chunkReader:    chunkReader,
	}, nil
}

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {

	// TODO:
	//  - global throttle requests / sec
	//  - per-connection throttle requests / sec
	//  - timeouts

	keys := []v2.BlobKey{v2.BlobKey(request.BlobKey)}
	mMap, err := s.metadataServer.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %w", err)
	}
	metadata := (*mMap)[v2.BlobKey(request.BlobKey)]
	if metadata == nil {
		return nil, fmt.Errorf("blob not found")
	}

	// TODO
	//  - global bytes/sec throttle
	//  - per-connection bytes/sec throttle

	key := v2.BlobKey(request.BlobKey)
	data, err := s.blobServer.GetBlob(key)
	if err != nil {
		return nil, fmt.Errorf("error fetching blob: %w", err)
	}

	reply := &pb.GetBlobReply{
		Blob: data,
	}

	return reply, nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// TODO:
	//  - authentication
	//  - global throttle requests / sec
	//  - per-connection throttle requests / sec
	//  - timeouts

	keys := make([]v2.BlobKey, 0) // TODO

	for _, chunkRequest := range request.ChunkRequests {
		var key v2.BlobKey
		if chunkRequest.GetByIndex() != nil {
			key = v2.BlobKey(chunkRequest.GetByIndex().GetBlobKey())
		} else {
			key = v2.BlobKey(chunkRequest.GetByRange().GetBlobKey())
		}
		keys = append(keys, key)
	}

	mMap, err := s.metadataServer.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %w", err)
	}

	frames, err := s.chunkServer.GetFrames(ctx, mMap)
	if err != nil {
		return nil, fmt.Errorf("error fetching frames: %w", err)
	}

	protoChunks := make([]*pb.Chunks, 0, len(*frames))

	// TODO encapsulate
	// return data in the order that it was requested
	for _, chunkRequest := range request.ChunkRequests {
		if chunkRequest.GetByIndex() != nil {
			key := v2.BlobKey(chunkRequest.GetByIndex().GetBlobKey())
			blobFrames := (*frames)[key]
			chunks := &pb.Chunks{
				Data: make([]*v2pb.Frame, 0, len(chunkRequest.GetByIndex().ChunkIndices)),
			}
			protoChunks = append(protoChunks, chunks)

			for index := range chunkRequest.GetByIndex().ChunkIndices {

				if index >= len(blobFrames) {
					return nil, fmt.Errorf(
						"chunk index %d out of range for key %s, chunk count %d",
						index, key, len(blobFrames))
				}
				chunks.Data = append(chunks.Data, blobFrames[index].ToProtobuf())
			}

		} else {
			key := v2.BlobKey(chunkRequest.GetByRange().GetBlobKey())

			blobFrames := (*frames)[key]
			chunks := &pb.Chunks{
				Data: make([]*v2pb.Frame, 0, len(chunkRequest.GetByIndex().ChunkIndices)),
			}
			protoChunks = append(protoChunks, chunks)

			if chunkRequest.GetByRange().EndIndex >= uint32(len(blobFrames)) {
				return nil, fmt.Errorf(
					"chunk range %d-%d is invald for key %s, chunk count %d",
					chunkRequest.GetByRange().StartIndex, chunkRequest.GetByRange().EndIndex, key, len(blobFrames))
			}
			for index := chunkRequest.GetByRange().StartIndex; index <= chunkRequest.GetByRange().EndIndex; index++ {
				chunks.Data = append(chunks.Data, blobFrames[index].ToProtobuf())
			}
		}

	}

	return &pb.GetChunksReply{
		Data: protoChunks,
	}, nil
}
