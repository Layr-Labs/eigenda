package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	pbv1 "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/signingrate"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	chainId           *big.Int
	blobStore         *blobstore.BlobStore
	blobMetadataStore blobstore.MetadataStore
	meterer           *meterer.Meterer

	chainReader              core.Reader
	blobRequestAuthenticator corev2.BlobRequestAuthenticator
	committer                *committer.Committer
	logger                   logging.Logger

	// state
	onchainState                atomic.Pointer[OnchainState]
	maxNumSymbolsPerBlob        uint32
	onchainStateRefreshInterval time.Duration

	// MaxDispersalAge is the maximum age a dispersal request can be before it is rejected.
	// Dispersals older than this duration are rejected by the API server at ingest.
	//
	// Age is determined by the BlobHeader.PaymentMetadata.Timestamp field, which is set by the
	// client at dispersal request creation time (in nanoseconds since Unix epoch).
	MaxDispersalAge time.Duration

	// MaxFutureDispersalTime is the maximum amount of time into the future a dispersal request can be
	// before it is rejected. Dispersals with timestamps more than this duration in the future are rejected
	// by the API server at ingest.
	MaxFutureDispersalTime time.Duration

	// getNow returns the current time
	getNow func() time.Time

	metricsConfig disperser.MetricsConfig
	metrics       *metricsV2

	// ReservedOnly mode doesn't support on-demand payments
	// This would be removed with decentralized ratelimiting
	ReservedOnly bool

	// Exists as a member variable so that the connection can be closed inside Stop().
	controllerConnection *grpc.ClientConn

	// Client for making gRPC calls to the controller for payment authorization.
	controllerClient controller.ControllerServiceClient

	// Pre-created listener for the gRPC server
	listener   net.Listener
	grpcServer *grpc.Server

	// DisableGetBlobCommitment, if true, causes the GetBlobCommitment gRPC endpoint to return
	// a deprecation error. This endpoint is deprecated and will be removed in a future release.
	disableGetBlobCommitment bool

	// Tracks signing rates for validators. This data is mirrored from the controller's signing rate tracker,
	// so that external requests can be serviced without involving the controller.
	signingRateTracker signingrate.SigningRateTracker
}

// NewDispersalServerV2 creates a new Server struct with the provided parameters.
func NewDispersalServerV2(
	serverConfig disperser.ServerConfig,
	getNow func() time.Time,
	chainId *big.Int,
	blobStore *blobstore.BlobStore,
	blobMetadataStore blobstore.MetadataStore,
	chainReader core.Reader,
	meterer *meterer.Meterer,
	blobRequestAuthenticator corev2.BlobRequestAuthenticator,
	committer *committer.Committer,
	maxNumSymbolsPerBlob uint32,
	onchainStateRefreshInterval time.Duration,
	maxDispersalAge time.Duration,
	maxFutureDispersalTime time.Duration,
	_logger logging.Logger,
	registry *prometheus.Registry,
	metricsConfig disperser.MetricsConfig,
	ReservedOnly bool,
	controllerConnection *grpc.ClientConn,
	controllerClient controller.ControllerServiceClient,
	listener net.Listener,
	signingRateTracker signingrate.SigningRateTracker,
) (*DispersalServerV2, error) {
	if listener == nil {
		return nil, errors.New("listener is required")
	}
	if serverConfig.GrpcPort == "" {
		return nil, errors.New("grpc port is required")
	}
	if getNow == nil {
		return nil, errors.New("getNow is required")
	}
	if chainId == nil {
		return nil, errors.New("chainId is required")
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
	if committer == nil {
		return nil, errors.New("committer is required")
	}
	if signingRateTracker == nil {
		return nil, errors.New("signingRateTracker is required")
	}
	if maxNumSymbolsPerBlob == 0 {
		return nil, errors.New("maxNumSymbolsPerBlob is required")
	}
	if _logger == nil {
		return nil, errors.New("logger is required")
	}
	if maxDispersalAge <= 0 {
		return nil, fmt.Errorf("maxDispersalAge must be positive (got: %v)", maxDispersalAge)
	}
	if maxFutureDispersalTime <= 0 {
		return nil, fmt.Errorf("maxFutureDispersalTime must be positive (got: %v)", maxFutureDispersalTime)
	}

	logger := _logger.With("component", "DispersalServerV2")

	if controllerClient == nil {
		return nil, errors.New("controller client is required")
	}

	return &DispersalServerV2{
		serverConfig:      serverConfig,
		chainId:           chainId,
		blobStore:         blobStore,
		blobMetadataStore: blobMetadataStore,

		chainReader:              chainReader,
		blobRequestAuthenticator: blobRequestAuthenticator,
		meterer:                  meterer,
		committer:                committer,
		logger:                   logger,

		maxNumSymbolsPerBlob:        maxNumSymbolsPerBlob,
		onchainStateRefreshInterval: onchainStateRefreshInterval,
		MaxDispersalAge:             maxDispersalAge,
		MaxFutureDispersalTime:      maxFutureDispersalTime,
		getNow:                      getNow,

		metricsConfig: metricsConfig,
		metrics:       newAPIServerV2Metrics(registry, metricsConfig, logger),

		ReservedOnly:             ReservedOnly,
		controllerConnection:     controllerConnection,
		controllerClient:         controllerClient,
		listener:                 listener,
		disableGetBlobCommitment: serverConfig.DisableGetBlobCommitment,
		signingRateTracker:       signingRateTracker,
	}, nil
}

func (s *DispersalServerV2) Start(ctx context.Context) error {
	// Start the metrics server
	if s.metricsConfig.EnableMetrics {
		s.metrics.Start(context.Background())
		// Set configuration gauges
		s.metrics.setDispersalTimestampConfig(
			s.MaxDispersalAge.Seconds(),
			s.MaxFutureDispersalTime.Seconds(),
		)
	}

	// Serve grpc requests
	keepAliveConfig := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     s.serverConfig.MaxIdleConnectionAge,
		MaxConnectionAge:      s.serverConfig.MaxConnectionAge,
		MaxConnectionAgeGrace: s.serverConfig.MaxConnectionAgeGrace,
	})

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.metrics.grpcMetrics.UnaryServerInterceptor(),
		), opt, keepAliveConfig)
	reflection.Register(s.grpcServer)
	pb.RegisterDisperserServer(s.grpcServer, s)

	// Unimplemented v1 server for grpcurl/reflection support
	pbv1.RegisterDisperserServer(s.grpcServer, &DispersalServerV1{})

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, s.grpcServer)

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

	s.logger.Info("GRPC Listening", "port", s.serverConfig.GrpcPort, "address", s.listener.Addr().String())

	if err := s.grpcServer.Serve(s.listener); err != nil {
		return fmt.Errorf("could not start GRPC server: %w", err)
	}

	return nil
}

func (s *DispersalServerV2) GetBlobCommitment(
	ctx context.Context,
	req *pb.BlobCommitmentRequest,
) (*pb.BlobCommitmentReply, error) {
	reply, st := s.getBlobCommitment(req)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *DispersalServerV2) getBlobCommitment(
	req *pb.BlobCommitmentRequest,
) (*pb.BlobCommitmentReply, *status.Status) {
	start := time.Now()
	defer func() {
		s.metrics.reportGetBlobCommitmentLatency(time.Since(start))
	}()

	if s.disableGetBlobCommitment {
		return nil, status.New(codes.Unimplemented, "GetBlobCommitment is deprecated and has been disabled. This service will be removed in a future release. Please compute blob commitments locally.")
	}

	if s.committer == nil {
		return nil, status.New(codes.Internal, "committer is not configured")
	}
	blobSize := uint32(len(req.GetBlob()))
	if blobSize == 0 {
		return nil, status.New(codes.InvalidArgument, "blob cannot be empty")
	}
	if encoding.GetBlobLengthPowerOf2(blobSize) > s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL {
		return nil, status.Newf(codes.InvalidArgument, "blob size cannot exceed %v bytes",
			s.maxNumSymbolsPerBlob*encoding.BYTES_PER_SYMBOL)
	}
	c, err := s.committer.GetCommitmentsForPaddedLength(req.GetBlob())
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to compute commitments: %v", err)
	}
	commitment, err := c.Commitment.Serialize()
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to serialize commitment: %v", err)
	}
	lengthCommitment, err := c.LengthCommitment.Serialize()
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to serialize length commitment: %v", err)
	}
	lengthProof, err := c.LengthProof.Serialize()
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to serialize length proof: %v", err)
	}

	return &pb.BlobCommitmentReply{
		BlobCommitment: &pbcommon.BlobCommitment{
			Commitment:       commitment,
			LengthCommitment: lengthCommitment,
			LengthProof:      lengthProof,
			Length:           uint32(c.Length),
		}}, status.New(codes.OK, "")
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
		BlobVersionParameters: corev2.NewBlobVersionParameterMap(blobParams),
		TTL:                   time.Duration((storeDurationBlocks+blockStaleMeasure)*12) * time.Second,
	}

	s.onchainState.Store(onchainState)

	return nil
}

func (s *DispersalServerV2) GetPaymentState(
	ctx context.Context,
	req *pb.GetPaymentStateRequest,
) (*pb.GetPaymentStateReply, error) {
	reply, st := s.getPaymentState(ctx, req)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *DispersalServerV2) getPaymentState(
	ctx context.Context,
	req *pb.GetPaymentStateRequest,
) (*pb.GetPaymentStateReply, *status.Status) {
	if s.meterer == nil {
		return nil, status.New(codes.Internal, "meterer is not configured")
	}
	start := time.Now()
	defer func() {
		s.metrics.reportGetPaymentStateLatency(time.Since(start))
	}()

	if !gethcommon.IsHexAddress(req.GetAccountId()) {
		return nil, status.New(codes.InvalidArgument, "invalid account ID")
	}

	accountID := gethcommon.HexToAddress(req.GetAccountId())

	// validate the signature
	if err := s.blobRequestAuthenticator.AuthenticatePaymentStateRequest(accountID, req); err != nil {
		s.logger.Debug("failed to validate signature", "err", err, "accountID", accountID)
		return nil, status.Newf(codes.Unauthenticated, "failed to validate signature: %s", err.Error())
	}
	// on-chain global payment parameters
	globalSymbolsPerSecond := s.meterer.ChainPaymentState.GetGlobalSymbolsPerSecond()
	minNumSymbols := s.meterer.ChainPaymentState.GetMinNumSymbols()
	pricePerSymbol := s.meterer.ChainPaymentState.GetPricePerSymbol()
	reservationWindow := s.meterer.ChainPaymentState.GetReservationWindow()

	// off-chain account specific payment state
	now := time.Now().Unix()
	currentReservationPeriod := meterer.GetReservationPeriod(now, reservationWindow)
	periodRecords, err := s.meterer.MeteringStore.GetPeriodRecords(ctx, accountID, currentReservationPeriod)
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
	return reply, status.New(codes.OK, "")
}

func (s *DispersalServerV2) GetValidatorSigningRate(
	ctx context.Context,
	request *pb.GetValidatorSigningRateRequest,
) (*pb.GetValidatorSigningRateReply, error) {

	if len(request.GetValidatorId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "validator id must be non-empty")
	}

	validatorId := core.OperatorID(request.GetValidatorId())

	signingRate, err := s.signingRateTracker.GetValidatorSigningRate(
		core.QuorumID(request.GetQuorum()),
		validatorId,
		time.Unix(int64(request.GetStartTimestamp()), 0),
		time.Unix(int64(request.GetEndTimestamp()), 0))

	if err != nil {
		return nil, fmt.Errorf("failed to get signing rate for validator %s: %w", validatorId.Hex(), err)
	}

	return &pb.GetValidatorSigningRateReply{
		ValidatorSigningRate: signingRate,
	}, nil
}

// Gracefully shuts down the server and closes any open connections
func (s *DispersalServerV2) Stop() error {
	if s.grpcServer != nil {
		// GracefulStop will close the listener that was passed to Serve()
		s.grpcServer.GracefulStop()
	}

	if s.controllerConnection != nil {
		if err := s.controllerConnection.Close(); err != nil {
			return fmt.Errorf("failed to close controller connection: %w", err)
		}
	}
	return nil
}
