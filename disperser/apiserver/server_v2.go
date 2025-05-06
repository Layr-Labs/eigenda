package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/eigenda/api"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbv1 "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	pbvalidator "github.com/Layr-Labs/eigenda/api/grpc/validator"
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
const requiredNodeVersion = "0.8.6"
const requiredNodeVersionStakeThreshold = 0.8 // 80%

type OnchainState struct {
	QuorumCount           uint8
	RequiredQuorums       []core.QuorumID
	BlobVersionParameters *corev2.BlobVersionParameterMap
	TTL                   time.Duration
	OperatorState         map[core.OperatorID]*OperatorOnchainState
}

// OperatorOnchainState holds operator info for apiserver
type OperatorOnchainState struct {
	OperatorID core.OperatorID
	Address    gethcommon.Address
	Socket     string
	Stake      *big.Int
	NodeInfo   *pbvalidator.GetNodeInfoReply // NodeInfo from api/grpc/validator/node_v2.pb.go
	Quorums    []core.QuorumID               // Quorums this operator is a member of
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

	ntpClock                        *core.NTPSyncedClock
	operatorSetRolloutReadyByQuorum map[core.QuorumID]bool
	currentRolloutStakePctByQuorum  map[core.QuorumID]float64 // for diagnostics and error messages
	operatorVersionCheck            bool                      // If true, enforce node version rollout check
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
	operatorVersionCheck ...bool, // variadic for backward compatibility
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

	check := true
	if len(operatorVersionCheck) > 0 {
		check = operatorVersionCheck[0]
	}

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

		metricsConfig: metricsConfig,
		metrics:       newAPIServerV2Metrics(registry, metricsConfig, logger),

		ntpClock:             ntpClock,
		operatorVersionCheck: check,
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

	// Periodic operator node info check (every hour)
	if s.operatorVersionCheck {
		go func() {
			ticker := time.NewTicker(time.Hour)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.periodicOperatorNodeInfoCheck(ctx)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

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

	var operatorState map[core.OperatorID]*OperatorOnchainState
	if s.operatorVersionCheck {
		operatorState, err = s.buildOperatorState(ctx, quorumCount, uint64(currentBlock))
		if err != nil {
			return err
		}
	} else {
		operatorState = make(map[core.OperatorID]*OperatorOnchainState)
	}

	onchainState := &OnchainState{
		QuorumCount:           quorumCount,
		RequiredQuorums:       requiredQuorums,
		BlobVersionParameters: v2.NewBlobVersionParameterMap(blobParams),
		TTL:                   time.Duration((storeDurationBlocks+blockStaleMeasure)*12) * time.Second,
		OperatorState:         operatorState,
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

// periodicOperatorNodeInfoCheck is a placeholder for the logic to fetch all operator endpoints and ping them for node info.
func (s *DispersalServerV2) periodicOperatorNodeInfoCheck(ctx context.Context) {
	s.logger.Info("Periodic operator node info check")
	quorumCount := s.onchainState.Load().QuorumCount
	currentBlock, err := s.chainReader.GetCurrentBlockNumber(ctx)
	if err != nil {
		s.logger.Error("failed to get current block number", "err", err)
		return
	}
	quorumIds := make([]core.QuorumID, quorumCount)
	for i := 0; i < int(quorumCount); i++ {
		quorumIds[i] = core.QuorumID(i)
	}
	operatorState, err := s.chainReader.GetOperatorVerboseState(ctx, quorumIds, currentBlock)
	if err != nil {
		s.logger.Error("failed to get operator info for quorums", "err", err)
		return
	}

	pctByQuorum, rolloutReady := CalculateQuorumRolloutReadiness(
		operatorState,
		requiredNodeVersion,
		requiredNodeVersionStakeThreshold,
	)
	for quorum, pct := range pctByQuorum {
		s.logger.Info("Operator version rollout check", "quorum", quorum, "required_version", requiredNodeVersion, "stake_pct", pct, "threshold", requiredNodeVersionStakeThreshold, "rollout_ready", rolloutReady[quorum])
	}
	s.currentRolloutStakePctByQuorum = pctByQuorum
	s.operatorSetRolloutReadyByQuorum = rolloutReady
}

// CalculateQuorumRolloutReadiness computes the stake percentage and readiness for each quorum.
func CalculateQuorumRolloutReadiness(
	ops core.OperatorStateVerbose,
	requiredVersion string,
	threshold float64,
) (map[core.QuorumID]float64, map[core.QuorumID]bool) {
	totalStakeByQuorum := make(map[core.QuorumID]*big.Int)
	upgradedStakeByQuorum := make(map[core.QuorumID]*big.Int)

	for quorumID, opMap := range ops {
		totalStake := big.NewInt(0)
		upgradedStake := big.NewInt(0)
		for _, opState := range opMap {
			if opState.Stake == nil {
				continue
			}
			totalStake.Add(totalStake, opState.Stake)
			if opState.NodeInfo != nil && opState.NodeInfo.Semver == requiredVersion {
				upgradedStake.Add(upgradedStake, opState.Stake)
			}
		}
		totalStakeByQuorum[quorumID] = totalStake
		upgradedStakeByQuorum[quorumID] = upgradedStake
	}

	pctByQuorum := make(map[core.QuorumID]float64)
	readyByQuorum := make(map[core.QuorumID]bool)
	for quorum, total := range totalStakeByQuorum {
		upgraded := upgradedStakeByQuorum[quorum]
		pct := 0.0
		if total.Cmp(big.NewInt(0)) > 0 && upgraded != nil {
			pct, _ = new(big.Rat).SetFrac(upgraded, total).Float64()
		}
		pctByQuorum[quorum] = pct
		readyByQuorum[quorum] = pct >= threshold
	}
	return pctByQuorum, readyByQuorum
}

// buildOperatorState build the operator state map
func (s *DispersalServerV2) buildOperatorState(ctx context.Context, quorumCount uint8, currentBlock uint64) (map[core.OperatorID]*OperatorOnchainState, error) {
	operatorState := make(map[core.OperatorID]*OperatorOnchainState)
	quorums := make([]core.QuorumID, quorumCount)
	for i := 0; i < int(quorumCount); i++ {
		quorums[i] = core.QuorumID(i)
	}
	stakesWithSocket, err := s.chainReader.GetOperatorStakesWithSocketForQuorums(ctx, quorums, uint32(currentBlock))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator stakes with socket: %w", err)
	}
	seen := make(map[core.OperatorID]struct{})

	// Build a map from OperatorID to set of quorums
	operatorQuorums := make(map[core.OperatorID]map[core.QuorumID]struct{})
	for quorum, ops := range stakesWithSocket {
		for _, op := range ops {
			if operatorQuorums[op.OperatorID] == nil {
				operatorQuorums[op.OperatorID] = make(map[core.QuorumID]struct{})
			}
			operatorQuorums[op.OperatorID][quorum] = struct{}{}
		}
	}
	for _, ops := range stakesWithSocket {
		for _, op := range ops {
			if _, ok := seen[op.OperatorID]; ok {
				continue
			}
			seen[op.OperatorID] = struct{}{}
			address, err := s.chainReader.OperatorIDToAddress(ctx, op.OperatorID)
			if err != nil {
				s.logger.Warn("Failed to get operator address", "operator", op.OperatorID.Hex(), "err", err)
				continue
			}
			// Collect all quorums for this operator
			quorums := make([]core.QuorumID, 0, len(operatorQuorums[op.OperatorID]))
			for q := range operatorQuorums[op.OperatorID] {
				quorums = append(quorums, q)
			}
			operatorState[op.OperatorID] = &OperatorOnchainState{
				OperatorID: op.OperatorID,
				Address:    address,
				Socket:     string(op.Socket),
				Stake:      op.Stake,
				NodeInfo:   nil, // To be filled by periodic pings
				Quorums:    quorums,
			}
		}
	}
	return operatorState, nil
}
