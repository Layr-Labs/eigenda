package encoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: Add EncodeMetrics
type Server struct {
	pb.UnimplementedEncoderServer

	config  ServerConfig
	logger  logging.Logger
	prover  encoding.Prover
	metrics *Metrics
	close   func()

	runningRequests chan struct{}
	requestPool     chan struct{}
}

func NewServer(config ServerConfig, logger logging.Logger, prover encoding.Prover, metrics *Metrics) *Server {
	return &Server{
		config:  config,
		logger:  logger.With("component", "EncoderServer"),
		prover:  prover,
		metrics: metrics,

		runningRequests: make(chan struct{}, config.MaxConcurrentRequests),
		requestPool:     make(chan struct{}, config.RequestPoolSize),
	}
}

func (s *Server) EncodeBlob(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
	select {
	case s.requestPool <- struct{}{}:
	default:
		s.metrics.IncrementRateLimitedBlobRequestNum(len(req.GetData()))
		s.logger.Warn("rate limiting as request pool is full", "requestPoolSize", s.config.RequestPoolSize, "maxConcurrentRequests", s.config.MaxConcurrentRequests)
		return nil, errors.New("too many requests")
	}
	s.runningRequests <- struct{}{}
	defer s.popRequest()

	if ctx.Err() != nil {
		s.metrics.IncrementCanceledBlobRequestNum(len(req.GetData()))
		return nil, ctx.Err()
	}

	reply, err := s.handleEncoding(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedBlobRequestNum(len(req.GetData()))
	} else {
		s.metrics.IncrementSuccessfulBlobRequestNum(len(req.GetData()))
	}
	return reply, err
}

func (s *Server) popRequest() {
	<-s.requestPool
	<-s.runningRequests
}

func (s *Server) handleEncoding(ctx context.Context, req *pb.EncodeBlobRequest) (*pb.EncodeBlobReply, error) {
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

	encodingTime := time.Since(begin)

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

	totalTime := time.Since(begin)
	s.metrics.TakeLatency(encodingTime, totalTime)

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

func (s *Server) Start() error {

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

func (s *Server) Close() {
	if s.close == nil {
		return
	}
	s.close()
}
