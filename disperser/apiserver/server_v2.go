package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
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

// Hardcoded required node version and threshold
const requiredNodeVersion = ">=0.9.0-rc.1"
const requiredNodeVersionStakeThreshold = 0.8 // 80%

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

	operatorSetRolloutReadyByQuorum map[core.QuorumID]bool
	currentRolloutStakePctByQuorum  map[core.QuorumID]float64 // for diagnostics and error messages
	operatorVersionCheck            bool                      // If true, enforce node version rollout check
	nodeInfoCheckInitOnce           sync.Once                 // Ensures we only initialize the node info check once
}

// NewDispersalServerV2 creates a new Server struct with the provided parameters.
func NewDispersalServerV2(
	serverConfig disperser.ServerConfig,
	blobStore *blobstore.BlobStore,
	blobMetadataStore *blobstore.BlobMetadataStore,
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
	operatorVersionCheck bool,
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
		serverConfig:                serverConfig,
		blobStore:                   blobStore,
		blobMetadataStore:           blobMetadataStore,
		chainReader:                 chainReader,
		blobRequestAuthenticator:    blobRequestAuthenticator,
		meterer:                     meterer,
		prover:                      prover,
		logger:                      logger,
		maxNumSymbolsPerBlob:        maxNumSymbolsPerBlob,
		onchainStateRefreshInterval: onchainStateRefreshInterval,
		operatorVersionCheck:        operatorVersionCheck,

		metricsConfig: metricsConfig,
		metrics:       newAPIServerV2Metrics(registry, metricsConfig, logger),

		ntpClock: ntpClock,
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
		s.logger.Info("Onchain state refresh ticker started")

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
	s.logger.Debug("RefreshOnchainState: Starting refresh")

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

	// off-chain account specific payment state
	now := time.Now().Unix()
	currentReservationPeriod := meterer.GetReservationPeriod(now, reservationWindow)
	periodRecords, err := s.meterer.OffchainStore.GetPeriodRecords(ctx, accountID, currentReservationPeriod)
	if err != nil {
		s.logger.Debug("failed to get reservation records, use placeholders", "err", err, "accountID", accountID)
	}
	var largestCumulativePaymentBytes []byte
	largestCumulativePayment, err := s.meterer.OffchainStore.GetLargestCumulativePayment(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get largest cumulative payment, use zero value", "err", err, "accountID", accountID)

	} else {
		largestCumulativePaymentBytes = largestCumulativePayment.Bytes()
	}
	// on-Chain account state
	var pbReservation *pb.Reservation
	reservation, err := s.meterer.ChainPaymentState.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		s.logger.Debug("failed to get onchain reservation, use zero values", "err", err, "accountID", accountID)
	} else {
		quorumNumbers := make([]uint32, len(reservation.QuorumNumbers))
		for i, v := range reservation.QuorumNumbers {
			quorumNumbers[i] = uint32(v)
		}
		quorumSplits := make([]uint32, len(reservation.QuorumSplits))
		for i, v := range reservation.QuorumSplits {
			quorumSplits[i] = uint32(v)
		}

		pbReservation = &pb.Reservation{
			SymbolsPerSecond: reservation.SymbolsPerSecond,
			StartTimestamp:   uint32(reservation.StartTimestamp),
			EndTimestamp:     uint32(reservation.EndTimestamp),
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
		PeriodRecords:            periodRecords[:],
		Reservation:              pbReservation,
		CumulativePayment:        largestCumulativePaymentBytes,
		OnchainCumulativePayment: onchainCumulativePaymentBytes,
	}
	return reply, nil
}

// periodicOperatorNodeInfoCheck performs a periodic check of the operator node info, including
// on-chain registration of operator ID, stake, socket, and offchain self-claimed node version.
func (s *DispersalServerV2) periodicOperatorNodeInfoCheck(ctx context.Context) {
	s.logger.Debug("Periodic operator node info check started")
	onchainState := s.onchainState.Load()
	if onchainState == nil {
		s.logger.Error("onchain state is nil during operator node info check")
		return
	}

	quorumCount := onchainState.QuorumCount

	currentBlock, err := s.chainReader.GetCurrentBlockNumber(ctx)
	if err != nil {
		s.logger.Error("failed to get current block number", "err", err)
		return
	}
	quorumIds := make([]core.QuorumID, quorumCount)
	for i := 0; i < int(quorumCount); i++ {
		quorumIds[i] = core.QuorumID(i)
	}

	stakesWithSocket, err := s.chainReader.GetOperatorStakesWithSocketForQuorums(ctx, quorumIds, uint32(currentBlock))
	if err != nil {
		s.logger.Error("failed to get operator stakes with socket", "err", err)
		return
	}

	operatorState, err := core.GetOperatorVerboseState(ctx, stakesWithSocket, quorumIds, currentBlock)
	if err != nil {
		s.logger.Error("failed to get operator info for quorums", "err", err)
		return
	}

	pctByQuorum, rolloutReady := core.CalculateQuorumRolloutReadiness(
		operatorState,
		requiredNodeVersion,
		requiredNodeVersionStakeThreshold,
	)

	// Log detailed results for debugging
	for quorum, pct := range pctByQuorum {
		ready := rolloutReady[quorum]
		s.logger.Debug("Operator version rollout check result",
			"quorum", quorum,
			"required_version", requiredNodeVersion,
			"stake_pct", pct,
			"threshold", requiredNodeVersionStakeThreshold,
			"rollout_ready", ready,
			"upgradedThreshold", pct >= requiredNodeVersionStakeThreshold)
	}

	// If no quorums were found, log a warning
	if len(pctByQuorum) == 0 {
		s.logger.Warn("No quorums found in rollout readiness calculation")
	}

	// Update the server state with the new data
	s.currentRolloutStakePctByQuorum = pctByQuorum
	s.operatorSetRolloutReadyByQuorum = rolloutReady
}

func (s *DispersalServerV2) checkQuorumRolloutReady(req *pb.DisperseBlobRequest) error {
	// Defensive check - initialize maps if they're nil
	if s.operatorSetRolloutReadyByQuorum == nil {
		s.operatorSetRolloutReadyByQuorum = make(map[core.QuorumID]bool)
	}
	if s.currentRolloutStakePctByQuorum == nil {
		s.currentRolloutStakePctByQuorum = make(map[core.QuorumID]float64)
	}

	notReady := make([]string, 0)
	for _, quorum := range req.GetBlobHeader().GetQuorumNumbers() {
		qid := core.QuorumID(quorum)
		ready, exists := s.operatorSetRolloutReadyByQuorum[qid]
		if !exists {
			s.logger.Warn("Quorum not found in rollout readiness map", "quorumID", qid)
			ready = false // Default to not ready if missing
		}

		pct, pctExists := s.currentRolloutStakePctByQuorum[qid]
		if !pctExists {
			s.logger.Warn("Quorum not found in stake percentage map", "quorumID", qid)
			pct = 0 // Default to 0% if missing
		}

		s.logger.Info("Checking quorum readiness", "quorumID", qid, "ready", ready, "stakePct", pct)
		if !ready {
			threshold := requiredNodeVersionStakeThreshold * 100
			notReady = append(notReady, fmt.Sprintf("quorum %d: %.2f%% of %.2f%% upgraded", quorum, pct*100, threshold))
		}
	}
	if len(notReady) > 0 {
		errMsg := fmt.Sprintf(
			"Operator new version %s rollout is in progress for: %s. Please be patient for the decentralized network to coordinate disperser and operator versions.",
			requiredNodeVersion,
			strings.Join(notReady, "; "),
		)
		s.logger.Warn("Quorum rollout not ready", "error", errMsg)
		return api.NewErrorResourceExhausted(errMsg)
	}
	s.logger.Info("All quorums are rollout ready")
	return nil
}
