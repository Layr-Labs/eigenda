package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/common/version"
	"github.com/Layr-Labs/eigenda/core"
	coreauthv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigenda/node/grpc/middleware"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/time/rate"
)

// ServerV2 implements the Node v2 proto APIs.
type ServerV2 struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	config             *node.Config
	node               *node.Node
	ratelimiter        common.RateLimiter
	logger             logging.Logger
	metrics            *MetricsV2
	chunkAuthenticator auth.RequestAuthenticator
	blobAuthenticator  corev2.BlobRequestAuthenticator
	replayGuardian     replay.ReplayGuardian

	// The current software version.
	softwareVersion *version.Semver

	// Pre-created listeners for the gRPC servers
	dispersalListener net.Listener
	retrievalListener net.Listener

	rateLimiter *middleware.DisperserRateLimiter
}

// NewServerV2 creates a new Server instance with the provided parameters.
func NewServerV2(
	ctx context.Context,
	config *node.Config,
	node *node.Node,
	logger logging.Logger,
	ratelimiter common.RateLimiter,
	registry *prometheus.Registry,
	reader core.Reader,
	softwareVersion *version.Semver,
	dispersalListener net.Listener,
	retrievalListener net.Listener) (*ServerV2, error) {

	metrics, err := NewV2Metrics(logger, registry)
	if err != nil {
		return nil, err
	}

	chunkAuthenticator, err := auth.NewRequestAuthenticator(
		ctx,
		reader,
		logger,
		config.DispersalAuthenticationKeyCacheSize,
		config.DisperserKeyTimeout,
		// TODO(litt3): once the checkpointed onchain config registry is ready, the authorized
		// on-demand dispersers should be read from there instead of being hardcoded.
		[]uint32{0}, // Default to disperser ID 0 for on-demand payments
		time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}
	blobAuthenticator := coreauthv2.NewBlobRequestAuthenticator()
	replayGuardian, err := replay.NewReplayGuardian(
		time.Now,
		config.StoreChunksRequestMaxPastAge,
		config.StoreChunksRequestMaxFutureAge)
	if err != nil {
		return nil, fmt.Errorf("failed to create replay guardian: %w", err)
	}

	return &ServerV2{
		config:             config,
		node:               node,
		ratelimiter:        ratelimiter,
		logger:             logger,
		metrics:            metrics,
		chunkAuthenticator: chunkAuthenticator,
		blobAuthenticator:  blobAuthenticator,
		replayGuardian:     replayGuardian,
		softwareVersion:    softwareVersion,
		dispersalListener:  dispersalListener,
		retrievalListener:  retrievalListener,
		rateLimiter: middleware.NewDisperserRateLimiter(
			logger,
			config.DisperserRateLimitPerSecond,
			config.DisperserRateLimitBurst,
		),
	}, nil
}

// GetDispersalPort returns the port number the dispersal listener is bound to.
func (s *ServerV2) GetDispersalPort() int {
	if s.dispersalListener == nil {
		return 0
	}
	return s.dispersalListener.Addr().(*net.TCPAddr).Port
}

// GetRetrievalPort returns the port number the retrieval listener is bound to.
func (s *ServerV2) GetRetrievalPort() int {
	if s.retrievalListener == nil {
		return 0
	}
	return s.retrievalListener.Addr().(*net.TCPAddr).Port
}

// Stop shuts down the listeners
func (s *ServerV2) Stop() {
	s.logger.Info("ServerV2 stop requested")

	if s.dispersalListener != nil {
		if err := s.dispersalListener.Close(); err != nil {
			s.logger.Warn("Failed to close dispersal listener", "error", err)
		}
	}

	if s.retrievalListener != nil {
		if err := s.retrievalListener.Close(); err != nil {
			s.logger.Warn("Failed to close retrieval listener", "error", err)
		}
	}
}

func (s *ServerV2) GetNodeInfo(ctx context.Context, in *pb.GetNodeInfoRequest) (*pb.GetNodeInfoReply, error) {
	if s.config.DisableNodeInfoResources {
		return &pb.GetNodeInfoReply{Semver: s.softwareVersion.String()}, nil
	}

	memBytes := uint64(0)
	v, err := mem.VirtualMemory()
	if err == nil {
		memBytes = v.Total
	}

	return &pb.GetNodeInfoReply{
		Semver:   s.softwareVersion.String(),
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCpu:   uint32(runtime.GOMAXPROCS(0)),
		MemBytes: memBytes,
	}, nil
}

func (s *ServerV2) StoreChunks(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	if !s.config.EnableV2 {
		return nil, api.NewErrorInvalidArg("v2 API is disabled")
	}

	if s.node.BLSSigner == nil {
		return nil, api.NewErrorInternal("missing bls signer")
	}

	probe := s.metrics.GetStoreChunksProbe()
	defer probe.End()

	probe.SetStage("validate")

	onDemandReservations := make([]*rate.Reservation, 0)
	success := false
	defer func() {
		if !success {
			for _, reservation := range onDemandReservations {
				s.node.CancelOnDemandDispersal(reservation)
			}
		}
	}()

	// Validate the request parameters (which is cheap) before starting any further
	// processing of the request.
	batch, err := s.validateStoreChunksRequest(in)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate store chunk request: %v", err))
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to serialize batch header hash: %v", err))
	}

	now := time.Now()
	if authenticatedID, ok := middleware.AuthenticatedDisperserIDFromContext(ctx); ok {
		// Defensive check: the interceptor should only set an ID that matches the request.
		if authenticatedID != in.GetDisperserID() {
			//nolint:wrapcheck
			return nil, api.NewErrorInvalidArg("authenticated disperser ID does not match request disperser ID")
		}
	} else {
		// Defense-in-depth: normally the gRPC interceptor authenticates StoreChunks and rate limits dispersers.
		// This fallback exists for direct calls (e.g. tests) or alternate wiring where the interceptor isn't installed.
		_, err = s.chunkAuthenticator.AuthenticateStoreChunksRequest(ctx, in, now)
		if err != nil {
			//nolint:wrapcheck
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to authenticate request: %v", err))
		}

		if s.rateLimiter != nil && !s.rateLimiter.Allow(in.GetDisperserID(), now) {
			//nolint:wrapcheck
			return nil, api.NewErrorResourceExhausted(
				fmt.Sprintf("disperser %d is rate limited", in.GetDisperserID()))
		}
	}

	if !s.chunkAuthenticator.IsDisperserAuthorized(in.GetDisperserID(), batch) {
		//nolint:wrapcheck
		return nil, api.NewErrorPermissionDenied(
			fmt.Sprintf("disperser %d not authorized for on-demand payments", in.GetDisperserID()))
	}

	for _, blobCert := range batch.BlobCertificates {
		_, err = s.validateDispersalRequest(blobCert)
		if err != nil {
			//nolint:wrapcheck
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate blob request: %v", err))
		}
	}

	blobHeadersAndTimestamps, err := hashing.HashBlobHeadersAndTimestamps(in)
	if err != nil {
		//nolint:wrapcheck
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to hash blob headers and timestamps: %v", err))
	}

	for i, blobHeader := range blobHeadersAndTimestamps {
		err = s.replayGuardian.VerifyRequest(blobHeader.Hash, blobHeader.Timestamp)
		if err != nil {
			//nolint:wrapcheck
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify blob header hash at index %d: %v", i, err))
		}
	}

	for _, blobCert := range batch.BlobCertificates {
		if blobCert.BlobHeader.PaymentMetadata.IsOnDemand() {
			length := blobCert.BlobHeader.BlobCommitments.Length
			reservation, meterErr := s.node.MeterOnDemandDispersal(length)
			if meterErr != nil {
				return nil, fmt.Errorf("global on-demand rate limit exceeded: %w", meterErr)
			}
			onDemandReservations = append(onDemandReservations, reservation)
		}
	}

	// Validate reservation payments (on-demand payments are validated on the the disperser's controller service)
	//
	// Note: the payment processing that occurs within this method is NOT reverted, even if something fails further
	// along. There are a couple reasons for this:
	// 1. At this stage, the dispersal request has already been sent to other validators. Even if this individual
	// validator were to revert the payment after some type of failure, there's no way to make sure that all other
	// validators would experience the same failure and revert. It is important to keep validator payment state in
	// sync, so the safest behavior is to just treat this as the point-of-no-return, from a payments perspective.
	// 2. Even if there were a way for all validators to agree on what payments to revert, non-trivial amounts of work
	// are being done shortly after this payment validation completes, for which the validators should be compensated.
	//
	// This accounting logic relies on each dispersal only arriving at this stage *once*. That is currently guaranteed
	// based on the replay guardian above. If the replay guardian were ever to be removed (for example, to enable
	// retried dispersals) then the accounting logic here would need to be revisited, and made retry tolerant.
	err = s.node.ValidateReservationPayment(ctx, batch, probe)
	if err != nil {
		return nil, fmt.Errorf("validate reservation payment: %w", err)
	}

	probe.SetStage("get_operator_state")
	s.logger.Info("new StoreChunks request",
		"batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]),
		"numBlobs", len(batch.BlobCertificates),
		"referenceBlockNumber", batch.BatchHeader.ReferenceBlockNumber)

	quorums := make(map[core.QuorumID]struct{}, len(batch.BlobCertificates))
	for _, blobCert := range batch.BlobCertificates {
		for _, quorum := range blobCert.BlobHeader.QuorumNumbers {
			quorums[quorum] = struct{}{}
		}
	}

	quorumList := make([]core.QuorumID, 0, len(quorums))
	for quorum := range quorums {
		quorumList = append(quorumList, quorum)
	}

	operatorState, err := s.node.OperatorStateCache.GetOperatorState(
		ctx, batch.BatchHeader.ReferenceBlockNumber, quorumList)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get the operator state: %v", err))
	}

	downloadSizeInBytes, relayRequests, err :=
		s.node.DetermineChunkLocations(batch, operatorState, probe)
	if err != nil {
		//nolint:wrapcheck
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to determine chunk locations: %v", err))
	}

	// storeChunksSemaphore can be nil during unit tests, since there are a bunch of places where the Node struct
	// is instantiated directly without using the constructor.
	if s.node.StoreChunksSemaphore != nil {
		// So far, we've only downloaded metadata for the blob. Before downloading the actual chunks, make sure there
		// is capacity in the store chunks buffer. This is an OOM safety measure.

		probe.SetStage("acquire_buffer_capacity")
		semaphoreCtx, cancel := context.WithTimeout(ctx, s.node.Config.StoreChunksBufferTimeout)
		defer cancel()
		err = s.node.StoreChunksSemaphore.Acquire(semaphoreCtx, int64(downloadSizeInBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to acquire buffer capacity: %w", err)
		}
		defer s.node.StoreChunksSemaphore.Release(int64(downloadSizeInBytes))
	}

	blobShards, rawBundles, err := s.node.DownloadChunksFromRelays(ctx, batch, relayRequests, probe)
	if err != nil {
		//nolint:wrapcheck
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to download chunks: %v", err))
	}

	err = s.validateAndStoreChunks(ctx, batch, blobShards, rawBundles, operatorState, batchHeaderHash, probe)
	if err != nil {
		return nil, err
	}

	probe.SetStage("sign")
	sig, err := s.node.BLSSigner.Sign(ctx, batchHeaderHash[:])
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to sign batch: %v", err))
	}

	success = true

	return &pb.StoreChunksReply{
		Signature: sig,
	}, nil
}

func (s *ServerV2) validateAndStoreChunks(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	rawBundles []*node.RawBundle,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte,
	probe *common.SequenceProbe,
) error {

	batchData := make([]*node.BundleToStore, 0, len(rawBundles))
	for _, bundle := range rawBundles {
		blobKey, err := bundle.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return api.NewErrorInternal("failed to get blob key")
		}

		// The current sampling scheme will store the same chunks for all quorums, so we always use quorum 0 as the quorum key in storage.
		quorum := core.QuorumID(0)

		bundleKey, err := node.BundleKey(blobKey, quorum)
		if err != nil {
			return api.NewErrorInternal("failed to get bundle key")
		}

		batchData = append(batchData, &node.BundleToStore{
			BundleKey:   bundleKey,
			BundleBytes: bundle.Bundle,
		})
	}

	return s.validateAndStoreChunksLittDB(
		ctx,
		batch,
		blobShards,
		batchData,
		operatorState,
		batchHeaderHash,
		probe)
}

func (s *ServerV2) validateAndStoreChunksLittDB(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	batchData []*node.BundleToStore,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte,
	probe *common.SequenceProbe,
) error {
	probe.SetStage("validate")
	err := s.node.ValidateBatchV2(ctx, batch, blobShards, operatorState)
	if err != nil {
		return api.NewErrorInternal(
			fmt.Sprintf("failed to validate batch %s: %v", hex.EncodeToString(batchHeaderHash[:]), err))
	}

	probe.SetStage("store")
	size, err := s.node.ValidatorStore.StoreBatch(batchData)
	if err != nil {
		return api.NewErrorInternal(
			fmt.Sprintf("failed to store batch %s: %v", hex.EncodeToString(batchHeaderHash[:]), err))
	}

	s.metrics.ReportStoreChunksRequestSize(size)

	return nil
}

// validateStoreChunksRequest validates the StoreChunksRequest and returns deserialized batch in the request
func (s *ServerV2) validateStoreChunksRequest(req *pb.StoreChunksRequest) (*corev2.Batch, error) {
	// The signature is created by go-ethereum library, which contains 1 additional byte (for
	// recovering the public key from signature), so it's 65 bytes.
	if len(req.GetSignature()) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes, found %d bytes", len(req.GetSignature()))
	}

	if req.GetBatch() == nil {
		return nil, errors.New("missing batch in request")
	}

	// BatchFromProtobuf internally validates the Batch while deserializing
	batch, err := corev2.BatchFromProtobuf(req.GetBatch(), s.config.EnforceSingleBlobBatches)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize batch: %v", err)
	}

	return batch, nil
}

func (s *ServerV2) GetChunks(ctx context.Context, in *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	start := time.Now()

	if !s.config.EnableV2 {
		return nil, api.NewErrorInvalidArg("v2 API is disabled")
	}

	blobKey, err := corev2.BytesToBlobKey(in.GetBlobKey())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob key: %v", err))
	}

	if corev2.MaxQuorumID < in.GetQuorumId() {
		//nolint: wrapcheck
		return nil, api.NewErrorInvalidArg(
			fmt.Sprintf("quorumID %d must be <= maxQuorumID %d", in.GetQuorumId(), corev2.MaxQuorumID))
	}

	// The current sampling scheme will store the same chunks for all quorums, so we always use quorum 0 as the quorum key in storage.
	quorumID := core.QuorumID(0)

	bundleKey, err := node.BundleKey(blobKey, quorumID)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to get bundle key: %v", err))
	}

	bundleData, err := s.node.ValidatorStore.GetBundleData(bundleKey)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get chunks: %v", err))
	}

	chunks, _, err := node.DecodeChunks(bundleData)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to decode chunks: %v", err))
	}

	size := 0
	if len(chunks) > 0 {
		size = len(chunks[0]) * len(chunks)
	}
	s.metrics.ReportGetChunksDataSize(size)

	s.metrics.ReportGetChunksLatency(time.Since(start))

	return &pb.GetChunksReply{
		Chunks:              chunks,
		ChunkEncodingFormat: pb.ChunkEncodingFormat_GNARK,
	}, nil
}

// validateDispersalRequest validates the DisperseBlobRequest and returns the blob header
// Differences between this and the DispersalServerV2 are:
// - Takes *corev2.BlobCertificate instead of DisperseBlobRequest
// - no encoding prover GetCommitmentsForPaddedLength check
// - directly take blob lengths (no blob data yet)
// - doesn't check every 32 bytes is a valid field element
// Node cannot make these checks because the checks require the blob data
func (s *ServerV2) validateDispersalRequest(
	blobCert *corev2.BlobCertificate,
) (*corev2.BlobHeader, error) {
	if len(blobCert.Signature) != 65 {
		return nil, fmt.Errorf("signature is expected to be 65 bytes, but got %d bytes", len(blobCert.Signature))
	}
	err := s.blobAuthenticator.AuthenticateBlobRequest(blobCert.BlobHeader, blobCert.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate blob request: %v", err)
	}

	// this is the length in SYMBOLS (32 byte field elements) of the blob. it must be a power of 2
	committedBlobLength := blobCert.BlobHeader.BlobCommitments.Length
	if committedBlobLength == 0 {
		return nil, errors.New("blob size must be greater than 0")
	}
	if uint64(committedBlobLength) != math.NextPowOf2u64(uint64(committedBlobLength)) {
		return nil, errors.New("invalid commitment length, must be a power of 2")
	}

	blobHeader := blobCert.BlobHeader
	if blobHeader.PaymentMetadata == (core.PaymentMetadata{}) {
		return nil, errors.New("payment metadata is required")
	}

	timestampIsNegative := blobHeader.PaymentMetadata.Timestamp < 0
	paymentIsNegative := blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == -1
	timestampIsZeroAndPaymentIsZero := blobHeader.PaymentMetadata.Timestamp == 0 && blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0
	if timestampIsNegative || paymentIsNegative || timestampIsZeroAndPaymentIsZero {
		return nil, errors.New("invalid payment metadata")
	}

	if len(blobHeader.QuorumNumbers) == 0 {
		return nil, errors.New("blob header must contain at least one quorum number")
	}

	if len(blobHeader.QuorumNumbers) > int(s.node.QuorumCount.Load()) {
		return nil, fmt.Errorf("too many quorum numbers specified: maximum is %d", s.node.QuorumCount.Load())
	}

	for _, quorum := range blobHeader.QuorumNumbers {
		if quorum > corev2.MaxQuorumID || quorum >= uint8(s.node.QuorumCount.Load()) {
			return nil, fmt.Errorf("invalid quorum number %d; maximum is %d", quorum, s.node.QuorumCount.Load())
		}
	}

	if _, ok := s.node.BlobVersionParams.Load().Get(corev2.BlobVersion(blobHeader.BlobVersion)); !ok {
		return nil, fmt.Errorf("invalid blob version %d; valid blob versions are: %v", blobHeader.BlobVersion, s.node.BlobVersionParams.Load().Keys())
	}

	return blobHeader, nil
}
