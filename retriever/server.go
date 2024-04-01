package retriever

import (
	"context"
	"errors"

	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
)

type Server struct {
	pb.UnimplementedRetrieverServer

	config          *Config
	retrievalClient clients.RetrievalClient
	chainClient     eth.ChainClient
	indexedState    core.IndexedChainState
	logger          logging.Logger
	metrics         *Metrics
}

func NewServer(
	config *Config,
	logger logging.Logger,
	retrievalClient clients.RetrievalClient,
	verifier encoding.Verifier,
	indexedState core.IndexedChainState,
	chainClient eth.ChainClient,
) *Server {
	metrics := NewMetrics(config.MetricsConfig.HTTPPort, logger)

	return &Server{
		config:          config,
		retrievalClient: retrievalClient,
		chainClient:     chainClient,
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
	s.logger.Info("Received request: ", "BatchHeaderHash", req.GetBatchHeaderHash(), "BlobIndex", req.GetBlobIndex())
	s.metrics.IncrementRetrievalRequestCounter()
	if len(req.GetBatchHeaderHash()) != 32 {
		return nil, errors.New("got invalid batch header hash")
	}
	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], req.GetBatchHeaderHash())

	batchHeader, err := s.chainClient.FetchBatchHeader(ctx, gcommon.HexToAddress(s.config.EigenDAServiceManagerAddr), req.GetBatchHeaderHash())
	if err != nil {
		return nil, err
	}

	data, err := s.retrievalClient.RetrieveBlob(
		ctx,
		batchHeaderHash,
		req.GetBlobIndex(),
		uint(batchHeader.ReferenceBlockNumber),
		batchHeader.BlobHeadersRoot,
		core.QuorumID(req.GetQuorumId()))
	if err != nil {
		return nil, err
	}
	return &pb.BlobReply{
		Data: data,
	}, nil
}
