package encoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/common/healthcheck"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

	gs := grpc.NewServer()
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

func (s *EncoderServerV2) EncodeBlobToChunkStore(ctx context.Context, req *pb.EncodeBlobToChunkStoreRequest) (*pb.EncodeBlobToChunkStoreReply, error) {
	totalStart := time.Now()
	defer func() {
		s.metrics.ObserveLatency("total", time.Since(totalStart))
	}()

	// Rate limit
	select {
	case s.requestPool <- struct{}{}:
	default:
		// TODO: Now that we no longer pass the data directly, should we pass in blob size as part of the request?
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

	s.metrics.ObserveLatency("queuing", time.Since(totalStart))
	reply, err := s.handleEncodingToChunkStore(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(1)
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(1)
	}

	return reply, err
}

func (s *EncoderServerV2) handleEncodingToChunkStore(ctx context.Context, req *pb.EncodeBlobToChunkStoreRequest) (*pb.EncodeBlobToChunkStoreReply, error) {
	// Validate request first
	blobKey, encodingParams, err := s.validateAndParseRequest(req)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Preparing to encode", "blobKey", blobKey, "encodingParams", encodingParams)

	// Fetch blob data
	fetchStart := time.Now()
	data, err := s.blobStore.GetBlob(ctx, blobKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob from blob store: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("handleEncodingToChunkStore: missing data")
	}
	s.logger.Info("fetched blob", "duration", time.Since(fetchStart))

	// Encode the data
	encodingStart := time.Now()
	frames, err := s.prover.GetFrames(data, encodingParams)
	if err != nil {
		return nil, fmt.Errorf("frame encoding failed: %w", err)
	}
	s.logger.Info("encoding frames", "duration", time.Since(encodingStart))

	// Process and store results
	return s.processAndStoreResults(ctx, blobKey, frames)
}

func (s *EncoderServerV2) popRequest() {
	<-s.requestPool
	<-s.runningRequests
}

func (s *EncoderServerV2) validateAndParseRequest(req *pb.EncodeBlobToChunkStoreRequest) (corev2.BlobKey, encoding.EncodingParams, error) {
	// Create zero values for return types
	var (
		blobKey corev2.BlobKey
		params  encoding.EncodingParams
	)

	if req == nil {
		return blobKey, params, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.BlobKey == nil {
		return blobKey, params, status.Error(codes.InvalidArgument, "blob key cannot be nil")
	}

	if req.EncodingParams == nil {
		return blobKey, params, status.Error(codes.InvalidArgument, "encoding parameters cannot be nil")
	}

	// Since these are uint32 in the proto, we only need to check for positive values
	if req.EncodingParams.ChunkLength == 0 {
		return blobKey, params, status.Error(codes.InvalidArgument, "chunk length must be greater than zero")
	}

	if req.EncodingParams.NumChunks == 0 {
		return blobKey, params, status.Error(codes.InvalidArgument, "number of chunks must be greater than zero")
	}

	blobKey, err := bytesToBlobKey(req.BlobKey)
	if err != nil {
		return blobKey, params, status.Errorf(codes.InvalidArgument, "invalid blob key: %v", err)
	}

	// Convert proto EncodingParams to our domain type
	params = encoding.EncodingParams{
		ChunkLength: req.EncodingParams.ChunkLength,
		NumChunks:   req.EncodingParams.NumChunks,
	}

	return blobKey, params, nil
}

func (s *EncoderServerV2) processAndStoreResults(ctx context.Context, blobKey corev2.BlobKey, frames []*encoding.Frame) (*pb.EncodeBlobToChunkStoreReply, error) {
	proofs, coeffs := extractProofsAndCoeffs(frames)

	// Store proofs
	storeStart := time.Now()
	if err := s.chunkWriter.PutChunkProofs(ctx, blobKey, proofs); err != nil {
		return nil, fmt.Errorf("failed to upload chunk proofs: %w", err)
	}
	s.logger.Info("stored proofs", "duration", time.Since(storeStart))

	// Store coefficients
	coeffStart := time.Now()
	fragmentInfo, err := s.chunkWriter.PutChunkCoefficients(ctx, blobKey, coeffs)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunk coefficients: %w", err)
	}
	s.logger.Info("stored coefficients", "duration", time.Since(coeffStart))

	return &pb.EncodeBlobToChunkStoreReply{
		FragmentInfo: &pb.FragmentInfo{
			TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		},
	}, nil
}

// Helper function to validate and convert bytes to BlobKey
func bytesToBlobKey(bytes []byte) (corev2.BlobKey, error) {
	// Validate length
	if len(bytes) != 32 {
		return corev2.BlobKey{}, fmt.Errorf("invalid blob key length: expected 32 bytes, got %d", len(bytes))
	}

	var blobKey corev2.BlobKey
	copy(blobKey[:], bytes)
	return blobKey, nil
}

func extractProofsAndCoeffs(frames []*encoding.Frame) ([]*encoding.Proof, []*rs.Frame) {
	proofs := make([]*encoding.Proof, len(frames))
	coeffs := make([]*rs.Frame, len(frames))

	for i, frame := range frames {
		proofs[i] = &frame.Proof
		coeffs[i] = &rs.Frame{Coeffs: frame.Coeffs}
	}
	return proofs, coeffs
}
