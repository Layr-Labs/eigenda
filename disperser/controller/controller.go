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
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/signingrate"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
)

var errNoBlobsToDispatch = errors.New("no blobs to dispatch")

type BlobCallback func(blobKey corev2.BlobKey) error

type Controller struct {
	*ControllerConfig

	blobMetadataStore blobstore.MetadataStore
	pool              common.WorkerPool
	chainState        core.IndexedChainState
	aggregator        core.SignatureAggregator
	nodeClientManager NodeClientManager
	logger            logging.Logger
	metrics           *controllerMetrics
	getNow            func() time.Time

	// beforeDispatch function is called before dispatching a blob
	beforeDispatch BlobCallback

	controllerLivenessChan chan<- healthcheck.HeartbeatMessage

	// A utility responsible for fetching batch metadata (i.e. reference block number and operator state).
	batchMetadataManager metadata.BatchMetadataManager

	// Tracks signing rates for validators and serves queries about signing rates.
	signingRateTracker signingrate.SigningRateTracker

	// Acquires blobs ready for dispersal from the encoder->controller pipeline.
	blobDispersalQueue BlobDispersalQueue
}

type batchData struct {
	Batch           *corev2.Batch
	BatchHeaderHash [32]byte
	BlobKeys        []corev2.BlobKey
	Metadata        map[corev2.BlobKey]*v2.BlobMetadata
	OperatorState   *core.IndexedOperatorState
	BatchSizeBytes  uint64
}

func NewController(
	ctx context.Context,
	config *ControllerConfig,
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
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage,
	signingRateTracker signingrate.SigningRateTracker,
	userAccountRemapping map[string]string,
	validatorIdRemapping map[string]string,
) (*Controller, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	metrics, err := newControllerMetrics(
		registry,
		config.SignificantSigningThresholdFraction,
		config.CollectDetailedValidatorSigningMetrics,
		config.EnablePerAccountBlobStatusMetrics,
		userAccountRemapping,
		validatorIdRemapping)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %v", err)
	}

	blobDispersalQueue, err := NewDynamodbBlobDispersalQueue(
		ctx,
		logger,
		blobMetadataStore,
		config.BlobDispersalQueueSize,
		config.BlobDispersalRequestBatchSize,
		config.BlobDispersalRequestBackoffPeriod,
		10*time.Minute, // TODO flags
		10*time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("NewDynamodbBlobDispersalQueue: %w", err)
	}

	return &Controller{
		ControllerConfig:       config,
		blobMetadataStore:      blobMetadataStore,
		pool:                   pool,
		chainState:             chainState,
		aggregator:             aggregator,
		nodeClientManager:      nodeClientManager,
		logger:                 logger.With("component", "controller"),
		metrics:                metrics,
		getNow:                 getNow,
		beforeDispatch:         beforeDispatch,
		controllerLivenessChan: controllerLivenessChan,
		batchMetadataManager:   batchMetadataManager,
		signingRateTracker:     signingRateTracker,
		blobDispersalQueue:     blobDispersalQueue,
	}, nil
}

func (c *Controller) Start(ctx context.Context) error {
	err := c.chainState.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start chain state: %w", err)
	}

	go func() {
		ticker := time.NewTicker(c.PullInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				attestationCtx, cancel := context.WithTimeout(ctx, c.BatchAttestationTimeout)
				probe := c.metrics.newBatchProbe()

				sigChan, batchData, err := c.HandleBatch(attestationCtx, probe)
				if err != nil {
					if errors.Is(err, errNoBlobsToDispatch) {
						c.logger.Debug("no blobs to dispatch")
					} else {
						c.logger.Error("failed to process a batch", "err", err)
					}
					cancel()
					probe.End()
					continue
				}
				go func() {
					probe.SetStage("handle_signatures")
					err := c.HandleSignatures(ctx, attestationCtx, batchData, sigChan)
					if err != nil {
						c.logger.Error("failed to handle signatures", "err", err)
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
func (c *Controller) HandleBatch(
	ctx context.Context,
	batchProbe *common.SequenceProbe,
) (chan core.SigningMessage, *batchData, error) {
	// Signal Liveness to indicate no stall
	healthcheck.SignalHeartbeat(c.logger, "dispatcher", c.controllerLivenessChan)

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob inclusion info for each blob in metadata store
	batchData, err := c.NewBatch(ctx, batchProbe)
	if err != nil {
		return nil, nil, err
	}

	batchProbe.SetStage("send_requests")

	signingResponseChan := make(chan core.SigningMessage, len(batchData.OperatorState.IndexedOperators))
	for validatorId, validatorInfo := range batchData.OperatorState.IndexedOperators {

		validatorProbe := c.metrics.newSendToValidatorProbe()
		validatorProbe.SetStage("pool_submission")

		c.pool.Submit(func() {
			signature, latency, err := c.sendChunksToValidator(
				ctx,
				batchData,
				validatorId,
				validatorInfo,
				validatorProbe)

			if err != nil {
				c.logger.Warn("error sending chunks to validator",
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
func (c *Controller) sendChunksToValidator(
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

	client, err := c.nodeClientManager.GetClient(host, v2DispersalPort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get node client for validator at host %s port %s: %w",
			host, v2DispersalPort, err)
	}

	validatorProbe.SetStage("put_dispersal_request")

	req := &corev2.DispersalRequest{
		OperatorID:  validatorId,
		Socket:      validatorInfo.Socket,
		DispersedAt: uint64(time.Now().UnixNano()),
		BatchHeader: *batchData.Batch.BatchHeader,
	}
	err = c.blobMetadataStore.PutDispersalRequest(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to put dispersal request for validator: %w", err)
	}

	validatorProbe.SetStage("send_chunks")

	start := time.Now()

	sig, err := c.sendChunks(ctx, client, batchData.Batch)
	if err != nil {
		storeErr := c.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
			DispersalRequest: req,
			RespondedAt:      uint64(time.Now().UnixNano()),
			Signature:        [32]byte{}, // all zero sig for failed dispersal
			Error:            err.Error(),
		})
		if storeErr != nil {
			c.logger.Error("failed to store a failed dispersal response", "err", storeErr)
		}
		return nil, 0, fmt.Errorf("failed to send chunks to validator: %w", err)
	}

	latency = time.Since(start)

	validatorProbe.SetStage("put_dispersal_response")
	storeErr := c.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
		DispersalRequest: req,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        sig.Bytes(),
	})
	if storeErr != nil {
		c.logger.Error("failed to store a succeeded dispersal response", "err", storeErr)
	}

	return sig, latency, nil
}

// HandleSignatures receives SigningMessages from operators for a given batch through the input sigChan. The signatures
// are validated, aggregated, and used to put an Attestation for the batch into the blobMetadataStore. The Attestation
// is periodically updated as additional signatures are gathered.
//
// This method will continue gathering signatures until a SigningMessage has been received from every operator, or until
// the global attestationCtx times out.
func (c *Controller) HandleSignatures(
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
		err := c.updateBlobStatus(ctx, key, v2.GatheringSignatures)
		if err != nil {
			c.logger.Error("failed to update blob status to 'gathering signatures'",
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
	err := c.blobMetadataStore.PutAttestation(ctx, attestation)
	if err != nil {
		// this error isn't fatal: a subsequent PutAttestation attempt might succeed
		c.logger.Error("error calling PutAttestation",
			"err", err,
			"batchHeaderHash", batchHeaderHash)
	}

	// This channel will remain open until the attestationTimeout triggers, or until signatures from all validators
	// have been received and processed. It will periodically yield QuorumAttestations with the latest set of received
	// signatures.
	attestationChan, err := ReceiveSignatures(
		attestationCtx,
		c.logger,
		c.metrics,
		c.signingRateTracker,
		batchData.OperatorState,
		batchData.BatchHeaderHash,
		sigChan,
		c.ControllerConfig.SignatureTickInterval,
		c.ControllerConfig.SignificantSigningThresholdFraction,
		batchData.BatchSizeBytes)
	if err != nil {
		receiveSignaturesErr := fmt.Errorf("receive and validate signatures for batch %s: %w", batchHeaderHash, err)

		dbErr := c.failBatch(ctx, batchData)
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
		err := c.updateAttestation(ctx, batchData, receivedQuorumAttestation)
		if err != nil {
			c.logger.Warnf("error updating attestation for batch %s: %v", batchHeaderHash, err)
			continue
		}

		finalAttestation = receivedQuorumAttestation
	}

	updateBatchStatusStartTime := time.Now()
	_, quorumPercentages := c.parseQuorumPercentages(finalAttestation.QuorumResults)
	err = c.updateBatchStatus(ctx, batchData, quorumPercentages)
	c.metrics.reportUpdateBatchStatusLatency(time.Since(updateBatchStatusStartTime))
	if err != nil {
		return fmt.Errorf("update batch status: %w", err)
	}

	return nil
}

// updateAttestation updates the QuorumAttestation in the blobMetadataStore
func (c *Controller) updateAttestation(
	ctx context.Context,
	batchData *batchData,
	quorumAttestation *core.QuorumAttestation,
) error {
	sortedNonZeroQuorums, quorumPercentages := c.parseQuorumPercentages(quorumAttestation.QuorumResults)
	if len(sortedNonZeroQuorums) == 0 {
		return errors.New("all quorums received no attestation for batch")
	}

	aggregationStartTime := time.Now()
	signatureAggregation, err := c.aggregator.AggregateSignatures(
		batchData.OperatorState,
		quorumAttestation,
		sortedNonZeroQuorums)
	c.metrics.reportAggregateSignaturesLatency(time.Since(aggregationStartTime))
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
	err = c.blobMetadataStore.PutAttestation(ctx, attestation)
	c.metrics.reportPutAttestationLatency(time.Since(putAttestationStartTime))
	if err != nil {
		return fmt.Errorf("put attestation: %w", err)
	}

	c.logAttestationUpdate(hex.EncodeToString(batchData.BatchHeaderHash[:]), quorumAttestation.QuorumResults)

	return nil
}

// parseQuorumPercentages iterates over the map of QuorumResults, and returns a sorted slice of nonZeroQuorums
// (quorums with >0 signing percentage), and a map from QuorumID to signing percentage.
func (c *Controller) parseQuorumPercentages(
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
func (c *Controller) logAttestationUpdate(batchHeaderHash string, quorumResults map[core.QuorumID]*core.QuorumResult) {
	quorumPercentagesBuilder := strings.Builder{}
	quorumPercentagesBuilder.WriteString("(")

	for quorumID, quorumResult := range quorumResults {
		quorumPercentagesBuilder.WriteString(
			fmt.Sprintf("quorum_%d: %d%%, ", quorumID, quorumResult.PercentSigned))
	}
	quorumPercentagesBuilder.WriteString(")")

	c.logger.Debug("attestation updated",
		"batchHeaderHash", batchHeaderHash,
		"quorumPercentages", quorumPercentagesBuilder.String())
}

// NewBatch creates a batch of blobs to dispatch
// Warning: This function is not thread-safe
func (c *Controller) NewBatch(
	ctx context.Context,
	probe *common.SequenceProbe,
) (*batchData, error) {

	batchMetadata := c.batchMetadataManager.GetMetadata()
	referenceBlockNumber := batchMetadata.ReferenceBlockNumber()
	operatorState := batchMetadata.OperatorState()

	probe.SetStage("get_blob_metadata")

	blobMetadatas := make([]*v2.BlobMetadata, 0, c.MaxBatchSize)
	keepLooking := true
	for keepLooking && int32(len(blobMetadatas)) < c.MaxBatchSize {
		var next *v2.BlobMetadata
		select {
		case next = <-c.blobDispersalQueue.GetBlobChannel():
		default:
			// no more blobs available right now
			keepLooking = false
		}

		blobKey, err := next.BlobHeader.BlobKey()
		if err != nil {
			c.logger.Errorf("failed to compute blob key for fetched blob, skipping: %v", err)
			continue
		}

		if c.checkAndHandleStaleBlob(
			ctx,
			blobKey,
			c.getNow(),
			next.BlobHeader.PaymentMetadata.Timestamp) {

			// discard stale blob
			continue
		}

		blobMetadatas = append(blobMetadatas, next)
	}

	// If we fail to finish batch creation, we need to go back and ensure that we mark all of the blobs
	// that were about to be in the batch as having failed.
	batchCreationSuccessful := false
	defer func() {
		if !batchCreationSuccessful {
			c.logger.Warnf("batch creation failed, marking %d blobs as failed", len(blobMetadatas))
			c.markBatchAsFailed(ctx, blobMetadatas)
		}
	}()

	if len(blobMetadatas) == 0 {
		return nil, errNoBlobsToDispatch
	}
	c.logger.Debug("got new metadatas to make batch",
		"numBlobs", len(blobMetadatas),
		"referenceBlockNumber", referenceBlockNumber)

	keys := make([]corev2.BlobKey, len(blobMetadatas))
	metadataMap := make(map[corev2.BlobKey]*v2.BlobMetadata, len(blobMetadatas))
	for i, metadata := range blobMetadatas {
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}
		keys[i] = blobKey
		metadataMap[blobKey] = metadata

		if c.beforeDispatch != nil {
			err = c.beforeDispatch(blobKey)
			if err != nil {
				c.logger.Error("beforeDispatch function failed", "blobKey", blobKey.Hex(), "err", err)
			}
		}
	}

	probe.SetStage("get_blob_certs")
	certs, _, err := c.blobMetadataStore.GetBlobCertificates(ctx, keys)
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
	err = c.blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch header: %w", err)
	}

	probe.SetStage("put_batch")
	batch := &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: certs,
	}
	err = c.blobMetadataStore.PutBatch(ctx, batch)
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
	err = c.blobMetadataStore.PutBlobInclusionInfos(ctx, inclusionInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to put blob inclusion infos: %w", err)
	}

	batchSizeBytes := uint64(0)
	for _, blobKey := range keys {
		blobMetadata, ok := metadataMap[blobKey]
		if !ok {
			c.logger.Warn("missing blob metadata for blob key when updating signing metrics",
				"blobKey", blobKey.Hex(),
				"batchHeaderHash", batchHeaderHash)
			continue
		}
		batchSizeBytes += blobMetadata.BlobSize
	}

	c.logger.Debug("new batch", "referenceBlockNumber", referenceBlockNumber, "numBlobs", len(certs))
	batchCreationSuccessful = true
	return &batchData{
		Batch:           batch,
		BatchHeaderHash: batchHeaderHash,
		BlobKeys:        keys,
		Metadata:        metadataMap,
		OperatorState:   operatorState,
		BatchSizeBytes:  batchSizeBytes,
	}, nil
}

// If when creating a batch we encounter a failure, we need to mark each blob that was planned to be a part of that
// batch as Failed.
func (c *Controller) markBatchAsFailed(
	ctx context.Context,
	blobsInBatch []*v2.BlobMetadata,
) {
	for _, blobMetadata := range blobsInBatch {
		blobKey, err := blobMetadata.BlobHeader.BlobKey()
		if err != nil {
			c.logger.Errorf("compute blob key: %w", err)
			continue
		}

		err = c.updateBlobStatus(ctx, blobKey, v2.Failed)
		if err != nil {
			c.logger.Errorf("update blob status to failed: %w", err)
		}
	}
}

// Checks if a blob is older than MaxDispersalAge and handles it accordingly.
// If the blob is stale, it increments metrics, logs a warning, and updates the database status to Failed.
// Returns true if the blob is stale, otherwise false.
func (c *Controller) checkAndHandleStaleBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	now time.Time,
	dispersalTimestamp int64,
) bool {
	dispersalTime := time.Unix(0, dispersalTimestamp)
	dispersalAge := now.Sub(dispersalTime)

	if dispersalAge <= c.MaxDispersalAge {
		return false
	}

	c.metrics.reportStaleDispersal()

	c.logger.Warnf(
		"discarding stale dispersal: blobKey=%s dispersalAge=%s maxAge=%s dispersalTime=%s",
		blobKey.Hex(),
		dispersalAge.String(),
		c.MaxDispersalAge.String(),
		dispersalTime.Format(time.RFC3339),
	)

	err := c.updateBlobStatus(ctx, blobKey, v2.Failed)
	if err != nil {
		c.logger.Errorf("update blob status: %w", err)
	} else {
		// Call beforeDispatch to clean up the blob from upstream encodingManager blobSet.
		// Since the stale check occurs before beforeDispatch would normally be called,
		// we must invoke it here to prevent orphaning the blob in the encoding manager's tracking.
		if c.beforeDispatch != nil {
			if err := c.beforeDispatch(blobKey); err != nil {
				c.logger.Errorf("beforeDispatch cleanup failed for stale blob: blobKey=%s err=%w", blobKey.Hex(), err)
			}
		}
	}

	return true
}

func (c *Controller) sendChunks(
	ctx context.Context,
	client clients.NodeClient,
	batch *corev2.Batch,
) (*core.Signature, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.AttestationTimeout)

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
func (c *Controller) updateBatchStatus(
	ctx context.Context,
	batch *batchData,
	quorumResults map[core.QuorumID]uint8,
) error {

	var multierr error
	for i, cert := range batch.Batch.BlobCertificates {
		blobKey := batch.BlobKeys[i]
		if cert == nil || cert.BlobHeader == nil {
			c.logger.Error("invalid blob certificate in batch")
			err := c.updateBlobStatus(ctx, blobKey, v2.Failed)
			if err != nil {
				multierr = multierror.Append(multierr, fmt.Errorf("update blob status: %w", err))
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				c.metrics.reportCompletedBlob(
					int(metadata.BlobSize), v2.Failed, metadata.BlobHeader.PaymentMetadata.AccountID.Hex())
			}
			continue
		}

		failed := false
		for _, q := range cert.BlobHeader.QuorumNumbers {
			if res, ok := quorumResults[q]; !ok || res == 0 {
				c.logger.Warn("quorum result not found", "quorumID", q, "blobKey", blobKey.Hex())
				failed = true
				break
			}
		}

		if failed {
			err := c.updateBlobStatus(ctx, blobKey, v2.Failed)
			if err != nil {
				multierr = multierror.Append(multierr, fmt.Errorf("update blob status: %w", err))
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				c.metrics.reportCompletedBlob(
					int(metadata.BlobSize), v2.Failed, metadata.BlobHeader.PaymentMetadata.AccountID.Hex())
			}
			continue
		}

		err := c.updateBlobStatus(ctx, blobKey, v2.Complete)

		if err != nil {
			multierr = multierror.Append(multierr, fmt.Errorf("update blob status: %w", err))
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			requestedAt := time.Unix(0, int64(metadata.RequestedAt))
			c.metrics.reportE2EDispersalLatency(time.Since(requestedAt))
			c.metrics.reportCompletedBlob(
				int(metadata.BlobSize), v2.Complete, metadata.BlobHeader.PaymentMetadata.AccountID.Hex())
		}
	}

	return multierr
}

func (c *Controller) failBatch(ctx context.Context, batch *batchData) error {
	var multierr error
	for _, blobKey := range batch.BlobKeys {
		err := c.updateBlobStatus(ctx, blobKey, v2.Failed)
		if err != nil {
			multierr = multierror.Append(multierr,
				fmt.Errorf("update blob status: %w", err))
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			c.metrics.reportCompletedBlob(
				int(metadata.BlobSize), v2.Failed, metadata.BlobHeader.PaymentMetadata.AccountID.Hex())
		}
	}

	return multierr
}

// Update the blob status. If the status is terminal, remove the blob from the blob set.
func (c *Controller) updateBlobStatus(ctx context.Context, blobKey corev2.BlobKey, status v2.BlobStatus) error {
	err := c.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, status)
	if err != nil {
		return fmt.Errorf("failed to update blob status for blob %s to %s: %w", blobKey.Hex(), status.String(), err)
	}

	return nil
}
