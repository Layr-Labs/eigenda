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
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
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
	blobMetadataStore blobstore.MetadataStore
	meterer           *meterer.Meterer

	chainReader              core.Reader
	blobRequestAuthenticator corev2.BlobRequestAuthenticator
	prover                   encoding.Prover
	logger                   logging.Logger

	// state
	onchainState                atomic.Pointer[OnchainState]
	maxNumSymbolsPerBlob        uint64
	onchainStateRefreshInterval time.Duration

	metricsConfig disperser.MetricsConfig
	metrics       *metricsV2

	ntpClock *core.NTPSyncedClock
	// ReservedOnly mode doesn't support on-demand payments
	// This would be removed with decentralized ratelimiting
	ReservedOnly bool
}

// NewDispersalServerV2 creates a new Server struct with the provided parameters.
func NewDispersalServerV2(
	serverConfig disperser.ServerConfig,
	blobStore *blobstore.BlobStore,
	blobMetadataStore blobstore.MetadataStore,
	chainReader core.Reader,
	meterer *meterer.Meterer,
	blobRequestAuthenticator corev2.BlobRequestAuthenticator,
	prover encoding.Prover,
	maxNumSymbolsPerBlob uint64,
	onchainStateRefreshInterval time.Duration,
	_logger logging.Logger,
	registry *prometheus.Registry,
	metricsConfig disperser.MetricsConfig,
	ntpClock *core.NTPSyncedClock,
	ReservedOnly bool,
) (*DispersalServerV2, error) {
	if serverConfig.GrpcPort == "" {
		return nil, errors.New("grpc port is required")
	}
	if blobStore == nil {
		return nil, errors.New("blob store is required")
	}
	if blobMetadataStore == nil {
		return nil, errors.New("blob metadata store is required")
	}
	if chainReader == nil {
		return nil, errors.New("chain reader is required")
	}
	if blobRequestAuthenticator == nil {
		return nil, errors.New("blobRequestAuthenticator is required")
	}
	if prover == nil {
		return nil, errors.New("prover is required")
	}
	if maxNumSymbolsPerBlob == 0 {
		return nil, errors.New("maxNumSymbolsPerBlob is required")
	}
	if _logger == nil {
		return nil, errors.New("logger is required")
	}

	logger := _logger.With("component", "DispersalServerV2")

	return &DispersalServerV2{
		serverConfig:      serverConfig,
		blobStore:         blobStore,
		blobMetadataStore: blobMetadataStore,

		chainReader:              chainReader,
		blobRequestAuthenticator: blobRequestAuthenticator,
		meterer:                  meterer,
		prover:                   prover,
		logger:                   logger,

		maxNumSymbolsPerBlob:        maxNumSymbolsPerBlob,
		onchainStateRefreshInterval: onchainStateRefreshInterval,

		metricsConfig: metricsConfig,
		metrics:       newAPIServerV2Metrics(registry, metricsConfig, logger),

		ntpClock:     ntpClock,
		ReservedOnly: ReservedOnly,
	}, nil
}

func (s *DispersalServerV2) Start(ctx context.Context) error {
	// Start the metrics server
	if s.metricsConfig.EnableMetrics {
		s.metrics.Start(context.Background())
	}

	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB

	gs := grpc.NewServer(opt, s.metrics.grpcServerOption)
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
				s.logger.Debug("Refreshed onchain quorum state", "onchainState", s.onchainState.Load())
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
	start := time.Now()
	defer func() {
		s.metrics.reportGetBlobCommitmentLatency(time.Since(start))
	}()

	if s.prover == nil {
		return nil, api.NewErrorUnimplemented()
	}
	blobSize := uint(len(req.GetBlob()))
	if blobSize == 0 {
		return nil, api.NewErrorInvalidArg("data is empty")
	}
	if uint64(encoding.GetBlobLengthPowerOf2(blobSize)) > s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("blob size cannot exceed %v bytes",
			s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL))
	}
	c, err := s.prover.GetCommitmentsForPaddedLength(req.GetBlob())
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

// GetPaymentState returns the payment state for a given account and the related on-chain parameters
// Deprecated: Use GetPaymentStateForAllQuorums instead.
func (s *DispersalServerV2) GetPaymentState(ctx context.Context, req *pb.GetPaymentStateRequest) (*pb.GetPaymentStateReply, error) {
	if s.meterer == nil {
		return nil, errors.New("payment meterer is not enabled")
	}
	start := time.Now()
	defer func() {
		s.metrics.reportGetPaymentStateLatency(time.Since(start))
	}()

	if !gethcommon.IsHexAddress(req.AccountId) {
		return nil, api.NewErrorInvalidArg("invalid account ID")
	}

	accountID := gethcommon.HexToAddress(req.AccountId)

	// validate the signature
	if err := s.blobRequestAuthenticator.AuthenticatePaymentStateRequest(accountID, req); err != nil {
		s.logger.Debug("failed to validate signature", "err", err, "accountID", accountID)
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
	}
	// on-chain global payment parameters
	params, err := s.meterer.ChainPaymentState.GetPaymentGlobalParams()
	if err != nil {
		s.logger.Error("failed to get payment global params", "err", err)
		return nil, api.NewErrorInternal("failed to get payment parameters")
	}

	paymentConfig, protocolConfig, err := params.GetQuorumConfigs(core.QuorumID(0))
	if err != nil {
		s.logger.Error("failed to get payment config for quorum 0", "err", err)
		return nil, api.NewErrorInternal("failed to get payment configuration")
	}

	globalSymbolsPerSecond := paymentConfig.OnDemandSymbolsPerSecond
	minNumSymbols := protocolConfig.MinNumSymbols
	pricePerSymbol := paymentConfig.OnDemandPricePerSymbol
	reservationWindow := protocolConfig.ReservationRateLimitWindow

	// off-chain account specific payment state
	now := time.Now().Unix()
	quorumIds := s.onchainState.Load().getAllQuorumIds()
	// Constructing the strictest period records across all quorums; a client should migrate to GetPaymentStateForAllQuorums for more precise state
	// TODO(hopeyen): remove this in a subsequent PR. The logic here is complicated and only temporary
	periods := make([]uint64, len(quorumIds))
	for i, quorumId := range quorumIds {
		_, quorumProtocolConfig, err := params.GetQuorumConfigs(quorumId)
		if err != nil {
			s.logger.Error("failed to get protocol config for quorum", "quorumId", quorumId, "err", err)
			return nil, api.NewErrorInternal("failed to get quorum protocol configuration")
		}
		periods[i] = meterer.GetReservationPeriod(now, quorumProtocolConfig.ReservationRateLimitWindow)
	}
	records, err := s.meterer.MeteringStore.GetPeriodRecords(ctx, accountID, quorumIds, periods, 3)
	if err != nil {
		s.logger.Debug("failed to get period records, use zero value", "err", err, "accountID", accountID)
	}

	highestPeriodRecords := make([]*pb.PeriodRecord, meterer.MinNumBins)
	for _, records := range records {
		for _, record := range records.Records {
			idx := record.Index % uint32(meterer.MinNumBins)
			if highestPeriodRecords[idx] == nil || record.Usage > highestPeriodRecords[idx].Usage {
				highestPeriodRecords[idx] = record
			}
		}
	}

	var largestCumulativePaymentBytes []byte
	largestCumulativePayment, err := s.meterer.MeteringStore.GetLargestCumulativePayment(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get largest cumulative payment, use zero value", "err", err, "accountID", accountID)

	} else {
		largestCumulativePaymentBytes = largestCumulativePayment.Bytes()
	}

	// on-Chain account state
	var pbReservation *pb.Reservation
	reservations, err := s.meterer.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, accountID, quorumIds)
	if err != nil {
		s.logger.Debug("failed to get onchain reservation, use zero values", "err", err, "accountID", accountID)
	} else {
		quorumNumbers := make([]uint32, len(reservations))
		for quorumNumber := range reservations {
			quorumNumbers[quorumNumber] = uint32(quorumNumber)
		}
		quorumSplits := make([]uint32, len(reservations))
		for quorumNumber := range reservations {
			quorumSplits[quorumNumber] = 0
		}

		// TODO: in a subsequent PR, we update PaymentState API types to include multiple quorum reservations;
		// For this PR, we return the first reservation as they are actually the same reservation
		pbReservation = &pb.Reservation{
			SymbolsPerSecond: reservations[0].SymbolsPerSecond,
			StartTimestamp:   uint32(reservations[0].StartTimestamp),
			EndTimestamp:     uint32(reservations[0].EndTimestamp),
			QuorumSplits:     quorumSplits,
			QuorumNumbers:    quorumNumbers,
		}
	}

	var onchainCumulativePaymentBytes []byte
	onDemandPayment, err := s.meterer.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get ondemand payment, use zero value", "err", err, "accountID", accountID)
	} else {
		onchainCumulativePaymentBytes = onDemandPayment.CumulativePayment.Bytes()
	}

	paymentGlobalParams := pb.PaymentGlobalParams{
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
		ReservationWindow:      reservationWindow,
	}

	// build reply
	reply := &pb.GetPaymentStateReply{
		PaymentGlobalParams:      &paymentGlobalParams,
		PeriodRecords:            highestPeriodRecords,
		Reservation:              pbReservation,
		CumulativePayment:        largestCumulativePaymentBytes,
		OnchainCumulativePayment: onchainCumulativePaymentBytes,
	}
	return reply, nil
}

// TODO(hopeyen): separate this into a subsequent PR
func (s *DispersalServerV2) GetPaymentStateForAllQuorums(ctx context.Context, req *pb.GetPaymentStateForAllQuorumsRequest) (*pb.GetPaymentStateForAllQuorumsReply, error) {
	return nil, api.NewErrorUnimplemented()
}

// getAllQuorumIds returns a slice of all quorum IDs (from 0 to quorumCount-1)
// Returns an empty slice if the onchain state is not loaded
func (o *OnchainState) getAllQuorumIds() []core.QuorumID {
	quorumCount := o.QuorumCount
	quorumIds := make([]core.QuorumID, quorumCount)
	for i := range quorumIds {
		quorumIds[i] = core.QuorumID(i)
	}

	return quorumIds
}
