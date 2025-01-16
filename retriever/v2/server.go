package v2

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/retriever"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type Config = retriever.Config

type Server struct {
	pb.UnimplementedRetrieverServer

	config          *Config
	retrievalClient clients.RetrievalClient
	indexedState    core.IndexedChainState
	logger          logging.Logger
	metrics         *retriever.Metrics
}

func NewServer(
	config *Config,
	logger logging.Logger,
	retrievalClient clients.RetrievalClient,
	indexedState core.IndexedChainState,
) *Server {
	metrics := retriever.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	return &Server{
		config:          config,
		retrievalClient: retrievalClient,
		indexedState:    indexedState,
		logger:          logger.With("component", "RetrieverServer"),
		metrics:         metrics,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.metrics.Start(ctx)
	return s.indexedState.Start(ctx)
}

func (s *Server) RetrieveBlob(ctx context.Context, req *pb.BlobRequest) (*pb.BlobReply, error) {
	if req.GetBlobHeader() == nil {
		return nil, errors.New("blob header is nil")
	}
	if req.GetReferenceBlockNumber() == 0 {
		return nil, errors.New("reference block number is 0")
	}

	blobHeader, err := corev2.BlobHeaderFromProtobuf(req.GetBlobHeader())
	if err != nil {
		return nil, err
	}

	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return nil, err
	}

	s.logger.Info("Received request: ", "blobKey", hex.EncodeToString(blobKey[:]), "referenceBlockNumber", req.GetReferenceBlockNumber(), "quorumId", req.GetQuorumId())
	s.metrics.IncrementRetrievalRequestCounter()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	data, err := s.retrievalClient.GetBlob(ctxWithTimeout, blobHeader, uint64(req.GetReferenceBlockNumber()), core.QuorumID(req.GetQuorumId()))
	if err != nil {
		return nil, err
	}
	restored := bytes.TrimRight(data, "\x00")
	restored = codec.RemoveEmptyByteFromPaddedBytes(restored)

	return &pb.BlobReply{
		Data: restored,
	}, nil
}
