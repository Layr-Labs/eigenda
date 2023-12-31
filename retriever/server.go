package retriever

import (
	"context"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	gcommon "github.com/ethereum/go-ethereum/common"
)

type Server struct {
	pb.UnimplementedRetrieverServer

	config          *Config
	retrievalClient clients.RetrievalClient
	chainClient     eth.ChainClient
	indexedState    core.IndexedChainState
	logger          common.Logger
	metrics         *Metrics
}

func NewServer(
	config *Config,
	logger common.Logger,
	retrievalClient clients.RetrievalClient,
	encoder core.Encoder,
	indexedState core.IndexedChainState,
	chainClient eth.ChainClient,
) *Server {
	metrics := NewMetrics(config.MetricsConfig.HTTPPort, logger)

	return &Server{
		config:          config,
		retrievalClient: retrievalClient,
		chainClient:     chainClient,
		indexedState:    indexedState,
		logger:          logger,
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
		return nil, fmt.Errorf("got invalid batch header hash")
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
