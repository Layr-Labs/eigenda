package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/core"
	coreauthv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
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
}

// NewServerV2 creates a new Server instance with the provided parameters.
func NewServerV2(
	ctx context.Context,
	config *node.Config,
	node *node.Node,
	logger logging.Logger,
	ratelimiter common.RateLimiter,
	registry *prometheus.Registry,
	reader core.Reader) (*ServerV2, error) {

	metrics, err := NewV2Metrics(logger, registry)
	if err != nil {
		return nil, err
	}

	var chunkAuthenticator auth.RequestAuthenticator
	var blobAuthenticator corev2.BlobRequestAuthenticator
	if !config.DisableDispersalAuthentication {
		chunkAuthenticator, err = auth.NewRequestAuthenticator(
			ctx,
			reader,
			config.DispersalAuthenticationKeyCacheSize,
			config.DisperserKeyTimeout,
			func(id uint32) bool {
				return id == api.EigenLabsDisperserID
			},
			time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator: %w", err)
		}
		blobAuthenticator = coreauthv2.NewBlobRequestAuthenticator()
	}
	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.StoreChunksRequestMaxPastAge,
		config.StoreChunksRequestMaxFutureAge)

	return &ServerV2{
		config:             config,
		node:               node,
		ratelimiter:        ratelimiter,
		logger:             logger,
		metrics:            metrics,
		chunkAuthenticator: chunkAuthenticator,
		blobAuthenticator:  blobAuthenticator,
		replayGuardian:     replayGuardian,
	}, nil
}

func (s *ServerV2) GetNodeInfo(ctx context.Context, in *pb.GetNodeInfoRequest) (*pb.GetNodeInfoReply, error) {
	if s.config.DisableNodeInfoResources {
		return &pb.GetNodeInfoReply{Semver: node.SemVer}, nil
	}

	memBytes := uint64(0)
	v, err := mem.VirtualMemory()
	if err == nil {
		memBytes = v.Total
	}

	return &pb.GetNodeInfoReply{
		Semver:   node.SemVer,
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

	// check the blacklist to see if we need to proceed or not
	disperserAddress, err := s.chunkAuthenticator.GetDisperserAddress(ctx, in.DisperserID, time.Now())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to get disperser address: %v", err))
	}

	if s.node.BlacklistStore.IsBlacklisted(ctx, disperserAddress.Bytes()) {
		return nil, api.NewErrorInvalidArg("disperser is blacklisted")
	}
	if s.chunkAuthenticator != nil {
		hash, err := s.chunkAuthenticator.AuthenticateStoreChunksRequest(ctx, in, time.Now())
		if err != nil {
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to authenticate request: %v", err))
		}

		timestamp := time.Unix(int64(in.Timestamp), 0)
		err = s.replayGuardian.VerifyRequest(hash, timestamp)
		if err != nil {
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request: %v", err))
		}
	}
	if s.blobAuthenticator != nil {
		// TODO: check the latency of request validation later; could be parallelized to avoid significant
		// impact to the request latency
		for _, blobCert := range batch.BlobCertificates {
			_, err = s.validateDispersalRequest(blobCert)
			if err != nil {
				// Blacklist the disperser if there's an invalid dispersal request
				s.blacklistDisperser(in, blobCert)
				return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate blob request: %v", err))
			}
		}
	}
	probe.SetStage("get_operator_state")
	s.logger.Info("new StoreChunks request",
		"batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]),
		"numBlobs", len(batch.BlobCertificates),
		"referenceBlockNumber", batch.BatchHeader.ReferenceBlockNumber)
	operatorState, err := s.node.ChainState.GetOperatorStateByOperator(
		ctx,
		uint(batch.BatchHeader.ReferenceBlockNumber),
		s.node.Config.ID)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get the operator state: %v", err))
	}

	blobShards, rawBundles, err := s.node.DownloadBundles(ctx, batch, operatorState, probe)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get the operator state: %v", err))
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

	return &pb.StoreChunksReply{
		Signature: sig,
	}, nil
}

func (s *ServerV2) validateAndStoreChunks(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	rawBundles []*node.RawBundles,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte,
	probe *common.SequenceProbe,
) error {

	bundleCount := 0
	for _, bundles := range rawBundles {
		bundleCount += len(bundles.Bundles)
	}

	batchData := make([]*node.BundleToStore, 0, bundleCount)
	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return api.NewErrorInternal("failed to get blob key")
		}

		for quorum, bundle := range bundles.Bundles {
			bundleKey, err := node.BundleKey(blobKey, quorum)
			if err != nil {
				return api.NewErrorInternal("failed to get bundle key")
			}

			batchData = append(batchData, &node.BundleToStore{
				BundleKey:   bundleKey,
				BundleBytes: bundle,
			})
		}
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

// blacklistDisperser blacklists a disperser by retrieving the disperser's public key from the request and storing it in the blacklist store
func (s *ServerV2) blacklistDisperser(request *pb.StoreChunksRequest, blobCert *corev2.BlobCertificate) error {

	ctx := context.Background()

	// using the pubkey here since disperserId can be claimed by others and has some edge cases, so will avoid that.
	disperserAddress, err := s.chunkAuthenticator.GetDisperserAddress(ctx, request.DisperserID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get disperser address: %w", err)
	}

	// Get blob key for context
	blobKey, err := blobCert.BlobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %w", err)
	}

	// Store bytes once to avoid repeated conversions
	disperserBytes := disperserAddress.Bytes()

	err = s.node.BlacklistStore.AddEntry(ctx, disperserBytes, fmt.Sprintf("blobKey: %x", blobKey), "blobCert validation failed")
	if err != nil {
		return fmt.Errorf("failed to add entry to blacklist: %w", err)
	}
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
	batch, err := corev2.BatchFromProtobuf(req.GetBatch())
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
		return nil, api.NewErrorInvalidArg("invalid quorum ID")
	}
	quorumID := core.QuorumID(in.GetQuorumId())

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
	commitedBlobLength := blobCert.BlobHeader.BlobCommitments.Length
	if commitedBlobLength == 0 {
		return nil, errors.New("blob size must be greater than 0")
	}
	if commitedBlobLength != encoding.NextPowerOf2(commitedBlobLength) {
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
