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
	commonpprof "github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/disperser/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type EncoderServer struct {
	pb.UnimplementedEncoderServer

	config      ServerConfig
	logger      logging.Logger
	prover      encoding.Prover
	metrics     *Metrics
	grpcMetrics *grpcprom.ServerMetrics
	close       func()

	runningRequests chan struct{}
	requestPool     chan blobRequest

	queueStats map[string]int
	queueLock  sync.Mutex
}

type blobRequest struct {
	blobSizeByte int
}

func NewEncoderServer(config ServerConfig, logger logging.Logger, prover encoding.Prover, metrics *Metrics, grpcMetrics *grpcprom.ServerMetrics) *EncoderServer {
	// Set initial queue capacity metric
	metrics.SetQueueCapacity(config.RequestPoolSize)

	return &EncoderServer{
		config:      config,
		logger:      logger.With("component", "EncoderServer"),
		prover:      prover,
		metrics:     metrics,
		grpcMetrics: grpcMetrics,

		runningRequests: make(chan struct{}, config.MaxConcurrentRequests),
		requestPool:     make(chan blobRequest, config.RequestPoolSize),
		queueStats:      make(map[string]int),
	}
}

func (s *EncoderServer) Start() error {
	pprofProfiler := commonpprof.NewPprofProfiler(s.config.PprofHttpPort, s.logger)
	if s.config.EnablePprof {
		go pprofProfiler.Start()
		s.logger.Info("Enabled pprof for encoder server", "port", s.config.PprofHttpPort)
	}

	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not start tcp listener: %v", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt,
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

func (s *EncoderServer) EncodeBlob(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
	startTime := time.Now()
	blobSize := len(req.GetData())
	sizeBucket := common.BlobSizeBucket(blobSize)

	select {
	case s.requestPool <- blobRequest{blobSizeByte: blobSize}:
		s.queueLock.Lock()
		s.queueStats[sizeBucket]++
		s.metrics.ObserveQueue(s.queueStats)
		s.queueLock.Unlock()
	default:
		s.metrics.IncrementRateLimitedBlobRequestNum(blobSize)
		s.logger.Warn("rate limiting as request pool is full", "requestPoolSize", s.config.RequestPoolSize, "maxConcurrentRequests", s.config.MaxConcurrentRequests)
		return nil, errors.New("too many requests")
	}

	s.runningRequests <- struct{}{}
	defer s.popRequest()

	if ctx.Err() != nil {
		s.metrics.IncrementCanceledBlobRequestNum(blobSize)
		return nil, ctx.Err()
	}

	s.metrics.ObserveLatency("queuing", time.Since(startTime))
	reply, err := s.handleEncoding(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(blobSize)
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(blobSize)
	}
	s.metrics.ObserveLatency("total", time.Since(startTime))

	return reply, err
}

func (s *EncoderServer) popRequest() {
	blobRequest := <-s.requestPool
	<-s.runningRequests
	s.queueLock.Lock()
	s.queueStats[common.BlobSizeBucket(blobRequest.blobSizeByte)]--
	s.metrics.ObserveQueue(s.queueStats)
	s.queueLock.Unlock()
}

func (s *EncoderServer) handleEncoding(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
	begin := time.Now()

	if len(req.Data) == 0 {
		return nil, errors.New("handleEncoding: missing data")
	}

	if req.EncodingParams == nil {
		return nil, errors.New("handleEncoding: missing encoding parameters")
	}

	// Convert to core EncodingParams
	var encodingParams = encoding.EncodingParams{
		ChunkLength: uint64(req.GetEncodingParams().GetChunkLength()),
		NumChunks:   uint64(req.GetEncodingParams().GetNumChunks()),
	}

	commits, chunks, err := s.prover.EncodeAndProve(req.GetData(), encodingParams)
	if err != nil {
		return nil, err
	}

	s.metrics.ObserveLatency("encoding", time.Since(begin))
	begin = time.Now()

	commitData, err := commits.Commitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthCommitData, err := commits.LengthCommitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthProofData, err := commits.LengthProof.Serialize()
	if err != nil {
		return nil, err
	}

	var chunksData [][]byte
	var format pb.ChunkEncodingFormat
	if s.config.EnableGnarkChunkEncoding {
		format = pb.ChunkEncodingFormat_GNARK
	} else {
		format = pb.ChunkEncodingFormat_GOB
	}

	for _, chunk := range chunks {
		var chunkSerialized []byte
		if s.config.EnableGnarkChunkEncoding {
			chunkSerialized, err = chunk.SerializeGnark()
		} else {
			chunkSerialized, err = chunk.Serialize()
		}
		if err != nil {
			return nil, err
		}
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}

	s.metrics.ObserveLatency("serialization", time.Since(begin))

	return &pb.EncodeBlobReply{
		Commitment: &pb.BlobCommitment{
			Commitment:       commitData,
			LengthCommitment: lengthCommitData,
			LengthProof:      lengthProofData,
			Length:           uint32(commits.Length),
		},
		Chunks:              chunksData,
		ChunkEncodingFormat: format,
	}, nil
}

func (s *EncoderServer) Close() {
	if s.close == nil {
		return
	}
	s.close()
}
