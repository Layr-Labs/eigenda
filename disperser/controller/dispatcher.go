package controller

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
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
	// SignatureTickInterval is the interval at which new Attestations will be submitted to the blobMetadataStore,
	// as signature gathering progresses.
	SignatureTickInterval time.Duration
	NumRequestRetries     int
	// MaxBatchSize is the maximum number of blobs to dispatch in a batch
	MaxBatchSize int32
}

type Dispatcher struct {
	*DispatcherConfig

	blobMetadataStore *blobstore.BlobMetadataStore
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
	blobSet BlobSet
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
	blobMetadataStore *blobstore.BlobMetadataStore,
	pool common.WorkerPool,
	chainState core.IndexedChainState,
	aggregator core.SignatureAggregator,
	nodeClientManager NodeClientManager,
	logger logging.Logger,
	registry *prometheus.Registry,
	beforeDispatch func(blobKey corev2.BlobKey) error,
	blobSet BlobSet,
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
	return &Dispatcher{
		DispatcherConfig: config,

		blobMetadataStore: blobMetadataStore,
		pool:              pool,
		chainState:        chainState,
		aggregator:        aggregator,
		nodeClientManager: nodeClientManager,
		logger:            logger.With("component", "Dispatcher"),
		metrics:           newDispatcherMetrics(registry),

		cursor:         nil,
		beforeDispatch: beforeDispatch,
		blobSet:        blobSet,
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

				sigChan, batchData, err := d.HandleBatch(attestationCtx)
				if err != nil {
					if errors.Is(err, errNoBlobsToDispatch) {
						d.logger.Debug("no blobs to dispatch")
					} else {
						d.logger.Error("failed to process a batch", "err", err)
					}
					cancel()
					continue
				}
				go func() {
					err := d.HandleSignatures(ctx, attestationCtx, batchData, sigChan)
					if err != nil {
						d.logger.Error("failed to handle signatures", "err", err)
					}
					cancel()
				}()
			}
		}
	}()

	return nil

}

func (d *Dispatcher) HandleBatch(ctx context.Context) (chan core.SigningMessage, *batchData, error) {
	start := time.Now()
	defer func() {
		d.metrics.reportHandleBatchLatency(time.Since(start))
	}()

	currentBlockNumber, err := d.chainState.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	referenceBlockNumber := uint64(currentBlockNumber) - d.FinalizationBlockDelay

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob inclusion info for each blob in metadata store
	batchData, err := d.NewBatch(ctx, referenceBlockNumber)
	if err != nil {
		return nil, nil, err
	}

	batch := batchData.Batch
	state := batchData.OperatorState
	sigChan := make(chan core.SigningMessage, len(state.IndexedOperators))
	for opID, op := range state.IndexedOperators {
		opID := opID
		op := op
		host, _, _, v2DispersalPort, _, err := core.ParseOperatorSocket(op.Socket)
		if err != nil {
			d.logger.Warn("failed to parse operator socket, check if the socket format is correct",
				"operator", opID.Hex(),
				"socket", op.Socket,
				"err", err)
			sigChan <- core.SigningMessage{
				Signature:            nil,
				Operator:             opID,
				BatchHeaderHash:      batchData.BatchHeaderHash,
				AttestationLatencyMs: 0,
				Err:                  fmt.Errorf("failed to parse operator socket (%s): %w", op.Socket, err),
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
				Signature:            nil,
				Operator:             opID,
				BatchHeaderHash:      batchData.BatchHeaderHash,
				AttestationLatencyMs: 0,
				Err:                  err,
			}
			continue
		}

		submissionStart := time.Now()

		d.pool.Submit(func() {

			req := &corev2.DispersalRequest{
				OperatorID: opID,
				// TODO: get OperatorAddress
				OperatorAddress: gethcommon.Address{},
				Socket:          op.Socket,
				DispersedAt:     uint64(time.Now().UnixNano()),
				BatchHeader:     *batch.BatchHeader,
			}
			putDispersalRequestStart := time.Now()
			err := d.blobMetadataStore.PutDispersalRequest(ctx, req)
			if err != nil {
				d.logger.Error("failed to put dispersal request", "err", err)
				sigChan <- core.SigningMessage{
					Signature:            nil,
					Operator:             opID,
					BatchHeaderHash:      batchData.BatchHeaderHash,
					AttestationLatencyMs: 0,
					Err:                  err,
				}
				return
			}

			d.metrics.reportPutDispersalRequestLatency(time.Since(putDispersalRequestStart))

			var i int
			var lastErr error
			for i = 0; i < d.NumRequestRetries+1; i++ {
				sendChunksStart := time.Now()
				sig, err := d.sendChunks(ctx, client, batch)
				lastErr = err
				sendChunksFinished := time.Now()
				d.metrics.reportSendChunksLatency(sendChunksFinished.Sub(sendChunksStart))
				if err == nil {
					storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
						DispersalRequest: req,
						RespondedAt:      uint64(time.Now().UnixNano()),
						Signature:        sig.Bytes(),
						Error:            "",
					})
					if storeErr != nil {
						d.logger.Error("failed to store a succeeded dispersal response", "err", storeErr)
					}

					d.metrics.reportPutDispersalResponseLatency(time.Since(sendChunksFinished))

					sigChan <- core.SigningMessage{
						Signature:            sig,
						Operator:             opID,
						BatchHeaderHash:      batchData.BatchHeaderHash,
						AttestationLatencyMs: float64(time.Since(sendChunksStart)),
						Err:                  nil,
					}
					break
				}

				d.logger.Warn("failed to send chunks",
					"operator", opID.Hex(),
					"NumAttempts", i,
					"batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]),
					"err", err)
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
			}

			if lastErr != nil {
				d.logger.Warn("failed to send chunks",
					"operator", opID.Hex(),
					"NumAttempts", i,
					"batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]),
					"err", lastErr)
				storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
					DispersalRequest: req,
					RespondedAt:      uint64(time.Now().UnixNano()),
					Signature:        [32]byte{}, // all zero sig for failed dispersal
					Error:            lastErr.Error(),
				})
				if storeErr != nil {
					d.logger.Error("failed to store a failed dispersal response", "err", storeErr)
				}

				sigChan <- core.SigningMessage{
					Signature:            nil,
					Operator:             opID,
					BatchHeaderHash:      batchData.BatchHeaderHash,
					AttestationLatencyMs: 0,
					Err:                  lastErr,
				}
			}
			d.metrics.reportSendChunksRetryCount(float64(i))
		})

		d.metrics.reportPoolSubmissionLatency(time.Since(submissionStart))
	}

	return sigChan, batchData, nil
}

// HandleSignatures receives signatures from operators, validates, and aggregates them.
//
// This method submits Attestations to the blobMetadataStore, containing signing data from the SigningMessages received
// through the sigChan. It periodically submits Attestations, as signatures are gathered.
func (d *Dispatcher) HandleSignatures(
	ctx context.Context,
	attestationCtx context.Context,
	batchData *batchData,
	sigChan chan core.SigningMessage,
) error {
	if batchData == nil {
		return errors.New("batchData is required")
	}

	handleSignaturesStart := time.Now()
	defer func() {
		d.metrics.reportHandleSignaturesLatency(time.Since(handleSignaturesStart))
	}()

	batchHeaderHash := hex.EncodeToString(batchData.BatchHeaderHash[:])
	for _, key := range batchData.BlobKeys {
		err := d.blobMetadataStore.UpdateBlobStatus(ctx, key, v2.GatheringSignatures)
		if err != nil {
			d.logger.Error("update blob status to 'gathering signatures'",
				"blobKey", key.Hex(),
				"batchHeaderHash", batchHeaderHash,
				"err", err)
		}
	}

	// submit an empty attestation before starting to gather signatures.
	// a new attestation will be periodically resubmitted as signatures are gathered.
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
		// TODO: this used to cause the HandleSignatures method to fail entirely. Is it ok to continue trying here?
		d.logger.Error("error calling PutAttestation",
			"err", err,
			"batchHeaderHash", batchHeaderHash)
	}

	// This channel will remain open until the attestationTimeout triggers, or until signatures from all validators
	// have been received and processed. It will periodically yield QuorumAttestations with the latest set of received
	// signatures.
	attestationChan, err := core.ReceiveSignatures(
		attestationCtx,
		d.logger,
		batchData.OperatorState,
		batchData.BatchHeaderHash,
		sigChan,
		d.DispatcherConfig.SignatureTickInterval)
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

	// keep track of the final attestation, since that's the attestation which will determine to final batch status
	finalAttestation := &core.QuorumAttestation{}
	// continue receiving attestations from the channel until it's closed
	for receivedQuorumAttestation := range attestationChan {
		err := d.submitAttestation(ctx, batchData, receivedQuorumAttestation)
		if err != nil {
			d.logger.Warnf("error submitting attestation for batch %s: %v", batchHeaderHash, err)
			continue
		}

		finalAttestation = receivedQuorumAttestation
	}

	d.metrics.reportReceiveSignaturesLatency(time.Since(handleSignaturesStart))

	updateBatchStatusStartTime := time.Now()
	_, quorumPercentages := d.parseAndLogQuorumPercentages(batchHeaderHash, finalAttestation.QuorumResults)
	err = d.updateBatchStatus(ctx, batchData, quorumPercentages)
	if err != nil {
		return fmt.Errorf("update batch status: %w", err)
	}
	d.metrics.reportUpdateBatchStatusLatency(time.Since(updateBatchStatusStartTime))

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

// submitAttestation submits a QuorumAttestation to the blobMetadataStore
func (d *Dispatcher) submitAttestation(
	ctx context.Context,
	batchData *batchData,
	quorumAttestation *core.QuorumAttestation,
) error {
	sortedNonZeroQuorums, quorumPercentages := d.parseAndLogQuorumPercentages(
		hex.EncodeToString(batchData.BatchHeaderHash[:]),
		quorumAttestation.QuorumResults)
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
	if err != nil {
		return fmt.Errorf("error calling PutAttestation: %w", err)
	}
	d.metrics.reportPutAttestationLatency(time.Since(putAttestationStartTime))

	return nil
}

// parseAndLogQuorumPercentages iterates over the map of QuorumResults, and logs the signing percentages of each quorum.
//
// This method returns a sorted slice of nonZeroQuorums (quorums with >0 signing percentage), and a map from QuorumID to
// signing percentage.
func (d *Dispatcher) parseAndLogQuorumPercentages(
	batchHeaderHash string,
	quorumResults map[core.QuorumID]*core.QuorumResult,
) ([]core.QuorumID, map[core.QuorumID]uint8) {
	nonZeroQuorums := make([]core.QuorumID, 0)
	quorumPercentages := make(map[core.QuorumID]uint8)

	messageBuilder := strings.Builder{}
	messageBuilder.WriteString(fmt.Sprintf("batchHeaderHash: %s (quorumID, percentSigned)", batchHeaderHash))

	for quorumID, quorumResult := range quorumResults {
		messageBuilder.WriteString(fmt.Sprintf("\n%d, %d%%", quorumID, quorumResult.PercentSigned))

		if quorumResult.PercentSigned > 0 {
			nonZeroQuorums = append(nonZeroQuorums, quorumID)
			quorumPercentages[quorumID] = quorumResult.PercentSigned
		}
	}

	d.logger.Debug(messageBuilder.String())

	slices.Sort(nonZeroQuorums)

	return nonZeroQuorums, quorumPercentages
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
func (d *Dispatcher) NewBatch(ctx context.Context, referenceBlockNumber uint64) (*batchData, error) {
	newBatchStart := time.Now()
	defer func() {
		d.metrics.reportNewBatchLatency(time.Since(newBatchStart))
	}()
	blobMetadatas, cursor, err := d.blobMetadataStore.GetBlobMetadataByStatusPaginated(
		ctx,
		v2.Encoded,
		d.cursor,
		d.MaxBatchSize,
	)
	getBlobMetadataFinished := time.Now()
	d.metrics.reportGetBlobMetadataLatency(getBlobMetadataFinished.Sub(newBatchStart))
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

	state, err := d.GetOperatorState(ctx, blobMetadatas, referenceBlockNumber)
	getOperatorStateFinished := time.Now()
	d.metrics.reportGetOperatorStateLatency(getOperatorStateFinished.Sub(getBlobMetadataFinished))
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

	certs, _, err := d.blobMetadataStore.GetBlobCertificates(ctx, keys)
	getBlobCertificatesFinished := time.Now()
	d.metrics.reportGetBlobCertificatesLatency(getBlobCertificatesFinished.Sub(getOperatorStateFinished))
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

	tree, err := corev2.BuildMerkleTree(certs)
	if err != nil {
		return nil, fmt.Errorf("failed to build merkle tree: %w", err)
	}

	copy(batchHeader.BatchRoot[:], tree.Root())

	buildMerkleTreeFinished := time.Now()
	d.metrics.reportBuildMerkleTreeLatency(buildMerkleTreeFinished.Sub(getBlobCertificatesFinished))

	batchHeaderHash, err := batchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}

	err = d.blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	putBatchHeaderFinished := time.Now()
	d.metrics.reportPutBatchHeaderLatency(putBatchHeaderFinished.Sub(buildMerkleTreeFinished))
	if err != nil {
		return nil, fmt.Errorf("failed to put batch header: %w", err)
	}

	batch := &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: certs,
	}
	err = d.blobMetadataStore.PutBatch(ctx, batch)
	putBatchFinished := time.Now()
	d.metrics.reportPutBatchLatency(putBatchFinished.Sub(putBatchHeaderFinished))
	if err != nil {
		return nil, fmt.Errorf("failed to put batch: %w", err)
	}

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

	proofGenerationFinished := time.Now()
	d.metrics.reportProofLatency(proofGenerationFinished.Sub(putBatchFinished))

	inclusionInfos := make([]*corev2.BlobInclusionInfo, len(inclusionInfoMap))
	i := 0
	for _, v := range inclusionInfoMap {
		inclusionInfos[i] = v
		i++
	}
	err = d.blobMetadataStore.PutBlobInclusionInfos(ctx, inclusionInfos)
	putBlobInclusionInfosFinished := time.Now()
	d.metrics.reportPutInclusionInfosLatency(putBlobInclusionInfosFinished.Sub(proofGenerationFinished))
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
