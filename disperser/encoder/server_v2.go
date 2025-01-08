package encoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common/healthcheck"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder/v2"
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

	runningRequests chan struct{}
	requestQueue    chan blobRequest

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
		config:          config,
		blobStore:       blobStore,
		chunkWriter:     chunkWriter,
		logger:          logger.With("component", "EncoderServerV2"),
		prover:          prover,
		metrics:         metrics,
		grpcMetrics:     grpcMetrics,
		runningRequests: make(chan struct{}, config.MaxConcurrentRequests),
		requestQueue:    make(chan blobRequest, config.RequestQueueSize),
		queueStats:      make(map[string]int),
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

	blobSize := req.GetBlobSize()
	sizeBucket := common.BlobSizeBucket(int(blobSize))

	// Rate limit
	select {
	case s.requestQueue <- blobRequest{blobSizeByte: int(blobSize)}:
		s.queueLock.Lock()
		s.queueStats[sizeBucket]++
		s.metrics.ObserveQueue(s.queueStats)
		s.queueLock.Unlock()
	default:
		s.metrics.IncrementRateLimitedBlobRequestNum(int(blobSize))
		s.logger.Warn("rate limiting as request queue is full", "requestQueueSize", s.config.RequestQueueSize, "maxConcurrentRequests", s.config.MaxConcurrentRequests)
		return nil, errors.New("too many requests")
	}

	// Limit the number of concurrent requests
	s.runningRequests <- struct{}{}
	defer s.popRequest()
	if ctx.Err() != nil {
		s.metrics.IncrementCanceledBlobRequestNum(int(blobSize))
		return nil, status.Error(codes.Canceled, "request was canceled")
	}

	s.metrics.ObserveLatency("queuing", time.Since(totalStart))
	reply, err := s.handleEncodingToChunkStore(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(int(blobSize))
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(int(blobSize))
	}

	return reply, err
}

func (s *EncoderServerV2) handleEncodingToChunkStore(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
	// Validate request first
	blobKey, encodingParams, err := s.validateAndParseRequest(req)
	if err != nil {
		return nil, err
	}

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

func (s *EncoderServerV2) popRequest() {
	blobRequest := <-s.requestQueue
	<-s.runningRequests
	s.queueLock.Lock()
	s.queueStats[common.BlobSizeBucket(blobRequest.blobSizeByte)]--
	s.metrics.ObserveQueue(s.queueStats)
	s.queueLock.Unlock()
}

func (s *EncoderServerV2) validateAndParseRequest(req *pb.EncodeBlobRequest) (corev2.BlobKey, encoding.EncodingParams, error) {
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

	blobKey, err := corev2.BytesToBlobKey(req.BlobKey)
	if err != nil {
		return blobKey, params, status.Errorf(codes.InvalidArgument, "invalid blob key: %v", err)
	}

	// Convert proto EncodingParams to our domain type
	params = encoding.EncodingParams{
		ChunkLength: req.EncodingParams.ChunkLength,
		NumChunks:   req.EncodingParams.NumChunks,
	}

	err = encoding.ValidateEncodingParams(params, s.prover.GetSRSOrder())
	if err != nil {
		return blobKey, params, status.Errorf(codes.InvalidArgument, "invalid encoding parameters: %v", err)
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
	if err := s.chunkWriter.PutChunkProofs(ctx, blobKey, proofs); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload chunk proofs: %v", err)
	}
	s.metrics.ObserveLatency("s3_upload_proofs", time.Since(storeStart))
	s.logger.Info("stored proofs", "duration", time.Since(storeStart).String())

	// Store coefficients
	coeffStart := time.Now()
	fragmentInfo, err := s.chunkWriter.PutChunkCoefficients(ctx, blobKey, coeffs)
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

func extractProofsAndCoeffs(frames []*encoding.Frame) ([]*encoding.Proof, []*rs.Frame) {
	proofs := make([]*encoding.Proof, len(frames))
	coeffs := make([]*rs.Frame, len(frames))

	for i, frame := range frames {
		proofs[i] = &frame.Proof
		coeffs[i] = &rs.Frame{Coeffs: frame.Coeffs}
	}
	return proofs, coeffs
}

func (s *EncoderServerV2) Close() {
	if s.close == nil {
		return
	}
	s.close()
}
