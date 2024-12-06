package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbv1 "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OnchainState struct {
	QuorumCount           uint8
	RequiredQuorums       []core.QuorumID
	BlobVersionParameters *corev2.BlobVersionParameterMap
	TTL                   time.Duration
}

// Include disperser v1 protos to support grpcurl/reflection of v1 APIs
type DispersalServerV1 struct {
	pbv1.UnimplementedDisperserServer
}

type DispersalServerV2 struct {
	pb.UnimplementedDisperserServer

	serverConfig      disperser.ServerConfig
	blobStore         *blobstore.BlobStore
	blobMetadataStore *blobstore.BlobMetadataStore
	meterer           *meterer.Meterer

	chainReader   core.Reader
	authenticator corev2.BlobRequestAuthenticator
	prover        encoding.Prover
	logger        logging.Logger

	// state
	onchainState                atomic.Pointer[OnchainState]
	maxNumSymbolsPerBlob        uint64
	onchainStateRefreshInterval time.Duration
}

// NewDispersalServerV2 creates a new Server struct with the provided parameters.
func NewDispersalServerV2(
	serverConfig disperser.ServerConfig,
	rateConfig RateConfig,
	blobStore *blobstore.BlobStore,
	blobMetadataStore *blobstore.BlobMetadataStore,
	chainReader core.Reader,
	ratelimiter common.RateLimiter,
	meterer *meterer.Meterer,
	authenticator corev2.BlobRequestAuthenticator,
	prover encoding.Prover,
	maxNumSymbolsPerBlob uint64,
	onchainStateRefreshInterval time.Duration,
	_logger logging.Logger,
) *DispersalServerV2 {
	logger := _logger.With("component", "DispersalServerV2")

	return &DispersalServerV2{
		serverConfig:      serverConfig,
		blobStore:         blobStore,
		blobMetadataStore: blobMetadataStore,

		chainReader:   chainReader,
		authenticator: authenticator,
		meterer:       meterer,
		prover:        prover,
		logger:        logger,

		maxNumSymbolsPerBlob:        maxNumSymbolsPerBlob,
		onchainStateRefreshInterval: onchainStateRefreshInterval,
	}
}

func (s *DispersalServerV2) Start(ctx context.Context) error {
	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB

	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Unimplemented v1 server for grpcurl/reflection support
	pbv1.RegisterDisperserServer(gs, &DispersalServerV1{})

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	if err := s.RefreshOnchainState(ctx); err != nil {
		return fmt.Errorf("failed to refresh onchain quorum state: %w", err)
	}

	go func() {
		ticker := time.NewTicker(s.onchainStateRefreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.RefreshOnchainState(ctx); err != nil {
					s.logger.Error("failed to refresh onchain quorum state", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	s.logger.Info("GRPC Listening", "port", s.serverConfig.GrpcPort, "address", listener.Addr().String())

	if err := gs.Serve(listener); err != nil {
		return errors.New("could not start GRPC server")
	}

	return nil
}

func (s *DispersalServerV2) GetBlobCommitment(ctx context.Context, req *pb.BlobCommitmentRequest) (*pb.BlobCommitmentReply, error) {
	if s.prover == nil {
		return nil, api.NewErrorUnimplemented()
	}
	blobSize := len(req.GetData())
	if blobSize == 0 {
		return nil, api.NewErrorInvalidArg("data is empty")
	}
	if uint64(blobSize) > s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("blob size cannot exceed %v bytes", s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL))
	}
	c, err := s.prover.GetCommitmentsForPaddedLength(req.GetData())
	if err != nil {
		return nil, api.NewErrorInternal("failed to get commitments")
	}
	commitment, err := c.Commitment.Serialize()
	if err != nil {
		return nil, api.NewErrorInternal("failed to serialize commitment")
	}
	lengthCommitment, err := c.LengthCommitment.Serialize()
	if err != nil {
		return nil, api.NewErrorInternal("failed to serialize length commitment")
	}
	lengthProof, err := c.LengthProof.Serialize()
	if err != nil {
		return nil, api.NewErrorInternal("failed to serialize length proof")
	}

	return &pb.BlobCommitmentReply{
		BlobCommitment: &pbcommon.BlobCommitment{
			Commitment:       commitment,
			LengthCommitment: lengthCommitment,
			LengthProof:      lengthProof,
			Length:           uint32(c.Length),
		}}, nil
}

// refreshOnchainState refreshes the onchain quorum state.
// It should be called periodically to keep the state up to date.
// **Note** that there is no lock. If the state is being updated concurrently, it may lead to inconsistent state.
func (s *DispersalServerV2) RefreshOnchainState(ctx context.Context) error {
	s.logger.Debug("Refreshing onchain quorum state")

	currentBlock, err := s.chainReader.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	quorumCount, err := s.chainReader.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		return fmt.Errorf("failed to get quorum count: %w", err)
	}
	requiredQuorums, err := s.chainReader.GetRequiredQuorumNumbers(ctx, currentBlock)
	if err != nil {
		return fmt.Errorf("failed to get required quorum numbers: %w", err)
	}

	blockStaleMeasure, err := s.chainReader.GetBlockStaleMeasure(ctx)
	if err != nil {
		return fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
	}
	storeDurationBlocks, err := s.chainReader.GetStoreDurationBlocks(ctx)
	if err != nil || storeDurationBlocks == 0 {
		return fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
	}

	blobParams, err := s.chainReader.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get blob version parameters: %w", err)
	}
	onchainState := &OnchainState{
		QuorumCount:           quorumCount,
		RequiredQuorums:       requiredQuorums,
		BlobVersionParameters: v2.NewBlobVersionParameterMap(blobParams),
		TTL:                   time.Duration((storeDurationBlocks+blockStaleMeasure)*12) * time.Second,
	}

	s.onchainState.Store(onchainState)

	return nil
}

func (s *DispersalServerV2) GetPaymentState(ctx context.Context, req *pb.GetPaymentStateRequest) (*pb.GetPaymentStateReply, error) {
	// validate the signature
	if err := s.authenticator.AuthenticatePaymentStateRequest(req.GetSignature(), req.GetAccountId()); err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
	}
	// on-chain global payment parameters
	globalSymbolsPerSecond := s.meterer.ChainPaymentState.GetGlobalSymbolsPerSecond()
	minNumSymbols := s.meterer.ChainPaymentState.GetMinNumSymbols()
	pricePerSymbol := s.meterer.ChainPaymentState.GetPricePerSymbol()
	reservationWindow := s.meterer.ChainPaymentState.GetReservationWindow()

	// off-chain account specific payment state
	now := uint64(time.Now().Unix())
	currentBinIndex := meterer.GetBinIndex(now, reservationWindow)
	binRecords, err := s.meterer.OffchainStore.GetBinRecords(ctx, req.AccountId, currentBinIndex)
	if err != nil {
		return nil, api.NewErrorNotFound("failed to get active reservation")
	}
	largestCumulativePayment, err := s.meterer.OffchainStore.GetLargestCumulativePayment(ctx, req.AccountId)
	if err != nil {
		return nil, api.NewErrorNotFound("failed to get largest cumulative payment")
	}
	// on-Chain account state
	reservation, err := s.meterer.ChainPaymentState.GetActiveReservationByAccount(ctx, req.AccountId)
	if err != nil {
		return nil, api.NewErrorNotFound("failed to get active reservation")
	}
	onDemandPayment, err := s.meterer.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, req.AccountId)
	if err != nil {
		return nil, api.NewErrorNotFound("failed to get on-demand payment")
	}

	paymentGlobalParams := pb.PaymentGlobalParams{
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
		ReservationWindow:      reservationWindow,
	}

	quorumNumbers := make([]uint32, len(reservation.QuorumNumbers))
	for i, v := range reservation.QuorumNumbers {
		quorumNumbers[i] = uint32(v)
	}
	// build reply
	reply := &pb.GetPaymentStateReply{
		PaymentGlobalParams: &paymentGlobalParams,
		BinRecords:          binRecords[:],
		Reservation: &pb.Reservation{
			SymbolsPerSecond: reservation.SymbolsPerSec,
			StartTimestamp:   uint32(reservation.StartTimestamp),
			EndTimestamp:     uint32(reservation.EndTimestamp),
			QuorumNumbers:    quorumNumbers,
		},
		CumulativePayment:        largestCumulativePayment.Bytes(),
		OnchainCumulativePayment: onDemandPayment.CumulativePayment.Bytes(),
	}
	return reply, nil
}
