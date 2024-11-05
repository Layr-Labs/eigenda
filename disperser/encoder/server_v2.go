package encoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/common/healthcheck"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type EncoderServerV2 struct {
	pb.UnimplementedEncoderServer

	config      ServerConfig
	blobStore   *blobstore.BlobStore
	chunkWriter chunkstore.ChunkWriter
	logger      logging.Logger
	prover      encoding.Prover
	metrics     *Metrics
	close       func()

	runningRequests chan struct{}
	requestPool     chan struct{}
}

func NewEncoderServerV2(config ServerConfig, blobStore *blobstore.BlobStore, chunkWriter chunkstore.ChunkWriter, logger logging.Logger, prover encoding.Prover, metrics *Metrics) *EncoderServerV2 {
	return &EncoderServerV2{
		config:      config,
		blobStore:   blobStore,
		chunkWriter: chunkWriter,
		logger:      logger.With("component", "EncoderServer"),
		prover:      prover,
		metrics:     metrics,

		runningRequests: make(chan struct{}, config.MaxConcurrentRequests),
		requestPool:     make(chan struct{}, config.RequestPoolSize),
	}
}

func (s *EncoderServerV2) Start() error {
	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not start tcp listener: %v", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterEncoderServer(gs, s)

	// Register Server for Health Checks
	name := pb.Encoder_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	s.close = func() {
		err := listener.Close()
		if err != nil {
			log.Printf("failed to close listener: %v", err)
		}
		gs.GracefulStop()
	}

	s.logger.Info("port", s.config.GrpcPort, "address", listener.Addr().String(), "GRPC Listening")
	return gs.Serve(listener)
}

func (s *EncoderServerV2) EncodeBlobToChunkStore(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobToChunkStoreReply, error) {
	startTime := time.Now()

	// Rate limit
	select {
	case s.requestPool <- struct{}{}:
	default:
		// TODO: Now that we no longer pass the data directly, should we pass in blob size as part of the request
		s.metrics.IncrementRateLimitedBlobRequestNum(1)
		s.logger.Warn("rate limiting as request pool is full", "requestPoolSize", s.config.RequestPoolSize, "maxConcurrentRequests", s.config.MaxConcurrentRequests)
		return nil, errors.New("too many requests")
	}

	// Limit the number of concurrent requests
	s.runningRequests <- struct{}{}
	defer s.popRequest()
	if ctx.Err() != nil {
		s.metrics.IncrementCanceledBlobRequestNum(1)
		return nil, ctx.Err()
	}

	s.metrics.ObserveLatency("queuing", time.Since(startTime))
	reply, err := s.handleEncodingToChunkStore(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(1)
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(1)
	}
	s.metrics.ObserveLatency("total", time.Since(startTime))

	return reply, err
}

func (s *EncoderServerV2) handleEncodingToChunkStore(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobToChunkStoreReply, error) {
	begin := time.Now()
	if req.BlobHeaderHash == nil {
		return nil, errors.New("handleEncoding: missing blob header hash")
	}

	if req.EncodingParams == nil {
		return nil, errors.New("handleEncoding: missing encoding parameters")
	}

	// Convert bytes to BlobKey
	blobKey, err := bytesToBlobKey(req.GetBlobHeaderHash())
	if err != nil {
		return nil, fmt.Errorf("invalid blob header hash: %w", err)
	}

	// Convert BlobKey to hex string for storage lookup
	blobKeyHex := blobKey.Hex()

	// Get the data from the blob store
	data, err := s.blobStore.GetBlob(ctx, blobKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob from blob store: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("handleEncoding: missing data")
	}

	// Convert to core EncodingParams
	var encodingParams = encoding.EncodingParams{
		ChunkLength: uint64(req.GetEncodingParams().GetChunkLength()),
		NumChunks:   uint64(req.GetEncodingParams().GetNumChunks()),
	}

	s.logger.Info("Preparing to encode", "blobKey", blobKey.Hex(), "encodingParams", encodingParams)
	s.logger.Info("time to get blob", "duration", time.Since(begin))

	begin = time.Now()
	// Get the encoding frames for the data
	frames, err := s.prover.GetFrames(data, encodingParams)
	if err != nil {
		return nil, err
	}

	s.metrics.ObserveLatency("encoding", time.Since(begin))
	s.logger.Info("encoding", "duration", time.Since(begin))

	begin = time.Now()
	// Get the proofs and coefficients from the encoding frames
	proofs := make([]*encoding.Proof, 0, len(frames))
	coeffs := make([]*rs.Frame, 0, len(frames))
	for _, frame := range frames {
		proofs = append(proofs, &frame.Proof)
		frameWithCoeffs := &rs.Frame{
			Coeffs: frame.Coeffs,
		}
		coeffs = append(coeffs, frameWithCoeffs)
	}
	s.logger.Info("copying frames", "duration", time.Since(begin))

	begin = time.Now()
	// Upload the coefficients to the chunk store
	metadata, err := s.chunkWriter.PutChunkCoefficients(ctx, blobKey.Hex()+"_coeffs", coeffs)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunks: %w", err)
	}
	s.logger.Info("uploading coefficients", "duration", time.Since(begin))

	begin = time.Now()
	// Upload the proofs
	err = s.chunkWriter.PutChunkProofs(ctx, blobKey.Hex()+"_proofs", proofs)
	if err != nil {
		return nil, fmt.Errorf("failed to upload proofs: %w", err)
	}
	s.logger.Info("uploading proofs", "duration", time.Since(begin))

	return &pb.EncodeBlobToChunkStoreReply{
		FragmentInfo: &pb.FragmentInfo{
			TotalChunkSizeBytes: int32(metadata.FragmentSize),
			NumFragments:        int32(metadata.DataSize),
		},
	}, nil
}

func (s *EncoderServerV2) popRequest() {
	<-s.requestPool
	<-s.runningRequests
}

// Helper function to validate and convert bytes to BlobKey
func bytesToBlobKey(bytes []byte) (v2.BlobKey, error) {
	// Validate length
	if len(bytes) != 32 {
		return v2.BlobKey{}, fmt.Errorf("invalid blob hash length: expected 32 bytes, got %d", len(bytes))
	}

	var blobKey v2.BlobKey
	copy(blobKey[:], bytes)
	return blobKey, nil
}
