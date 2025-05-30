package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

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
	globalSymbolsPerSecond := s.meterer.ChainPaymentState.GetGlobalSymbolsPerSecond()
	minNumSymbols := s.meterer.ChainPaymentState.GetMinNumSymbols()
	pricePerSymbol := s.meterer.ChainPaymentState.GetPricePerSymbol()
	reservationWindow := s.meterer.ChainPaymentState.GetReservationWindow()

	// on-Chain account state
	var pbReservation *pb.Reservation
	reservations, err := s.meterer.ChainPaymentState.GetReservedPaymentByAccount(ctx, accountID)
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

		// for all reservations, find the lowest SymbolsPerSecond
		lowestSymbolsPerSecond := uint64(math.MaxUint64)
		latestStartTimestamp := uint64(math.MaxUint64)
		earliestEndTimestamp := uint64(0)
		for _, reservation := range reservations {
			if reservation.SymbolsPerSecond < lowestSymbolsPerSecond {
				lowestSymbolsPerSecond = reservation.SymbolsPerSecond
			}
			if reservation.StartTimestamp > latestStartTimestamp {
				latestStartTimestamp = reservation.StartTimestamp
			}
			if reservation.EndTimestamp < earliestEndTimestamp {
				earliestEndTimestamp = reservation.EndTimestamp
			}
		}

		// TODO: in a subsequent PR, we update PaymentState API types to include multiple quorum reservations;
		// For this PR, we return the first reservation as they are actually the same reservation
		pbReservation = &pb.Reservation{
			SymbolsPerSecond: lowestSymbolsPerSecond,
			StartTimestamp:   uint32(latestStartTimestamp),
			EndTimestamp:     uint32(earliestEndTimestamp),
			QuorumSplits:     quorumSplits,
			QuorumNumbers:    quorumNumbers,
		}
	}

	// off-chain account specific payment state
	now := time.Now().Unix()
	currentReservationPeriod := meterer.GetReservationPeriod(now, reservationWindow)
	// take the first quorum number from the reservations as all the records are the same for current quorum-agnostic PaymentState
	var quorumNumber uint8
	if len(pbReservation.QuorumNumbers) > 0 {
		quorumNumber = uint8(pbReservation.QuorumNumbers[0])
	} else {
		quorumNumber = 0
	}
	periodRecords, err := s.meterer.MeteringStore.GetPeriodRecords(ctx, accountID, currentReservationPeriod, quorumNumber)
	if err != nil {
		s.logger.Debug("failed to get reservation records, use placeholders", "err", err, "accountID", accountID)
	}
	var largestCumulativePaymentBytes []byte
	largestCumulativePayment, err := s.meterer.MeteringStore.GetLargestCumulativePayment(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get largest cumulative payment, use zero value", "err", err, "accountID", accountID)

	} else {
		largestCumulativePaymentBytes = largestCumulativePayment.Bytes()
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
		PeriodRecords:            periodRecords[:],
		Reservation:              pbReservation,
		CumulativePayment:        largestCumulativePaymentBytes,
		OnchainCumulativePayment: onchainCumulativePaymentBytes,
	}
	return reply, nil
}

// getAllQuorumIds returns a slice of all quorum IDs (from 0 to quorumCount-1)
// Returns an empty slice if the onchain state is not loaded
func (s *DispersalServerV2) getAllQuorumIds() []uint8 {
	state := s.onchainState.Load()
	if state == nil {
		s.logger.Debug("onchain state not loaded yet")
		return []uint8{}
	}

	quorumCount := state.QuorumCount
	quorumIds := make([]uint8, quorumCount)
	for i := range quorumIds {
		quorumIds[i] = uint8(i)
	}

	return quorumIds
}

func (s *DispersalServerV2) GetPaymentStateQuorumSpecific(ctx context.Context, req *pb.GetPaymentStateQuorumSpecificRequest) (*pb.GetPaymentStateQuorumSpecificReply, error) {
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
	if err := s.blobRequestAuthenticator.AuthenticatePaymentStateQuorumSpecificRequest(accountID, req); err != nil {
		s.logger.Debug("failed to validate signature", "err", err, "accountID", accountID)
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
	}

	// Get on-Chain account state for reservations
	// TODO: update this to be pulling specific quorum IDs from the chain, after payment vault interface updates
	pbReservation := make(map[uint32]*pb.QuorumReservation)
	reservations, err := s.meterer.ChainPaymentState.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get onchain reservation, use zero values", "err", err, "accountID", accountID)
	} else {
		for quorumId, reservation := range reservations {
			pbReservation[uint32(quorumId)] = &pb.QuorumReservation{
				SymbolsPerSecond: reservation.SymbolsPerSecond,
				StartTimestamp:   uint32(reservation.StartTimestamp),
				EndTimestamp:     uint32(reservation.EndTimestamp),
			}
		}
	}

	// off-chain account specific payment state
	now := time.Now().Unix()
	reservationWindow := s.meterer.ChainPaymentState.GetReservationWindow()
	currentReservationPeriod := meterer.GetReservationPeriod(now, reservationWindow)

	// Get all quorum IDs from the system
	periodRecords := make(map[uint32]*pb.PeriodRecords)
	quorumIds := s.getAllQuorumIds()

	// Get all period records for this account across all quorums
	numQuorums := uint8(len(reservations))
	reservedQuorums := make([]uint8, numQuorums)
	for i := range reservedQuorums {
		reservedQuorums[i] = uint8(i)
	}

	//TODO(hopeyen): temporarily repeat the same record for all reserved quorums
	// update the MeteringStore interface in a subsequent PR for quorum specific period records
	quorumNumber := uint8(0)
	if len(reservedQuorums) > 0 {
		quorumNumber = reservedQuorums[0]
	}
	records, err := s.meterer.MeteringStore.GetPeriodRecords(ctx, accountID, currentReservationPeriod, quorumNumber)
	s.logger.Debug("offchain stored period records", "records", records)
	if err != nil {
		s.logger.Debug("failed to get reservation records for multiple quorums",
			"err", err, "accountID", accountID)
		return nil, err
	}
	pbPeriodRecords := &pb.PeriodRecords{
		Records: records,
	}
	for quorumId := range reservedQuorums {
		periodRecords[uint32(quorumId)] = pbPeriodRecords
	}

	// Get largest cumulative payment
	var largestCumulativePaymentBytes []byte
	largestCumulativePayment, err := s.meterer.MeteringStore.GetLargestCumulativePayment(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get largest cumulative payment, use zero value", "err", err, "accountID", accountID)
	} else {
		largestCumulativePaymentBytes = largestCumulativePayment.Bytes()
	}

	// Get on-demand payment information
	var onchainCumulativePaymentBytes []byte
	onDemandPayment, err := s.meterer.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get ondemand payment, use zero value", "err", err, "accountID", accountID)
	} else {
		onchainCumulativePaymentBytes = onDemandPayment.CumulativePayment.Bytes()
	}

	// Prepare payment vault parameters for the response
	// Get all quorum configurations
	quorumPaymentConfigs := make(map[uint32]*pb.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[uint32]*pb.PaymentQuorumProtocolConfig)

	// Build configurations for all available quorums
	for _, quorumNumber := range quorumIds {
		quorumID := uint32(quorumNumber)

		// Payment configuration for this quorum
		quorumPaymentConfigs[quorumID] = &pb.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: 0, // TODO: get from chain when available
			OnDemandSymbolsPerSecond:    s.meterer.ChainPaymentState.GetGlobalSymbolsPerSecond(),
			OnDemandPricePerSymbol:      s.meterer.ChainPaymentState.GetPricePerSymbol(),
		}

		// Protocol configuration for this quorum
		quorumProtocolConfigs[quorumID] = &pb.PaymentQuorumProtocolConfig{
			MinNumSymbols:              s.meterer.ChainPaymentState.GetMinNumSymbols(),
			ReservationAdvanceWindow:   0, // TODO: get from chain when available
			ReservationRateLimitWindow: s.meterer.ChainPaymentState.GetReservationWindow(),
			OnDemandRateLimitWindow:    s.meterer.ChainPaymentState.GetGlobalRatePeriodInterval(),
			OnDemandEnabled:            true, // TODO: get from chain when available
		}
	}

	// Get on-demand quorum numbers
	onDemandQuorumNumbers, err := s.meterer.ChainPaymentState.GetOnDemandQuorumNumbers(ctx)
	if err != nil {
		s.logger.Debug("failed to get on-demand quorum numbers, using default", "err", err)
		onDemandQuorumNumbers = []uint8{0} // fallback to quorum 0
	}
	onDemandQuorumNumbers32 := make([]uint32, len(onDemandQuorumNumbers))
	for i, qn := range onDemandQuorumNumbers {
		onDemandQuorumNumbers32[i] = uint32(qn)
	}

	paymentVaultParams := &pb.PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers32,
	}

	// Build reply
	reply := &pb.GetPaymentStateQuorumSpecificReply{
		PaymentVaultParams:       paymentVaultParams,
		PeriodRecords:            periodRecords,
		Reservations:             pbReservation,
		CumulativePayment:        largestCumulativePaymentBytes,
		OnchainCumulativePayment: onchainCumulativePaymentBytes,
	}
	return reply, nil
}
