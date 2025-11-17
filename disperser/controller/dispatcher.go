package controller

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
)

var errNoBlobsToDispatch = errors.New("no blobs to dispatch")

type BlobCallback func(blobKey corev2.BlobKey) error

// DispatcherConfig contains configuration parameters for the Dispatcher.
// The Dispatcher is responsible for batching encoded blobs, dispersing them to DA nodes,
// collecting signatures, and creating attestations.
type DispatcherConfig struct {
	// PullInterval is how frequently the Dispatcher polls for new encoded blobs to batch and dispatch.
	// Must be positive.
	PullInterval time.Duration

	// DisperserID is the unique identifier for this disperser instance.
	DisperserID uint32

	// FinalizationBlockDelay is the number of blocks to wait before using operator state.
	// This provides a hedge against chain reorganizations.
	FinalizationBlockDelay uint64

	// BatchMetadataUpdatePeriod is the interval between attempts to refresh batch metadata
	// (reference block number and operator state).
	// Since this changes at most once per eth block, values shorter than 10 seconds are not useful.
	// In practice, checking every several minutes is sufficient.
	// Must be positive.
	BatchMetadataUpdatePeriod time.Duration

	// AttestationTimeout is the maximum time to wait for a single node to provide a signature.
	// Must be positive.
	AttestationTimeout time.Duration

	// BatchAttestationTimeout is the maximum time to wait for all nodes to provide signatures for a batch.
	// Must be positive and must be longer or equal to the AttestationTimeout.
	BatchAttestationTimeout time.Duration

	// SignatureTickInterval is how frequently attestations are updated in the blob metadata store
	// as signature gathering progresses.
	// Must be positive.
	SignatureTickInterval time.Duration

	// MaxBatchSize is the maximum number of blobs to include in a single batch for dispersal.
	// Must be at least 1.
	MaxBatchSize int32

	// SignificantSigningThresholdFraction is a configurable "important" signing threshold fraction.
	// Used to track signing metrics and understand system performance.
	// If the value is 0, special handling for this threshold is disabled.
	// Must be between 0.0 and 1.0.
	SignificantSigningThresholdFraction float64

	// Whether or not to collect detailed validator signing metrics.
	CollectDetailedValidatorSigningMetrics bool

	// NumConcurrentRequests is the size of the worker pool for processing dispersal requests concurrently.
	// Must be at least 1.
	NumConcurrentRequests int

	// NodeClientCacheSize is the maximum number of node clients to cache for reuse.
	// Must be at least 1.
	NodeClientCacheSize int

	// MaxDispersalAge is the maximum age a dispersal request can be before it is discarded.
	// Dispersals older than this duration are marked as Failed and not processed.
	//
	// Age is determined by the BlobHeader.PaymentMetadata.Timestamp field, which is set by the
	// client at dispersal request creation time (in nanoseconds since Unix epoch).
	MaxDispersalAge time.Duration
}

var _ config.VerifiableConfig = &DispatcherConfig{}

func (c *DispatcherConfig) Verify() error {
	if c.PullInterval <= 0 {
		return fmt.Errorf("PullInterval must be positive, got %v", c.PullInterval)
	}
	if c.BatchMetadataUpdatePeriod <= 0 {
		return fmt.Errorf("BatchMetadataUpdatePeriod must be positive, got %v", c.BatchMetadataUpdatePeriod)
	}
	if c.AttestationTimeout <= 0 {
		return fmt.Errorf("AttestationTimeout must be positive, got %v", c.AttestationTimeout)
	}
	if c.BatchAttestationTimeout <= 0 {
		return fmt.Errorf("BatchAttestationTimeout must be positive, got %v", c.BatchAttestationTimeout)
	}
	if c.BatchAttestationTimeout < c.AttestationTimeout {
		return fmt.Errorf("BatchAttestationTimeout must be longer than AttestationTimeout, got %v < %v",
			c.BatchAttestationTimeout, c.AttestationTimeout)
	}
	if c.SignatureTickInterval <= 0 {
		return fmt.Errorf("SignatureTickInterval must be positive, got %v", c.SignatureTickInterval)
	}
	if c.MaxBatchSize < 1 {
		return fmt.Errorf("MaxBatchSize must be at least 1, got %d", c.MaxBatchSize)
	}
	if c.SignificantSigningThresholdFraction > 1.0 || c.SignificantSigningThresholdFraction < 0.0 {
		return fmt.Errorf(
			"SignificantSigningThresholdFraction must be between 0.0 and 1.0, got %f",
			c.SignificantSigningThresholdFraction)
	}
	if c.NumConcurrentRequests < 1 {
		return fmt.Errorf("NumConcurrentRequests must be at least 1, got %d", c.NumConcurrentRequests)
	}
	if c.NodeClientCacheSize < 1 {
		return fmt.Errorf("NodeClientCacheSize must be at least 1, got %d", c.NodeClientCacheSize)
	}
	if c.MaxDispersalAge <= 0 {
		return fmt.Errorf("MaxDispersalAge must be positive, got %v", c.MaxDispersalAge)
	}
	return nil
}

func DefaultDispatcherConfig() *DispatcherConfig {
	return &DispatcherConfig{
		PullInterval:                        1 * time.Second,
		FinalizationBlockDelay:              75,
		AttestationTimeout:                  45 * time.Second,
		BatchMetadataUpdatePeriod:           time.Minute,
		BatchAttestationTimeout:             55 * time.Second,
		SignatureTickInterval:               50 * time.Millisecond,
		MaxBatchSize:                        32,
		SignificantSigningThresholdFraction: 0.55,
		NumConcurrentRequests:               600,
		NodeClientCacheSize:                 400,
		MaxDispersalAge:                     45 * time.Second,
	}
}

type Dispatcher struct {
	*DispatcherConfig

	blobMetadataStore blobstore.MetadataStore
	pool              common.WorkerPool
	chainState        core.IndexedChainState
	aggregator        core.SignatureAggregator
	nodeClientManager NodeClientManager
	logger            logging.Logger
	metrics           *controllerMetrics
	getNow            func() time.Time

	cursor *blobstore.StatusIndexCursor

	// beforeDispatch function is called before dispatching a blob
	beforeDispatch BlobCallback

	// blobSet keeps track of blobs that are being dispatched
	// This is used to deduplicate blobs to prevent the same blob from being dispatched multiple times
	// Blobs are removed from the queue when they are in a terminal state (Complete or Failed)
	blobSet BlobSet

	controllerLivenessChan chan<- healthcheck.HeartbeatMessage

	// A utility responsible for fetching batch metadata (i.e. reference block number and operator state).
	batchMetadataManager metadata.BatchMetadataManager
}

type batchData struct {
	Batch           *corev2.Batch
	BatchHeaderHash [32]byte
	BlobKeys        []corev2.BlobKey
	Metadata        map[corev2.BlobKey]*v2.BlobMetadata
	OperatorState   *core.IndexedOperatorState
	BatchSizeBytes  uint64
}

func NewDispatcher(
	config *DispatcherConfig,
	getNow func() time.Time,
	blobMetadataStore blobstore.MetadataStore,
	pool common.WorkerPool,
	chainState core.IndexedChainState,
	batchMetadataManager metadata.BatchMetadataManager,
	aggregator core.SignatureAggregator,
	nodeClientManager NodeClientManager,
	logger logging.Logger,
	registry *prometheus.Registry,
	beforeDispatch func(blobKey corev2.BlobKey) error,
	blobSet BlobSet,
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage,
) (*Dispatcher, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	metrics, err := newControllerMetrics(
		registry,
		config.SignificantSigningThresholdFraction,
		config.CollectDetailedValidatorSigningMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %v", err)
	}

	return &Dispatcher{
		DispatcherConfig: config,

		blobMetadataStore: blobMetadataStore,
		pool:              pool,
		chainState:        chainState,
		aggregator:        aggregator,
		nodeClientManager: nodeClientManager,
		logger:            logger.With("component", "Dispatcher"),
		metrics:           metrics,
		getNow:            getNow,

		cursor:                 nil,
		beforeDispatch:         beforeDispatch,
		blobSet:                blobSet,
		controllerLivenessChan: controllerLivenessChan,
		batchMetadataManager:   batchMetadataManager,
	}, nil
}

func (d *Dispatcher) Start(ctx context.Context) error {
	err := d.chainState.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start chain state: %w", err)
	}

	go func() {
		ticker := time.NewTicker(d.PullInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				attestationCtx, cancel := context.WithTimeout(ctx, d.BatchAttestationTimeout)
				probe := d.metrics.newBatchProbe()

				sigChan, batchData, err := d.HandleBatch(attestationCtx, probe)
				if err != nil {
					if errors.Is(err, errNoBlobsToDispatch) {
						d.logger.Debug("no blobs to dispatch")
					} else {
						d.logger.Error("failed to process a batch", "err", err)
					}
					cancel()
					probe.End()
					continue
				}
				go func() {
					probe.SetStage("handle_signatures")
					err := d.HandleSignatures(ctx, attestationCtx, batchData, sigChan)
					if err != nil {
						d.logger.Error("failed to handle signatures", "err", err)
					}
					cancel()
					probe.End()
				}()
			}
		}
	}()

	return nil

}

// For each blob in a batch, send a StoreChunks request to each validator, collecting responses and putting those
// responses in the returned channel.
func (d *Dispatcher) HandleBatch(
	ctx context.Context,
	batchProbe *common.SequenceProbe,
) (chan core.SigningMessage, *batchData, error) {
	// Signal Liveness to indicate no stall
	healthcheck.SignalHeartbeat(d.logger, "dispatcher", d.controllerLivenessChan)

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob inclusion info for each blob in metadata store
	batchData, err := d.NewBatch(ctx, batchProbe)
	if err != nil {
		return nil, nil, err
	}

	batchProbe.SetStage("send_requests")

	signingResponseChan := make(chan core.SigningMessage, len(batchData.OperatorState.IndexedOperators))
	for validatorId, validatorInfo := range batchData.OperatorState.IndexedOperators {

		validatorProbe := d.metrics.newSendToValidatorProbe()
		validatorProbe.SetStage("pool_submission")

		d.pool.Submit(func() {
			signature, latency, err := d.sendChunksToValidator(
				ctx,
				batchData,
				validatorId,
				validatorInfo,
				validatorProbe)

			if err != nil {
				d.logger.Warn("error sending chunks to validator",
					"validator", validatorId.Hex(),
					"batchHeaderHash", hex.EncodeToString(batchData.BatchHeaderHash[:]),
					"err", err)
			}

			signingResponseChan <- core.SigningMessage{
				ValidatorId:     validatorId,
				Signature:       signature,
				BatchHeaderHash: batchData.BatchHeaderHash,
				Latency:         latency,
				Err:             err,
			}
		})
	}

	batchProbe.SetStage("await_responses")

	return signingResponseChan, batchData, nil
}

// Send a StoreChunks request for a batch to a specific validator, returning the result.
func (d *Dispatcher) sendChunksToValidator(
	ctx context.Context,
	batchData *batchData,
	validatorId core.OperatorID,
	validatorInfo *core.IndexedOperatorInfo,
	validatorProbe *common.SequenceProbe,
) (signature *core.Signature, latency time.Duration, err error) {

	defer validatorProbe.End()

	validatorProbe.SetStage("get_client")

	host, _, _, v2DispersalPort, _, err := core.ParseOperatorSocket(validatorInfo.Socket)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse operator socket %s: %w", validatorInfo.Socket, err)
	}

	client, err := d.nodeClientManager.GetClient(host, v2DispersalPort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get node client for validator at port %s: %w", v2DispersalPort, err)
	}

	validatorProbe.SetStage("put_dispersal_request")

	req := &corev2.DispersalRequest{
		OperatorID: validatorId,
		// TODO: get OperatorAddress
		OperatorAddress: gethcommon.Address{},
		Socket:          validatorInfo.Socket,
		DispersedAt:     uint64(time.Now().UnixNano()),
		BatchHeader:     *batchData.Batch.BatchHeader,
	}
	err = d.blobMetadataStore.PutDispersalRequest(ctx, req)
	if err != nil {
		d.logger.Error("failed to put dispersal request", "err", err)

		return nil, 0, fmt.Errorf("failed to put dispersal request for validator: %w", err)
	}

	validatorProbe.SetStage("send_chunks")

	start := time.Now()

	sig, err := d.sendChunks(ctx, client, batchData.Batch)
	if err != nil {
		storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
			DispersalRequest: req,
			RespondedAt:      uint64(time.Now().UnixNano()),
			Signature:        [32]byte{}, // all zero sig for failed dispersal
			Error:            err.Error(),
		})
		if storeErr != nil {
			d.logger.Error("failed to store a failed dispersal response", "err", storeErr)
		}
		return nil, 0, fmt.Errorf("failed to send chunks to validator: %w", err)
	}

	latency = time.Since(start)

	validatorProbe.SetStage("put_dispersal_response")
	storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
		DispersalRequest: req,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        sig.Bytes(),
	})
	if storeErr != nil {
		d.logger.Error("failed to store a succeeded dispersal response", "err", storeErr)
	}

	return sig, latency, nil
}

// HandleSignatures receives SigningMessages from operators for a given batch through the input sigChan. The signatures
// are validated, aggregated, and used to put an Attestation for the batch into the blobMetadataStore. The Attestation
// is periodically updated as additional signatures are gathered.
//
// This method will continue gathering signatures until a SigningMessage has been received from every operator, or until
// the global attestationCtx times out.
func (d *Dispatcher) HandleSignatures(
	ctx context.Context,
	attestationCtx context.Context,
	batchData *batchData,
	sigChan chan core.SigningMessage,
) error {
	if batchData == nil {
		return errors.New("batchData is required")
	}

	batchHeaderHash := hex.EncodeToString(batchData.BatchHeaderHash[:])
	for _, key := range batchData.BlobKeys {
		err := d.blobMetadataStore.UpdateBlobStatus(ctx, key, v2.GatheringSignatures)
		if err != nil {
			d.logger.Error("failed to update blob status to 'gathering signatures'",
				"blobKey", key.Hex(),
				"batchHeaderHash", batchHeaderHash,
				"err", err)
		}
	}

	// write an empty attestation before starting to gather signatures, so that it can be queried right away.
	// the attestation will be periodically updated as signatures are gathered.
	attestation := &corev2.Attestation{
		BatchHeader:      batchData.Batch.BatchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: nil,
		APKG2:            nil,
		QuorumAPKs:       nil,
		Sigma:            nil,
		QuorumNumbers:    nil,
		QuorumResults:    nil,
	}
	err := d.blobMetadataStore.PutAttestation(ctx, attestation)
	if err != nil {
		// this error isn't fatal: a subsequent PutAttestation attempt might succeed
		d.logger.Error("error calling PutAttestation",
			"err", err,
			"batchHeaderHash", batchHeaderHash)
	}

	// This channel will remain open until the attestationTimeout triggers, or until signatures from all validators
	// have been received and processed. It will periodically yield QuorumAttestations with the latest set of received
	// signatures.
	attestationChan, err := ReceiveSignatures(
		attestationCtx,
		d.logger,
		d.metrics,
		batchData.OperatorState,
		batchData.BatchHeaderHash,
		sigChan,
		d.DispatcherConfig.SignatureTickInterval,
		d.DispatcherConfig.SignificantSigningThresholdFraction,
		batchData.BatchSizeBytes)
	if err != nil {
		receiveSignaturesErr := fmt.Errorf("receive and validate signatures for batch %s: %w", batchHeaderHash, err)

		dbErr := d.failBatch(ctx, batchData)
		if dbErr != nil {
			return multierror.Append(
				receiveSignaturesErr,
				fmt.Errorf("update blob statuses for batch to 'failed': %w", dbErr))
		}

		return receiveSignaturesErr
	}

	// keep track of the final attestation, since that's the attestation which will determine the final batch status
	finalAttestation := &core.QuorumAttestation{}
	// continue receiving attestations from the channel until it's closed
	for receivedQuorumAttestation := range attestationChan {
		err := d.updateAttestation(ctx, batchData, receivedQuorumAttestation)
		if err != nil {
			d.logger.Warnf("error updating attestation for batch %s: %v", batchHeaderHash, err)
			continue
		}

		finalAttestation = receivedQuorumAttestation
	}

	updateBatchStatusStartTime := time.Now()
	_, quorumPercentages := d.parseQuorumPercentages(finalAttestation.QuorumResults)
	err = d.updateBatchStatus(ctx, batchData, quorumPercentages)
	d.metrics.reportUpdateBatchStatusLatency(time.Since(updateBatchStatusStartTime))
	if err != nil {
		return fmt.Errorf("update batch status: %w", err)
	}

	return nil
}

// updateAttestation updates the QuorumAttestation in the blobMetadataStore
func (d *Dispatcher) updateAttestation(
	ctx context.Context,
	batchData *batchData,
	quorumAttestation *core.QuorumAttestation,
) error {
	sortedNonZeroQuorums, quorumPercentages := d.parseQuorumPercentages(quorumAttestation.QuorumResults)
	if len(sortedNonZeroQuorums) == 0 {
		return errors.New("all quorums received no attestation for batch")
	}

	aggregationStartTime := time.Now()
	signatureAggregation, err := d.aggregator.AggregateSignatures(
		batchData.OperatorState,
		quorumAttestation,
		sortedNonZeroQuorums)
	d.metrics.reportAggregateSignaturesLatency(time.Since(aggregationStartTime))
	if err != nil {
		return fmt.Errorf("aggregate signatures: %w", err)
	}

	attestation := &corev2.Attestation{
		BatchHeader:      batchData.Batch.BatchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: signatureAggregation.NonSigners,
		APKG2:            signatureAggregation.AggPubKey,
		QuorumAPKs:       signatureAggregation.QuorumAggPubKeys,
		Sigma:            signatureAggregation.AggSignature,
		QuorumNumbers:    sortedNonZeroQuorums,
		QuorumResults:    quorumPercentages,
	}

	putAttestationStartTime := time.Now()
	err = d.blobMetadataStore.PutAttestation(ctx, attestation)
	d.metrics.reportPutAttestationLatency(time.Since(putAttestationStartTime))
	if err != nil {
		return fmt.Errorf("put attestation: %w", err)
	}

	d.logAttestationUpdate(hex.EncodeToString(batchData.BatchHeaderHash[:]), quorumAttestation.QuorumResults)

	return nil
}

// parseQuorumPercentages iterates over the map of QuorumResults, and returns a sorted slice of nonZeroQuorums
// (quorums with >0 signing percentage), and a map from QuorumID to signing percentage.
func (d *Dispatcher) parseQuorumPercentages(
	quorumResults map[core.QuorumID]*core.QuorumResult,
) ([]core.QuorumID, map[core.QuorumID]uint8) {
	nonZeroQuorums := make([]core.QuorumID, 0)
	quorumPercentages := make(map[core.QuorumID]uint8)

	for quorumID, quorumResult := range quorumResults {
		if quorumResult.PercentSigned > 0 {
			nonZeroQuorums = append(nonZeroQuorums, quorumID)
			quorumPercentages[quorumID] = quorumResult.PercentSigned
		}
	}

	slices.Sort(nonZeroQuorums)

	return nonZeroQuorums, quorumPercentages
}

// logAttestationUpdate logs the attestation details, including batch header hash and quorum signing percentages
func (d *Dispatcher) logAttestationUpdate(batchHeaderHash string, quorumResults map[core.QuorumID]*core.QuorumResult) {
	quorumPercentagesBuilder := strings.Builder{}
	quorumPercentagesBuilder.WriteString("(")

	for quorumID, quorumResult := range quorumResults {
		quorumPercentagesBuilder.WriteString(
			fmt.Sprintf("quorum_%d: %d%%, ", quorumID, quorumResult.PercentSigned))
	}
	quorumPercentagesBuilder.WriteString(")")

	d.logger.Debug("attestation updated",
		"batchHeaderHash", batchHeaderHash,
		"quorumPercentages", quorumPercentagesBuilder.String())
}

// Iterates over the input metadata slice, and returns a new slice with stale and duplicate metadatas filtered out
func (d *Dispatcher) filterStaleAndDedupBlobs(
	ctx context.Context,
	inputMetadatas []*v2.BlobMetadata,
) []*v2.BlobMetadata {
	outputMetadatas := make([]*v2.BlobMetadata, 0, len(inputMetadatas))
	now := d.getNow()

	for _, metadata := range inputMetadatas {
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			d.logger.Errorf("compute blob key: %w", err)
			// we must discard if we cannot compute key, since it's used for deduplication
			continue
		}

		if d.checkAndHandleStaleBlob(ctx, blobKey, now, metadata.BlobHeader.PaymentMetadata.Timestamp) {
			// discard stale blob
			continue
		}

		if d.blobSet.Contains(blobKey) {
			// discard duplicate blob
			continue
		}

		outputMetadatas = append(outputMetadatas, metadata)
	}

	return outputMetadatas
}

// NewBatch creates a batch of blobs to dispatch
// Warning: This function is not thread-safe
func (d *Dispatcher) NewBatch(
	ctx context.Context,
	probe *common.SequenceProbe,
) (*batchData, error) {

	batchMetadata := d.batchMetadataManager.GetMetadata()
	referenceBlockNumber := batchMetadata.ReferenceBlockNumber()
	operatorState := batchMetadata.OperatorState()

	probe.SetStage("get_blob_metadata")
	blobMetadatas, cursor, err := d.blobMetadataStore.GetBlobMetadataByStatusPaginated(
		ctx,
		v2.Encoded,
		d.cursor,
		d.MaxBatchSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata by status: %w", err)
	}

	blobMetadatas = d.filterStaleAndDedupBlobs(ctx, blobMetadatas)
	d.metrics.reportBlobSetSize(d.blobSet.Size())
	if len(blobMetadatas) == 0 {
		return nil, errNoBlobsToDispatch
	}
	d.logger.Debug("got new metadatas to make batch",
		"numBlobs", len(blobMetadatas),
		"referenceBlockNumber", referenceBlockNumber)

	keys := make([]corev2.BlobKey, len(blobMetadatas))
	metadataMap := make(map[corev2.BlobKey]*v2.BlobMetadata, len(blobMetadatas))
	for i, metadata := range blobMetadatas {
		if metadata == nil || metadata.BlobHeader == nil {
			return nil, fmt.Errorf("invalid blob metadata")
		}
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}
		keys[i] = blobKey
		metadataMap[blobKey] = metadata

		if d.beforeDispatch != nil {
			err = d.beforeDispatch(blobKey)
			if err != nil {
				d.logger.Error("beforeDispatch function failed", "blobKey", blobKey.Hex(), "err", err)
			}
		}
	}

	probe.SetStage("get_blob_certs")
	certs, _, err := d.blobMetadataStore.GetBlobCertificates(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob certificates: %w", err)
	}

	if len(certs) != len(keys) {
		return nil, fmt.Errorf("blob certificates (%d) not found for all blob keys (%d)", len(certs), len(keys))
	}

	certsMap := make(map[corev2.BlobKey]*corev2.BlobCertificate, len(certs))
	for _, cert := range certs {
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		certsMap[blobKey] = cert
	}

	// Keep the order of certs the same as the order of keys
	for i, key := range keys {
		c, ok := certsMap[key]
		if !ok {
			return nil, fmt.Errorf("blob certificate not found for blob key %s", key.Hex())
		}
		certs[i] = c
	}

	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{},
		ReferenceBlockNumber: referenceBlockNumber,
	}

	probe.SetStage("build_merkle_tree")
	tree, err := corev2.BuildMerkleTree(certs)
	if err != nil {
		return nil, fmt.Errorf("failed to build merkle tree: %w", err)
	}

	copy(batchHeader.BatchRoot[:], tree.Root())

	batchHeaderHash, err := batchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}

	probe.SetStage("put_batch_header")
	err = d.blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch header: %w", err)
	}

	probe.SetStage("put_batch")
	batch := &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: certs,
	}
	err = d.blobMetadataStore.PutBatch(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch: %w", err)
	}

	probe.SetStage("generate_proof")
	// accumulate inclusion infos in a map to avoid duplicate entries
	// batch write operation fails if there are duplicate entries
	inclusionInfoMap := make(map[corev2.BlobKey]*corev2.BlobInclusionInfo)
	for i, cert := range certs {
		if cert == nil || cert.BlobHeader == nil {
			return nil, fmt.Errorf("invalid blob certificate")
		}
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		merkleProof, err := tree.GenerateProofWithIndex(uint64(i), 0)
		if err != nil {
			return nil, fmt.Errorf("failed to generate merkle proof: %w", err)
		}

		inclusionInfoMap[blobKey] = &corev2.BlobInclusionInfo{
			BatchHeader:    batchHeader,
			BlobKey:        blobKey,
			BlobIndex:      uint32(i),
			InclusionProof: core.SerializeMerkleProof(merkleProof),
		}
	}

	probe.SetStage("put_inclusion_info")
	inclusionInfos := make([]*corev2.BlobInclusionInfo, len(inclusionInfoMap))
	i := 0
	for _, v := range inclusionInfoMap {
		inclusionInfos[i] = v
		i++
	}
	err = d.blobMetadataStore.PutBlobInclusionInfos(ctx, inclusionInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to put blob inclusion infos: %w", err)
	}

	d.cursor = cursor

	// Add blobs to the blob set to deduplicate blobs
	for _, blobKey := range keys {
		d.blobSet.AddBlob(blobKey)
	}

	batchSizeBytes := uint64(0)
	for _, blobKey := range keys {
		blobMetadata, ok := metadataMap[blobKey]
		if !ok {
			d.logger.Warn("missing blob metadata for blob key when updating signing metrics",
				"blobKey", blobKey.Hex(),
				"batchHeaderHash", batchHeaderHash)
			continue
		}
		batchSizeBytes += blobMetadata.BlobSize
	}

	d.logger.Debug("new batch", "referenceBlockNumber", referenceBlockNumber, "numBlobs", len(certs))
	return &batchData{
		Batch:           batch,
		BatchHeaderHash: batchHeaderHash,
		BlobKeys:        keys,
		Metadata:        metadataMap,
		OperatorState:   operatorState,
		BatchSizeBytes:  batchSizeBytes,
	}, nil
}

// Checks if a blob is older than MaxDispersalAge and handles it accordingly.
// If the blob is stale, it increments metrics, logs a warning, and updates the database status to Failed.
// Returns true if the blob is stale, otherwise false.
func (d *Dispatcher) checkAndHandleStaleBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	now time.Time,
	dispersalTimestamp int64,
) bool {
	dispersalTime := time.Unix(0, dispersalTimestamp)
	dispersalAge := now.Sub(dispersalTime)

	if dispersalAge <= d.MaxDispersalAge {
		return false
	}

	d.metrics.reportStaleDispersal()

	d.logger.Warnf(
		"discarding stale dispersal: blobKey=%s dispersalAge=%s maxAge=%s dispersalTime=%s",
		blobKey.Hex(),
		dispersalAge.String(),
		d.MaxDispersalAge.String(),
		dispersalTime.Format(time.RFC3339),
	)

	err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
	if err != nil {
		d.logger.Errorf("update stale blob status to Failed: blobKey=%s err=%w", blobKey.Hex(), err)
	} else {
		// Call beforeDispatch to clean up the blob from upstream encodingManager blobSet.
		// Since the stale check occurs before beforeDispatch would normally be called,
		// we must invoke it here to prevent orphaning the blob in the encoding manager's tracking.
		if d.beforeDispatch != nil {
			if err := d.beforeDispatch(blobKey); err != nil {
				d.logger.Errorf("beforeDispatch cleanup failed for stale blob: blobKey=%s err=%w", blobKey.Hex(), err)
			}
		}
		d.blobSet.RemoveBlob(blobKey)
	}

	return true
}

// GetOperatorState returns the operator state for the given quorums at the given block number
func (d *Dispatcher) GetOperatorState(
	ctx context.Context,
	metadatas []*v2.BlobMetadata,
	blockNumber uint64,
) (*core.IndexedOperatorState, error) {

	quorums := make(map[core.QuorumID]struct{}, 0)
	for _, m := range metadatas {
		for _, quorum := range m.BlobHeader.QuorumNumbers {
			quorums[quorum] = struct{}{}
		}
	}

	quorumIds := make([]core.QuorumID, len(quorums))
	i := 0
	for id := range quorums {
		quorumIds[i] = id
		i++
	}

	// GetIndexedOperatorState should return state for valid quorums only
	return d.chainState.GetIndexedOperatorState(ctx, uint(blockNumber), quorumIds)
}

func (d *Dispatcher) sendChunks(
	ctx context.Context,
	client clients.NodeClient,
	batch *corev2.Batch,
) (*core.Signature, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.AttestationTimeout)

	defer cancel()

	sig, err := client.StoreChunks(ctxWithTimeout, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to store chunks: %w", err)
	}

	return sig, nil
}

// updateBatchStatus updates the status of the blobs in the batch based on the quorum results
// If a blob is not included in the quorum results or runs into any unexpected errors, it is marked as failed
// If a blob is included in the quorum results, it is marked as complete
// This function also removes the blobs from the blob set indicating that this blob has been processed
// If the blob is removed from the blob set after the time it is retrieved as part of a batch
// for processing by `NewBatch` (when it's in `ENCODED` state) and before the time the batch
// is deduplicated against the blobSet, it will be dispatched again in a different batch.
func (d *Dispatcher) updateBatchStatus(
	ctx context.Context,
	batch *batchData,
	quorumResults map[core.QuorumID]uint8,
) error {

	var multierr error
	for i, cert := range batch.Batch.BlobCertificates {
		blobKey := batch.BlobKeys[i]
		if cert == nil || cert.BlobHeader == nil {
			d.logger.Error("invalid blob certificate in batch")
			err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
			if err != nil {
				multierr = multierror.Append(multierr,
					fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
			} else {
				d.blobSet.RemoveBlob(blobKey)
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Failed)
			}
			continue
		}

		failed := false
		for _, q := range cert.BlobHeader.QuorumNumbers {
			if res, ok := quorumResults[q]; !ok || res == 0 {
				d.logger.Warn("quorum result not found", "quorumID", q, "blobKey", blobKey.Hex())
				failed = true
				break
			}
		}

		if failed {
			err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
			if err != nil {
				multierr = multierror.Append(multierr,
					fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
			} else {
				d.blobSet.RemoveBlob(blobKey)
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Failed)
			}
			continue
		}

		err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Complete)
		if err != nil {
			multierr = multierror.Append(multierr,
				fmt.Errorf("failed to update blob status for blob %s to complete: %w", blobKey.Hex(), err))
		} else {
			d.blobSet.RemoveBlob(blobKey)
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			requestedAt := time.Unix(0, int64(metadata.RequestedAt))
			d.metrics.reportE2EDispersalLatency(time.Since(requestedAt))
			d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Complete)
		}
	}

	return multierr
}

func (d *Dispatcher) failBatch(ctx context.Context, batch *batchData) error {
	var multierr error
	for _, blobKey := range batch.BlobKeys {
		err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
		if err != nil {
			multierr = multierror.Append(multierr,
				fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Failed)
		}
		d.blobSet.RemoveBlob(blobKey)
	}

	return multierr
}
