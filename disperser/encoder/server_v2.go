package encoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/encoder/v2"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
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
	grpcMetrics *grpcprom.ServerMetrics
	close       func()

	// This channel is used to limit the number of concurrent requests executed by the server. If its capacity
	// is smaller than the capacity of the backlogLimiter, then the server will process all enqueued requests
	// in parallel.
	concurrencyLimiter chan struct{}

	// This channel is used to limit the number of requests that can be enqueued. If this channel is at its limit
	// and new work is submitted, the server will immediately reject the new request.
	backlogLimiter chan struct{}

	queueStats map[string]int
	queueLock  sync.Mutex
}

func NewEncoderServerV2(
	config ServerConfig,
	blobStore *blobstore.BlobStore,
	chunkWriter chunkstore.ChunkWriter,
	logger logging.Logger,
	prover encoding.Prover,
	metrics *Metrics,
	grpcMetrics *grpcprom.ServerMetrics,
) *EncoderServerV2 {
	metrics.SetQueueCapacity(config.RequestQueueSize)

	return &EncoderServerV2{
		config:             config,
		blobStore:          blobStore,
		chunkWriter:        chunkWriter,
		logger:             logger.With("component", "EncoderServerV2"),
		prover:             prover,
		metrics:            metrics,
		grpcMetrics:        grpcMetrics,
		concurrencyLimiter: make(chan struct{}, config.MaxConcurrentRequests),
		backlogLimiter:     make(chan struct{}, config.RequestQueueSize),
		queueStats:         make(map[string]int),
	}
}

func (s *EncoderServerV2) Start() error {
	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not start tcp listener: %v", err)
	}

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(
			s.grpcMetrics.UnaryServerInterceptor(),
		),
	)
	reflection.Register(gs)
	pb.RegisterEncoderServer(gs, s)
	s.grpcMetrics.InitializeMetrics(gs)

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

func (s *EncoderServerV2) EncodeBlob(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
	totalStart := time.Now()
	defer func() {
		s.metrics.ObserveLatency("total", time.Since(totalStart))
	}()

	// Validate the request.
	blobKey, encodingParams, err := s.validateAndParseRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	blobSize := int(req.GetBlobSize())

	// If we have too large of a backlog, refuse to accept new work.
	err = s.pushBacklogLimiter(blobSize)
	if err != nil {
		return nil, err
	}
	defer s.popBacklogLimiter(blobSize)

	// Limit the number of concurrent requests.
	err = s.pushConcurrencyLimiter(ctx, blobSize)
	if err != nil {
		return nil, err
	}
	defer s.popConcurrencyLimiter()

	s.metrics.ObserveLatency("queuing", time.Since(totalStart))
	reply, err := s.handleEncodingToChunkStore(ctx, blobKey, encodingParams)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(blobSize)
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(blobSize)
	}

	return reply, err
}

func (s *EncoderServerV2) handleEncodingToChunkStore(ctx context.Context, blobKey corev2.BlobKey, encodingParams encoding.EncodingParams) (*pb.EncodeBlobReply, error) {
	s.logger.Info("Preparing to encode", "blobKey", blobKey.Hex(), "encodingParams", encodingParams)

	// Check if the blob has already been encoded
	if s.config.PreventReencoding && s.chunkWriter.ProofExists(ctx, blobKey) {
		coefExist, fragmentInfo := s.chunkWriter.CoefficientsExists(ctx, blobKey)
		if coefExist {
			s.logger.Info("blob already encoded", "blobKey", blobKey.Hex())
			return &pb.EncodeBlobReply{
				FragmentInfo: &pb.FragmentInfo{
					TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
					FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
				},
			}, nil
		}
	}

	// Fetch blob data
	fetchStart := time.Now()
	data, err := s.blobStore.GetBlob(ctx, blobKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get blob from blob store: %v", err)
	}
	if len(data) == 0 {
		return nil, status.Error(codes.NotFound, "blob length is zero")
	}
	s.metrics.ObserveLatency("s3_download", time.Since(fetchStart))
	s.logger.Info("fetched blob", "duration", time.Since(fetchStart).String())

	// Encode the data
	encodingStart := time.Now()
	frames, err := s.prover.GetFrames(data, encodingParams)
	if err != nil {
		s.logger.Error("failed to encode frames", "error", err)
		return nil, status.Errorf(codes.Internal, "encoding failed: %v", err)
	}
	s.metrics.ObserveLatency("encoding", time.Since(encodingStart))
	s.logger.Info("encoding frames", "duration", time.Since(encodingStart).String())

	return s.processAndStoreResults(ctx, blobKey, frames)
}

// pushBacklogLimiter pushes a token to the backlog limiter and increments the queue stats accordingly.
// If there is no capacity in the backlog limiter, an error is returned.
func (s *EncoderServerV2) pushBacklogLimiter(blobSizeBytes int) error {
	sizeBucket := common.BlobSizeBucket(blobSizeBytes)

	select {
	case s.backlogLimiter <- struct{}{}:
		s.queueLock.Lock()
		s.queueStats[sizeBucket]++
		s.metrics.ObserveQueue(s.queueStats)
		s.queueLock.Unlock()

		return nil
	default:
		s.metrics.IncrementRateLimitedBlobRequestNum(blobSizeBytes)
		s.logger.Warn("rate limiting as request queue is full",
			"requestQueueSize", s.config.RequestQueueSize,
			"maxConcurrentRequests", s.config.MaxConcurrentRequests)
		return api.NewErrorResourceExhausted(fmt.Sprintf(
			"request queue is full, max queue size: %d", s.config.RequestQueueSize))
	}
}

// popBacklogLimiter pops a token from the backlog limiter and decrements the queue stats accordingly.
func (s *EncoderServerV2) popBacklogLimiter(blobSizeBytes int) {
	<-s.backlogLimiter
	s.queueLock.Lock()
	s.queueStats[common.BlobSizeBucket(blobSizeBytes)]--
	s.metrics.ObserveQueue(s.queueStats)
	s.queueLock.Unlock()
}

// pushConcurrencyLimiter pushes a token to the concurrency limiter.
func (s *EncoderServerV2) pushConcurrencyLimiter(ctx context.Context, blobSizeBytes int) error {
	select {
	case s.concurrencyLimiter <- struct{}{}:
		return nil
	case <-ctx.Done():
		s.metrics.IncrementCanceledBlobRequestNum(blobSizeBytes)
		return status.Error(codes.Canceled, "request was canceled")
	}
}

// popConcurrencyLimiter pops a token from the concurrency limiter.
func (s *EncoderServerV2) popConcurrencyLimiter() {
	<-s.concurrencyLimiter
}

func (s *EncoderServerV2) validateAndParseRequest(req *pb.EncodeBlobRequest) (corev2.BlobKey, encoding.EncodingParams, error) {
	// Create zero values for return types
	var (
		blobKey corev2.BlobKey
		params  encoding.EncodingParams
	)

	if req == nil {
		return blobKey, params, errors.New("request cannot be nil")
	}

	if req.BlobKey == nil {
		return blobKey, params, errors.New("blob key cannot be nil")
	}

	if req.GetEncodingParams() == nil {
		return blobKey, params, errors.New("encoding parameters cannot be nil")
	}

	// Since these are uint32 in the proto, we only need to check for positive values
	if req.GetEncodingParams().GetChunkLength() == 0 {
		return blobKey, params, errors.New("chunk length must be greater than zero")
	}
	if req.GetEncodingParams().GetChunkLength()&(req.GetEncodingParams().GetChunkLength()-1) != 0 {
		return blobKey, params, errors.New("chunk length must be power of 2")
	}

	if req.GetEncodingParams().GetNumChunks() == 0 {
		return blobKey, params, errors.New("number of chunks must be greater than zero")
	}

	if req.GetBlobSize() == 0 || uint64(encoding.GetBlobLength(uint(req.GetBlobSize()))) > req.GetEncodingParams().GetChunkLength()*req.GetEncodingParams().GetNumChunks() {
		return blobKey, params, errors.New("blob size is invalid")
	}

	blobKey, err := corev2.BytesToBlobKey(req.GetBlobKey())
	if err != nil {
		return blobKey, params, fmt.Errorf("invalid blob key: %v", err)
	}

	// Convert proto EncodingParams to our domain type
	params = encoding.EncodingParams{
		ChunkLength: req.GetEncodingParams().GetChunkLength(),
		NumChunks:   req.GetEncodingParams().GetNumChunks(),
	}

	err = encoding.ValidateEncodingParams(params, s.prover.GetSRSOrder())
	if err != nil {
		return blobKey, params, fmt.Errorf("invalid encoding parameters: %v", err)
	}

	return blobKey, params, nil
}

func (s *EncoderServerV2) processAndStoreResults(ctx context.Context, blobKey corev2.BlobKey, frames []*encoding.Frame) (*pb.EncodeBlobReply, error) {
	// Store proofs
	storeStart := time.Now()
	defer func() {
		s.metrics.ObserveLatency("process_and_store_results", time.Since(storeStart))
	}()

	proofs, coeffs := extractProofsAndCoeffs(frames)
	if err := s.chunkWriter.PutFrameProofs(ctx, blobKey, proofs); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload chunk proofs: %v", err)
	}
	s.metrics.ObserveLatency("s3_upload_proofs", time.Since(storeStart))
	s.logger.Info("stored proofs", "duration", time.Since(storeStart).String())

	// Store coefficients
	coeffStart := time.Now()
	fragmentInfo, err := s.chunkWriter.PutFrameCoefficients(ctx, blobKey, coeffs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload chunk coefficients: %v", err)
	}
	s.metrics.ObserveLatency("s3_upload_coefficients", time.Since(coeffStart))
	s.logger.Info("stored coefficients", "duration", time.Since(coeffStart).String())

	return &pb.EncodeBlobReply{
		FragmentInfo: &pb.FragmentInfo{
			TotalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			FragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		},
	}, nil
}

func extractProofsAndCoeffs(frames []*encoding.Frame) ([]*encoding.Proof, []rs.FrameCoeffs) {
	proofs := make([]*encoding.Proof, len(frames))
	coeffs := make([]rs.FrameCoeffs, len(frames))

	for i, frame := range frames {
		proofs[i] = &frame.Proof
		coeffs[i] = frame.Coeffs
	}
	return proofs, coeffs
}

func (s *EncoderServerV2) Close() {
	if s.close == nil {
		return
	}
	s.close()
}
