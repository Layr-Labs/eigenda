package relay

import (
	"context"
	"errors"
	"fmt"
	v2pb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	// the logger for the server
	logger logging.Logger

	// grpcPort is the port that the relay server listens on.
	grpcPort int

	// maxProtoSize is the maximum size of a gRPC message that the server will accept.
	maxProtoSize int

	// metadataServer encapsulates logic for fetching metadata for blobs.
	metadataServer *metadataManager

	// blobServer encapsulates logic for fetching blobs.
	blobServer *blobManager

	// chunkServer encapsulates logic for fetching chunks.
	chunkServer *chunkManager

	// grpcServer is the gRPC server.
	grpcServer *grpc.Server
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader chunkstore.ChunkReader) (*Server, error) {

	ms, err := newMetadataManager(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataWorkPoolSize,
		config.Shards)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata server: %w", err)
	}

	bs, err := newBlobManager(
		ctx,
		logger,
		blobStore,
		config.BlobCacheSize,
		config.BlobWorkPoolSize)
	if err != nil {
		return nil, fmt.Errorf("error creating blob server: %w", err)
	}

	cs, err := newChunkManager(
		ctx,
		logger,
		chunkReader,
		config.ChunkCacheSize,
		config.ChunkWorkPoolSize)
	if err != nil {
		return nil, fmt.Errorf("error creating chunk server: %w", err)
	}

	return &Server{
		logger:         logger,
		grpcPort:       config.GRPCPort,
		maxProtoSize:   config.MaxGRPCMessageSize,
		metadataServer: ms,
		blobServer:     bs,
		chunkServer:    cs,
	}, nil
}

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {

	// Future work:
	//  - global throttle
	//  - per-connection throttle
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

	key := v2.BlobKey(request.BlobKey)
	data, err := s.blobServer.GetBlob(key)
	if err != nil {
		return nil, fmt.Errorf("error fetching blob %s: %w", key.Hex(), err)
	}

	reply := &pb.GetBlobReply{
		Blob: data,
	}

	return reply, nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// Future work:
	//  - authentication
	//  - global throttle
	//  - per-connection throttle
	//  - timeouts

	keys := make([]v2.BlobKey, 0, len(request.GetBlobKeys()))

	for _, keyBytes := range request.GetBlobKeys() {
		keys = append(keys, v2.BlobKey(keyBytes))
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

	// return data in the order that it was requested
	for _, keyBytes := range request.GetBlobKeys() {
		key := v2.BlobKey(keyBytes)
		if request.GetByIndex() != nil {
			blobFrames := (*frames)[key]
			chunks := &pb.Chunks{
				Data: make([]*v2pb.Frame, 0, len(request.GetByIndex().ChunkIndices)),
			}
			protoChunks = append(protoChunks, chunks)

			for index := range request.GetByIndex().ChunkIndices {

				if index >= len(blobFrames) {
					return nil, fmt.Errorf(
						"chunk index %d out of range for key %s, chunk count %d",
						index, key.Hex(), len(blobFrames))
				}
				chunks.Data = append(chunks.Data, blobFrames[index].ToProtobuf())
			}

		} else {
			startIndex := request.GetByRange().StartIndex
			endIndex := request.GetByRange().EndIndex

			blobFrames := (*frames)[key]

			if startIndex > endIndex {
				return nil, fmt.Errorf(
					"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
					startIndex, endIndex, key.Hex())
			}
			if endIndex > uint32(len((*frames)[key])) {
				return nil, fmt.Errorf(
					"chunk range %d-%d is invald for key %s, chunk count %d",
					request.GetByRange().StartIndex, request.GetByRange().EndIndex, key, len(blobFrames))
			}

			chunks := &pb.Chunks{
				Data: make([]*v2pb.Frame, 0, endIndex-startIndex),
			}
			protoChunks = append(protoChunks, chunks)

			for index := startIndex; index < endIndex; index++ {
				chunks.Data = append(chunks.Data, blobFrames[index].ToProtobuf())
			}
		}
	}

	return &pb.GetChunksReply{
		Data: protoChunks,
	}, nil
}

// Start starts the server listening for requests. This method will block until the server is stopped.
func (s *Server) Start() error {

	// Serve grpc requests
	addr := fmt.Sprintf("0.0.0.0:%d", s.grpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener on %s: %w", addr, err)
	}

	opt := grpc.MaxRecvMsgSize(s.maxProtoSize)

	s.grpcServer = grpc.NewServer(opt)
	reflection.Register(s.grpcServer)
	pb.RegisterRelayServer(s.grpcServer, s)

	// Register Server for Health Checks
	name := pb.Relay_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, s.grpcServer)

	s.logger.Info("GRPC Listening", "port", s.grpcPort, "address", listener.Addr().String())

	if err = s.grpcServer.Serve(listener); err != nil {
		return errors.New("could not start GRPC server")
	}

	return nil
}

// Stop stops the server.
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
}
