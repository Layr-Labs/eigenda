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
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/core"
	authv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/meterer"
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
	meterer            *meterer.Meterer
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
	meterer *meterer.Meterer,
	registry *prometheus.Registry,
	reader core.Reader) (*ServerV2, error) {

	metrics, err := NewV2Metrics(logger, registry)
	if err != nil {
		return nil, err
	}

	chunkAuthenticator, err := auth.NewRequestAuthenticator(
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
	blobAuthenticator := authv2.NewPaymentStateAuthenticator(
		//TODO: add as configs
		time.Duration(10)*time.Second,
		time.Duration(10)*time.Second,
	)

	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.StoreChunksRequestMaxPastAge,
		config.StoreChunksRequestMaxFutureAge)

	return &ServerV2{
		config:             config,
		node:               node,
		ratelimiter:        ratelimiter,
		meterer:            meterer,
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
	batch, err := s.validateStoreChunksRequest(ctx, in)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate store chunk request: %v", err))
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to serialize batch header hash: %v", err))
	}

	hash, err := s.chunkAuthenticator.AuthenticateStoreChunksRequest(ctx, in, time.Now())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to authenticate request: %v", err))
	}

	timestamp := time.Unix(int64(in.Timestamp), 0)
	err = s.replayGuardian.VerifyRequest(hash, timestamp)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request: %v", err))
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

	batchData := make([]*node.BundleToStore, 0, len(rawBundles))
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

	if s.config.LittDBEnabled {
		return s.validateAndStoreChunksLittDB(
			ctx,
			batch,
			blobShards,
			batchData,
			operatorState,
			batchHeaderHash,
			probe)
	} else {
		probe.SetStage("validate_and_store")
		return s.validateAndStoreChunksLevelDB(ctx, batch, blobShards, batchData, operatorState, batchHeaderHash)
	}
}

func (s *ServerV2) validateAndStoreChunksLevelDB(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	batchData []*node.BundleToStore,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte) error {

	type storeResult struct {
		keys []kvstore.Key
		err  error
	}
	storeChan := make(chan storeResult)
	go func() {
		keys, size, err := s.node.ValidatorStore.StoreBatch(batchHeaderHash[:], batchData)
		if err != nil {
			storeChan <- storeResult{
				keys: nil,
				err:  err,
			}
			return
		}

		s.metrics.ReportStoreChunksRequestSize(size)
		storeChan <- storeResult{
			keys: keys,
			err:  nil,
		}
	}()

	err := s.node.ValidateBatchV2(ctx, batch, blobShards, operatorState)
	if err != nil {
		res := <-storeChan
		if len(res.keys) > 0 {
			if deleteErr := s.node.ValidatorStore.DeleteKeys(res.keys); deleteErr != nil {
				s.logger.Error(
					"failed to delete keys",
					"err", deleteErr,
					"batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]))
			}
		}
		return api.NewErrorInternal(fmt.Sprintf("failed to validate batch: %v", err))
	}

	res := <-storeChan
	if res.err != nil {
		return api.NewErrorInternal(fmt.Sprintf("failed to store batch: %v", res.err))
	}

	return nil
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
	_, size, err := s.node.ValidatorStore.StoreBatch(batchHeaderHash[:], batchData)
	if err != nil {
		return api.NewErrorInternal(
			fmt.Sprintf("failed to store batch %s: %v", hex.EncodeToString(batchHeaderHash[:]), err))
	}

	s.metrics.ReportStoreChunksRequestSize(size)

	return nil
}

// validateStoreChunksRequest validates the StoreChunksRequest and returns deserialized batch in the request
func (s *ServerV2) validateStoreChunksRequest(ctx context.Context, req *pb.StoreChunksRequest) (*corev2.Batch, error) {
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

	for _, cert := range batch.BlobCertificates {
		err := s.validateDispersalRequest(ctx, cert, s.node.LoadOnchainState())
		if err != nil {
			return nil, fmt.Errorf("failed to validate dispersal request: %v", err)
		}
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

func (s *ServerV2) validateDispersalRequest(ctx context.Context, blobCert *corev2.BlobCertificate, onchainState *eth.OnchainState) error {
	if len(blobCert.Signature) != 65 {
		return api.NewErrorInvalidArg(fmt.Sprintf("signature is expected to be 65 bytes, but got %d bytes", len(blobCert.Signature)))
	}

	blobLength := blobCert.BlobHeader.BlobCommitments.Length
	if blobLength == 0 {
		return api.NewErrorInvalidArg("blob length must be greater than 0")
	}
	if blobLength > uint(onchainState.MaxNumSymbolsPerBlob) {
		return api.NewErrorInvalidArg("blob length too big")
	}
	if blobLength != encoding.NextPowerOf2(blobLength) {
		return api.NewErrorInvalidArg("invalid blob length, must be a power of 2")
	}

	blobHeader := blobCert.BlobHeader
	if blobHeader.PaymentMetadata == (core.PaymentMetadata{}) {
		return api.NewErrorInvalidArg("payment metadata is required")
	}

	if len(blobHeader.PaymentMetadata.AccountID) == 0 || (blobHeader.PaymentMetadata.Timestamp == 0 && blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0) {
		return api.NewErrorInvalidArg("invalid payment metadata")
	}

	quorumNumbers := blobHeader.QuorumNumbers
	if len(quorumNumbers) == 0 {
		return api.NewErrorInvalidArg("blob header must contain at least one quorum number")
	}

	if len(quorumNumbers) > int(onchainState.QuorumCount) {
		return api.NewErrorInvalidArg(fmt.Sprintf("too many quorum numbers specified: maximum is %d", onchainState.QuorumCount))
	}

	for _, quorum := range quorumNumbers {
		if quorum > corev2.MaxQuorumID || uint8(quorum) >= onchainState.QuorumCount {
			return api.NewErrorInvalidArg(fmt.Sprintf("invalid quorum number %d; maximum is %d", quorum, onchainState.QuorumCount))
		}
	}

	if _, ok := onchainState.BlobVersionParameters.Get(corev2.BlobVersion(blobHeader.BlobVersion)); !ok {
		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob version %d; valid blob versions are: %v", blobHeader.BlobVersion, onchainState.BlobVersionParameters.Keys()))
	}

	if err := s.blobAuthenticator.AuthenticateBlobRequest(blobHeader, blobCert.Signature); err != nil {
		return api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
	}

	// handle payments and check rate limits
	timestamp := blobHeader.PaymentMetadata.Timestamp
	cumulativePayment := blobHeader.PaymentMetadata.CumulativePayment
	accountID := blobHeader.PaymentMetadata.AccountID

	paymentHeader := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	// TODO: meterer must be updated to limit on-demand requests to the EigenDA disperser
	if _, err := s.meterer.MeterRequest(ctx, paymentHeader, uint64(blobLength), blobHeader.QuorumNumbers, time.Now()); err != nil {
		return api.NewErrorResourceExhausted(err.Error())
	}

	// should node run prover here? doesn't have data yet though
	return nil
}
