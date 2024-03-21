package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"net"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
)

const localhost = "0.0.0.0"

// Server implements the Node proto APIs.
type Server struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	node   *node.Node
	config *node.Config
	logger logging.Logger

	ratelimiter common.RateLimiter

	mu *sync.Mutex
}

// NewServer creates a new Server instance with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewServer(config *node.Config, node *node.Node, logger logging.Logger, ratelimiter common.RateLimiter) *Server {

	return &Server{
		config:      config,
		logger:      logger,
		node:        node,
		ratelimiter: ratelimiter,
		mu:          &sync.Mutex{},
	}
}

func (s *Server) Start() {

	// TODO: In order to facilitate integration testing with multiple nodes, we need to be able to set the port.
	// TODO: Properly implement the health check.
	// go func() {
	// 	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(http.StatusOK)
	// 	})
	// }()

	// TODO: Add monitoring
	go func() {
		for {
			err := s.serveDispersal()
			s.logger.Error("dispersal server failed; restarting.", "err", err)
		}
	}()

	go func() {
		for {
			err := s.serveRetrieval()
			s.logger.Error("retrieval server failed; restarting.", "err", err)
		}
	}()

}

func (s *Server) serveDispersal() error {

	addr := fmt.Sprintf("%s:%s", localhost, s.config.InternalDispersalPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatalf("Could not start tcp listener: %v", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 1024) // 1 GiB
	gs := grpc.NewServer(opt)

	// Register reflection service on gRPC server
	// This makes "grpcurl -plaintext localhost:9000 list" command work
	reflection.Register(gs)

	pb.RegisterDispersalServer(gs, s)

	s.logger.Info("port", s.config.InternalDispersalPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return err
	}
	return nil

}

func (s *Server) serveRetrieval() error {
	addr := fmt.Sprintf("%s:%s", localhost, s.config.InternalRetrievalPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatalf("Could not start tcp listener: %v", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt)

	// Register reflection service on gRPC server
	// This makes "grpcurl -plaintext localhost:9000 list" command work
	reflection.Register(gs)

	pb.RegisterRetrievalServer(gs, s)

	s.logger.Info("port", s.config.InternalRetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return err
	}
	return nil

}

func (s *Server) handleStoreChunksRequest(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	// Get batch header hash
	batchHeader, err := GetBatchHeader(in)
	if err != nil {
		return nil, err
	}

	blobs, err := GetBlobMessages(in)
	if err != nil {
		return nil, err
	}

	sig, err := s.node.ProcessBatch(ctx, batchHeader, blobs, in.GetBlobs())
	if err != nil {
		return nil, err
	}

	sigData := sig.Serialize()

	return &pb.StoreChunksReply{Signature: sigData[:]}, nil
}

func (s *Server) validateStoreChunkRequest(in *pb.StoreChunksRequest) error {
	if in.GetBatchHeader() == nil {
		return api.NewInvalidArgError("missing batch_header in request")
	}
	if in.GetBatchHeader().GetBatchRoot() == nil {
		return api.NewInvalidArgError("missing batch_root in request")
	}
	if in.GetBatchHeader().GetReferenceBlockNumber() == 0 {
		return api.NewInvalidArgError("missing reference_block_number in request")
	}

	if len(in.GetBlobs()) == 0 {
		return api.NewInvalidArgError("missing blobs in request")
	}
	for _, blob := range in.Blobs {
		if blob.GetHeader() == nil {
			return api.NewInvalidArgError("missing blob header in request")
		}
		if len(blob.GetHeader().GetQuorumHeaders()) == 0 {
			return api.NewInvalidArgError("missing quorum headers in request")
		}
		if len(blob.GetHeader().GetQuorumHeaders()) != len(blob.GetBundles()) {
			return api.NewInvalidArgError("the number of quorums must be the same as the number of bundles")
		}
		for _, q := range blob.GetHeader().GetQuorumHeaders() {
			if q.GetAdversaryThreshold() >= q.GetConfirmationThreshold() {
				return api.NewInvalidArgError("adversary_threshold must be less than confirmation_threshold")
			}
			if q.GetConfirmationThreshold() > 100 {
				return api.NewInvalidArgError("confirmation threshold exceeds 100")
			}
			if q.AdversaryThreshold == 0 {
				return api.NewInvalidArgError("adversary threshold equals 0")
			}
		}
	}
	return nil
}

// StoreChunks is called by dispersers to store data.
func (s *Server) StoreChunks(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(sec float64) {
		s.node.Metrics.ObserveLatency("StoreChunks", "total", sec*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	// Validate the request.
	if err := s.validateStoreChunkRequest(in); err != nil {
		return nil, err
	}

	// Process the request.
	reply, err := s.handleStoreChunksRequest(ctx, in)

	// Record metrics.
	if err != nil {
		s.node.Metrics.RecordRPCRequest("StoreChunks", "failure")
		s.node.Logger.Error("StoreChunks failed", "err", err)
	} else {
		s.node.Metrics.RecordRPCRequest("StoreChunks", "success")
	}

	return reply, err
}

func (s *Server) RetrieveChunks(ctx context.Context, in *pb.RetrieveChunksRequest) (*pb.RetrieveChunksReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(sec float64) {
		s.node.Metrics.ObserveLatency("RetrieveChunks", "total", sec*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	if in.GetQuorumId() > core.MaxQuorumID {
		return nil, fmt.Errorf("invalid request: quorum ID must be in range [0, %d], but found %d", core.MaxQuorumID, in.GetQuorumId())
	}

	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], in.GetBatchHeaderHash())

	blobHeader, _, err := s.getBlobHeader(ctx, batchHeaderHash, int(in.BlobIndex), uint8(in.GetQuorumId()))
	if err != nil {
		return nil, err
	}

	retrieverID, err := common.GetClientAddress(ctx, s.config.ClientIPHeader, 1, false)
	if err != nil {
		return nil, err
	}

	quorumInfo := blobHeader.GetQuorumInfo(uint8(in.GetQuorumId()))
	if quorumInfo == nil {
		return nil, fmt.Errorf("invalid request: quorum ID %d not found in blob header", in.GetQuorumId())
	}
	encodedBlobSize := encoding.GetBlobSize(encoding.GetEncodedBlobLength(blobHeader.Length, quorumInfo.ConfirmationThreshold, quorumInfo.AdversaryThreshold))
	rate := quorumInfo.QuorumRate

	s.mu.Lock()
	allow, err := s.ratelimiter.AllowRequest(ctx, retrieverID, encodedBlobSize, rate)
	s.mu.Unlock()
	if err != nil {
		return nil, err
	}

	if !allow {
		return nil, errors.New("request rate limited")
	}

	chunks, ok := s.node.Store.GetChunks(ctx, batchHeaderHash, int(in.GetBlobIndex()), uint8(in.GetQuorumId()))
	if !ok {
		s.node.Metrics.RecordRPCRequest("RetrieveChunks", "failure")
		return nil, fmt.Errorf("could not find chunks for batchHeaderHash %v, blob index: %v, quorumID: %v", batchHeaderHash, in.GetBlobIndex(), in.GetQuorumId())
	}
	s.node.Metrics.RecordRPCRequest("RetrieveChunks", "success")
	return &pb.RetrieveChunksReply{Chunks: chunks}, nil
}

func (s *Server) GetBlobHeader(ctx context.Context, in *pb.GetBlobHeaderRequest) (*pb.GetBlobHeaderReply, error) {
	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], in.GetBatchHeaderHash())

	blobHeader, protoBlobHeader, err := s.getBlobHeader(ctx, batchHeaderHash, int(in.BlobIndex), uint8(in.GetQuorumId()))
	if err != nil {
		return nil, err
	}

	blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
	if err != nil {
		return nil, err
	}

	tree, err := s.rebuildMerkleTree(batchHeaderHash, uint8(in.GetQuorumId()))
	if err != nil {
		return nil, err
	}

	proof, err := tree.GenerateProof(blobHeaderHash[:], 0)
	if err != nil {
		return nil, err
	}

	return &pb.GetBlobHeaderReply{
		BlobHeader: protoBlobHeader,
		Proof: &pb.MerkleProof{
			Hashes: proof.Hashes,
			Index:  uint32(proof.Index),
		},
	}, nil
}

func (s *Server) getBlobHeader(ctx context.Context, batchHeaderHash [32]byte, blobIndex int, quorumId uint8) (*core.BlobHeader, *pb.BlobHeader, error) {

	blobHeaderBytes, err := s.node.Store.GetBlobHeader(ctx, batchHeaderHash, blobIndex)
	if err != nil {
		return nil, nil, errors.New("failed to get the blob header from Store")
	}

	var protoBlobHeader pb.BlobHeader
	err = proto.Unmarshal(blobHeaderBytes, &protoBlobHeader)
	if err != nil {
		return nil, nil, err
	}

	blobHeader, err := GetBlobHeaderFromProto(&protoBlobHeader)
	if err != nil {
		return nil, nil, err
	}

	return blobHeader, &protoBlobHeader, nil

}
