package controller

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
)

var errNoBlobsToDispatch = errors.New("no blobs to dispatch")

type BlobCallback func(blobKey corev2.BlobKey) error

type DispatcherConfig struct {
	PullInterval time.Duration

	FinalizationBlockDelay uint64
	// The maximum time permitted to wait for a node to provide a signature for a batch.
	AttestationTimeout time.Duration
	// The maximum time permitted to wait for all nodes to provide signatures for a batch.
	BatchAttestationTimeout time.Duration
	// SignatureTickInterval is the interval at which Attestations will be updated in the blobMetadataStore,
	// as signature gathering progresses.
	SignatureTickInterval time.Duration
	NumRequestRetries     int
	// MaxBatchSize is the maximum number of blobs to dispatch in a batch
	MaxBatchSize int32
	// SignificantSigningThresholdPercentage is a configurable "important" signing threshold. Right now, it's being
	// used to track signing metrics, to understand system performance. If the value is 0, then special handling for
	// the threshold is disabled.
	SignificantSigningThresholdPercentage uint8
	// Important signing thresholds for metrics reporting.
	// Values should be between 0.0 (0% signed) and 1.0 (100% signed).
	SignificantSigningMetricsThresholds []string
}

type Dispatcher struct {
	*DispatcherConfig

	blobMetadataStore blobstore.MetadataStore
	pool              common.WorkerPool
	chainState        core.IndexedChainState
	aggregator        core.SignatureAggregator
	nodeClientManager NodeClientManager
	logger            logging.Logger
	metrics           *dispatcherMetrics

	cursor *blobstore.StatusIndexCursor
	// beforeDispatch function is called before dispatching a blob
	beforeDispatch BlobCallback
	// blobSet keeps track of blobs that are being dispatched
	// This is used to deduplicate blobs to prevent the same blob from being dispatched multiple times
	// Blobs are removed from the queue when they are in a terminal state (Complete or Failed)
	blobSet                BlobSet
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage
}

type batchData struct {
	Batch           *corev2.Batch
	BatchHeaderHash [32]byte
	BlobKeys        []corev2.BlobKey
	Metadata        map[corev2.BlobKey]*v2.BlobMetadata
	OperatorState   *core.IndexedOperatorState
}

func NewDispatcher(
	config *DispatcherConfig,
	blobMetadataStore blobstore.MetadataStore,
	pool common.WorkerPool,
	chainState core.IndexedChainState,
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
	if config.PullInterval == 0 ||
		config.AttestationTimeout == 0 ||
		config.BatchAttestationTimeout == 0 ||
		config.SignatureTickInterval == 0 ||
		config.MaxBatchSize == 0 {
		return nil, errors.New("invalid config")
	}

	// CLI library doesn't support float slices at current version, parsing must happen manually
	significantThresholds := make([]float64, 0, len(config.SignificantSigningMetricsThresholds))
	for _, threshold := range config.SignificantSigningMetricsThresholds {
		significantThreshold, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse significant threshold %s: %v", threshold, err)
		}
		significantThresholds = append(significantThresholds, significantThreshold)
	}

	metrics, err := newDispatcherMetrics(registry, significantThresholds)
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

		cursor:                 nil,
		beforeDispatch:         beforeDispatch,
		blobSet:                blobSet,
		controllerLivenessChan: controllerLivenessChan,
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

func (d *Dispatcher) HandleBatch(
	ctx context.Context,
	batchProbe *common.SequenceProbe,
) (chan core.SigningMessage, *batchData, error) {
	// Signal Liveness to indicate no stall
	healthcheck.SignalHeartbeat("dispatcher", d.controllerLivenessChan, d.logger)

	batchProbe.SetStage("get_reference_block")
	currentBlockNumber, err := d.chainState.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	referenceBlockNumber := uint64(currentBlockNumber) - d.FinalizationBlockDelay

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob inclusion info for each blob in metadata store
	batchData, err := d.NewBatch(ctx, referenceBlockNumber, batchProbe)
	if err != nil {
		return nil, nil, err
	}

	batchProbe.SetStage("send_requests")

	batch := batchData.Batch
	state := batchData.OperatorState
	sigChan := make(chan core.SigningMessage, len(state.IndexedOperators))
	for opID, op := range state.IndexedOperators {

		validatorProbe := d.metrics.newSendToValidatorProbe()
		validatorProbe.SetStage("get_client")

		opID := opID
		op := op
		host, _, _, v2DispersalPort, _, err := core.ParseOperatorSocket(op.Socket)
		if err != nil {
			d.logger.Warn("failed to parse operator socket, check if the socket format is correct",
				"operator", opID.Hex(),
				"socket", op.Socket,
				"err", err)
			sigChan <- core.SigningMessage{
				Signature:       nil,
				Operator:        opID,
				BatchHeaderHash: batchData.BatchHeaderHash,
				TimeReceived:    time.Now(),
				Err:             fmt.Errorf("failed to parse operator socket (%s): %w", op.Socket, err),
			}
			continue
		}

		client, err := d.nodeClientManager.GetClient(host, v2DispersalPort)
		if err != nil {
			d.logger.Warn("failed to get node client; node may not be reachable",
				"operator", opID.Hex(),
				"host", host,
				"v2DispersalPort", v2DispersalPort,
				"err", err)
			sigChan <- core.SigningMessage{
				Signature:       nil,
				Operator:        opID,
				BatchHeaderHash: batchData.BatchHeaderHash,
				TimeReceived:    time.Now(),
				Err:             err,
			}
			continue
		}

		validatorProbe.SetStage("pool_submission")

		d.pool.Submit(func() {
			defer validatorProbe.End()
			validatorProbe.SetStage("put_dispersal_request")

			req := &corev2.DispersalRequest{
				OperatorID: opID,
				// TODO: get OperatorAddress
				OperatorAddress: gethcommon.Address{},
				Socket:          op.Socket,
				DispersedAt:     uint64(time.Now().UnixNano()),
				BatchHeader:     *batch.BatchHeader,
			}
			err := d.blobMetadataStore.PutDispersalRequest(ctx, req)
			if err != nil {
				d.logger.Error("failed to put dispersal request", "err", err)
				sigChan <- core.SigningMessage{
					Signature:       nil,
					Operator:        opID,
					BatchHeaderHash: batchData.BatchHeaderHash,
					TimeReceived:    time.Now(),
					Err:             err,
				}
				return
			}

			var i int
			var lastErr error
			for i = 0; i < d.NumRequestRetries+1; i++ {
				validatorProbe.SetStage("send_chunks")

				sig, err := d.sendChunks(ctx, client, batch)
				lastErr = err
				if err == nil {
					validatorProbe.SetStage("put_dispersal_response")
					storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
						DispersalRequest: req,
						RespondedAt:      uint64(time.Now().UnixNano()),
						Signature:        sig.Bytes(),
						Error:            "",
					})
					if storeErr != nil {
						d.logger.Error("failed to store a succeeded dispersal response", "err", storeErr)
					}

					sigChan <- core.SigningMessage{
						Signature:       sig,
						Operator:        opID,
						BatchHeaderHash: batchData.BatchHeaderHash,
						TimeReceived:    time.Now(),
						Err:             nil,
					}
					break
				}

				// Parse batch meterer error for metrics collection
				if bmErr, parsed := ParseBatchMeterError(err.Error()); parsed {
					LogBatchMeterError(d.logger, bmErr, err)

					// Report metrics for batch meterer errors
					if d.metrics != nil {
						category := GetBatchMeterErrorCategory(bmErr)
						willRetry := i < d.NumRequestRetries
						d.metrics.reportBatchMeterError(bmErr.Code, category, willRetry)
					}
				} else {
					d.logger.Warn("failed to send chunks",
						"operator", opID.Hex(),
						"NumAttempts", i,
						"batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]),
						"err", err)
				}
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
			}

			if lastErr != nil {
				// Enhanced error logging with batch meterer error details
				errorStr := lastErr.Error()
				if bmErr, parsed := ParseBatchMeterError(errorStr); parsed {
					summary := GetBatchMeterErrorSummary(bmErr)
					d.logger.Error("batch meterer validation failed after retries",
						"operator", opID.Hex(),
						"NumAttempts", i,
						"batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]),
						"errorCode", bmErr.Code,
						"accountID", bmErr.AccountID,
						"quorumID", bmErr.QuorumID,
						"summary", summary,
						"originalError", lastErr)
					// Store the structured error summary for better analysis
					errorStr = fmt.Sprintf("[%s] %s", bmErr.Code, summary)
				} else {
					d.logger.Warn("failed to send chunks",
						"operator", opID.Hex(),
						"NumAttempts", i,
						"batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]),
						"err", lastErr)
				}

				storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
					DispersalRequest: req,
					RespondedAt:      uint64(time.Now().UnixNano()),
					Signature:        [32]byte{}, // all zero sig for failed dispersal
					Error:            errorStr,
				})
				if storeErr != nil {
					d.logger.Error("failed to store a failed dispersal response", "err", storeErr)
				}

				sigChan <- core.SigningMessage{
					Signature:       nil,
					Operator:        opID,
					BatchHeaderHash: batchData.BatchHeaderHash,
					TimeReceived:    time.Now(),
					Err:             lastErr,
				}
			}
			d.metrics.reportSendChunksRetryCount(float64(i))
		})
	}

	batchProbe.SetStage("await_responses")

	return sigChan, batchData, nil
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
		d.DispatcherConfig.SignificantSigningThresholdPercentage)
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

	// Track attestation metrics
	operatorCount := make(map[core.QuorumID]int)
	signerCount := make(map[core.QuorumID]int)
	for quorumID, opState := range batchData.OperatorState.Operators {
		operatorCount[quorumID] = len(opState)
		if _, ok := signerCount[quorumID]; !ok {
			signerCount[quorumID] = 0
		}
		for opID := range opState {
			if _, ok := finalAttestation.SignerMap[opID]; ok {
				signerCount[quorumID]++
			}
		}
	}
	d.metrics.reportAttestation(operatorCount, signerCount, finalAttestation.QuorumResults)

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
		ctx,
		d.chainState,
		uint(batchData.Batch.BatchHeader.ReferenceBlockNumber),
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

func (d *Dispatcher) dedupBlobs(blobs []*v2.BlobMetadata) []*v2.BlobMetadata {
	dedupedBlobs := make([]*v2.BlobMetadata, 0)
	for _, blob := range blobs {
		key, err := blob.BlobHeader.BlobKey()
		if err != nil {
			d.logger.Error("failed to get blob key", "err", err, "requestedAt", blob.RequestedAt)
			continue
		}
		if !d.blobSet.Contains(key) {
			dedupedBlobs = append(dedupedBlobs, blob)
		}
	}
	return dedupedBlobs
}

// NewBatch creates a batch of blobs to dispatch
// Warning: This function is not thread-safe
func (d *Dispatcher) NewBatch(
	ctx context.Context,
	referenceBlockNumber uint64,
	probe *common.SequenceProbe,
) (*batchData, error) {

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

	blobMetadatas = d.dedupBlobs(blobMetadatas)
	d.metrics.reportBlobSetSize(d.blobSet.Size())
	if len(blobMetadatas) == 0 {
		return nil, errNoBlobsToDispatch
	}
	d.logger.Debug("got new metadatas to make batch",
		"numBlobs", len(blobMetadatas),
		"referenceBlockNumber", referenceBlockNumber)

	probe.SetStage("get_operator_state")
	state, err := d.GetOperatorState(ctx, blobMetadatas, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state at block %d: %w", referenceBlockNumber, err)
	}

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

	d.logger.Debug("new batch", "referenceBlockNumber", referenceBlockNumber, "numBlobs", len(certs))
	return &batchData{
		Batch:           batch,
		BatchHeaderHash: batchHeaderHash,
		BlobKeys:        keys,
		Metadata:        metadataMap,
		OperatorState:   state,
	}, nil
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
