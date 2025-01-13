package controller

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
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

type DispatcherConfig struct {
	PullInterval time.Duration

	FinalizationBlockDelay uint64
	NodeRequestTimeout     time.Duration
	NumRequestRetries      int
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
) (*Dispatcher, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}
	if config.PullInterval == 0 || config.NodeRequestTimeout == 0 || config.MaxBatchSize == 0 {
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

		cursor: nil,
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
				sigChan, batchData, err := d.HandleBatch(ctx)
				if err != nil {
					if errors.Is(err, errNoBlobsToDispatch) {
						d.logger.Debug("no blobs to dispatch")
					} else {
						d.logger.Error("failed to process a batch", "err", err)
					}
					continue
				}
				go func() {
					err := d.HandleSignatures(ctx, batchData, sigChan)
					if err != nil {
						d.logger.Error("failed to handle signatures", "err", err)
					}
					close(sigChan)
					// TODO(ian-shim): handle errors and mark failed
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

	currentBlockNumber, err := d.chainState.GetCurrentBlockNumber()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	referenceBlockNumber := uint64(currentBlockNumber) - d.FinalizationBlockDelay

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob verification info for each blob in metadata store
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
		host, dispersalPort, _, err := core.ParseOperatorSocket(op.Socket)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse operator socket: %w", err)
		}

		client, err := d.nodeClientManager.GetClient(host, dispersalPort)
		if err != nil {
			d.logger.Error("failed to get node client", "operator", opID.Hex(), "err", err)
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
						d.logger.Error("failed to put dispersal response", "err", storeErr)
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

				d.logger.Warn("failed to send chunks", "operator", opID.Hex(), "NumAttempts", i, "batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]), "err", err)
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
			}

			if lastErr != nil {
				d.logger.Error("failed to send chunks", "operator", opID.Hex(), "NumAttempts", i, "batchHeader", hex.EncodeToString(batchData.BatchHeaderHash[:]), "err", lastErr)
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

// HandleSignatures receives signatures from operators, validates, and aggregates them
func (d *Dispatcher) HandleSignatures(ctx context.Context, batchData *batchData, sigChan chan core.SigningMessage) error {
	if batchData == nil {
		return errors.New("batchData is required")
	}
	handleSignaturesStart := time.Now()
	defer func() {
		d.metrics.reportHandleSignaturesLatency(time.Since(handleSignaturesStart))
	}()

	batchHeaderHash := hex.EncodeToString(batchData.BatchHeaderHash[:])
	quorumAttestation, err := d.aggregator.ReceiveSignatures(ctx, batchData.OperatorState, batchData.BatchHeaderHash, sigChan)
	if err != nil {
		dbErr := d.failBatch(ctx, batchData)
		if dbErr != nil {
			return fmt.Errorf("failed to update blob statuses for batch %s to failed: %w", batchHeaderHash, dbErr)
		}
		return fmt.Errorf("failed to receive and validate signatures for batch %s: %w", batchHeaderHash, err)
	}
	receiveSignaturesFinished := time.Now()
	d.metrics.reportReceiveSignaturesLatency(receiveSignaturesFinished.Sub(handleSignaturesStart))

	nonZeroQuorums := make([]core.QuorumID, 0)
	quorumResults := make(map[core.QuorumID]uint8)
	for quorumID, quorumResult := range quorumAttestation.QuorumResults {
		d.logger.Debug("quorum attestation results", "quorumID", quorumID, "result", quorumResult)
		if quorumResult.PercentSigned > 0 {
			nonZeroQuorums = append(nonZeroQuorums, quorumID)
			quorumResults[quorumID] = quorumResult.PercentSigned
		}
	}
	if len(nonZeroQuorums) == 0 {
		err = d.updateBatchStatus(ctx, batchData, quorumResults)
		if err != nil {
			return fmt.Errorf("failed to update blob statuses for batch %s: %w", batchHeaderHash, err)
		}
		return fmt.Errorf("all quorums received no attestation for batch %s", batchHeaderHash)
	}

	aggSig, err := d.aggregator.AggregateSignatures(ctx, d.chainState, uint(batchData.Batch.BatchHeader.ReferenceBlockNumber), quorumAttestation, nonZeroQuorums)
	aggregateSignaturesFinished := time.Now()
	d.metrics.reportAggregateSignaturesLatency(aggregateSignaturesFinished.Sub(receiveSignaturesFinished))
	if err != nil {
		dbErr := d.failBatch(ctx, batchData)
		if dbErr != nil {
			return fmt.Errorf("failed to update blob statuses for batch %s to failed: %w", batchHeaderHash, dbErr)
		}
		return fmt.Errorf("failed to aggregate signatures for batch %s: %w", batchHeaderHash, err)
	}

	attestation := &corev2.Attestation{
		BatchHeader:      batchData.Batch.BatchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: aggSig.NonSigners,
		APKG2:            aggSig.AggPubKey,
		QuorumAPKs:       aggSig.QuorumAggPubKeys,
		Sigma:            aggSig.AggSignature,
		QuorumNumbers:    nonZeroQuorums,
		QuorumResults:    quorumResults,
	}
	err = d.blobMetadataStore.PutAttestation(ctx, attestation)
	putAttestationFinished := time.Now()
	d.metrics.reportPutAttestationLatency(putAttestationFinished.Sub(aggregateSignaturesFinished))
	if err != nil {
		dbErr := d.failBatch(ctx, batchData)
		if dbErr != nil {
			return fmt.Errorf("failed to update blob statuses for batch %s to failed: %w", batchHeaderHash, dbErr)
		}
		return fmt.Errorf("failed to put attestation for batch %s: %w", batchHeaderHash, err)
	}

	err = d.updateBatchStatus(ctx, batchData, attestation.QuorumResults)
	updateBatchStatusFinished := time.Now()
	d.metrics.reportUpdateBatchStatusLatency(updateBatchStatusFinished.Sub(putAttestationFinished))
	if err != nil {
		return fmt.Errorf("failed to update blob statuses for batch %s: %w", batchHeaderHash, err)
	}

	d.logger.Debug("successfully processed batch", "batchHeader", batchHeaderHash)
	return nil
}

// NewBatch creates a batch of blobs to dispatch
// Warning: This function is not thread-safe
func (d *Dispatcher) NewBatch(ctx context.Context, referenceBlockNumber uint64) (*batchData, error) {
	newBatchStart := time.Now()
	defer func() {
		d.metrics.reportNewBatchLatency(time.Since(newBatchStart))
	}()

	blobMetadatas, cursor, err := d.blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, d.cursor, d.MaxBatchSize)
	getBlobMetadataFinished := time.Now()
	d.metrics.reportGetBlobMetadataLatency(getBlobMetadataFinished.Sub(newBatchStart))
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata by status: %w", err)
	}

	if len(blobMetadatas) == 0 {
		return nil, errNoBlobsToDispatch
	}
	d.logger.Debug("got new metadatas to make batch", "numBlobs", len(blobMetadatas), "referenceBlockNumber", referenceBlockNumber)

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

	batchHeaderHash, err := batchHeader.Hash()
	buildMerkleTreeFinished := time.Now()
	d.metrics.reportBuildMerkleTreeLatency(buildMerkleTreeFinished.Sub(getBlobCertificatesFinished))
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}

	err = d.blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	putBatchHeaderFinished := time.Now()
	d.metrics.reportPutBatchHeaderLatency(putBatchHeaderFinished.Sub(buildMerkleTreeFinished))
	if err != nil {
		return nil, fmt.Errorf("failed to put batch header: %w", err)
	}

	// accumulate verification infos in a map to avoid duplicate entries
	// batch write operation fails if there are duplicate entries
	verificationInfoMap := make(map[corev2.BlobKey]*corev2.BlobVerificationInfo)
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

		verificationInfoMap[blobKey] = &corev2.BlobVerificationInfo{
			BatchHeader:    batchHeader,
			BlobKey:        blobKey,
			BlobIndex:      uint32(i),
			InclusionProof: core.SerializeMerkleProof(merkleProof),
		}
	}

	proofGenerationFinished := time.Now()
	d.metrics.reportProofLatency(proofGenerationFinished.Sub(putBatchHeaderFinished))

	verificationInfos := make([]*corev2.BlobVerificationInfo, len(verificationInfoMap))
	i := 0
	for _, v := range verificationInfoMap {
		verificationInfos[i] = v
		i++
	}
	err = d.blobMetadataStore.PutBlobVerificationInfos(ctx, verificationInfos)
	putBlobVerificationInfosFinished := time.Now()
	d.metrics.reportPutVerificationInfosLatency(putBlobVerificationInfosFinished.Sub(proofGenerationFinished))
	if err != nil {
		return nil, fmt.Errorf("failed to put blob verification infos: %w", err)
	}

	if cursor != nil {
		d.cursor = cursor
	}

	d.logger.Debug("new batch", "referenceBlockNumber", referenceBlockNumber, "numBlobs", len(certs))
	return &batchData{
		Batch: &corev2.Batch{
			BatchHeader:      batchHeader,
			BlobCertificates: certs,
		},
		BatchHeaderHash: batchHeaderHash,
		BlobKeys:        keys,
		Metadata:        metadataMap,
		OperatorState:   state,
	}, nil
}

// GetOperatorState returns the operator state for the given quorums at the given block number
func (d *Dispatcher) GetOperatorState(ctx context.Context, metadatas []*v2.BlobMetadata, blockNumber uint64) (*core.IndexedOperatorState, error) {
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

func (d *Dispatcher) sendChunks(ctx context.Context, client clients.NodeClient, batch *corev2.Batch) (*core.Signature, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.NodeRequestTimeout)
	defer cancel()

	sig, err := client.StoreChunks(ctxWithTimeout, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to store chunks: %w", err)
	}

	return sig, nil
}

func (d *Dispatcher) updateBatchStatus(ctx context.Context, batch *batchData, quorumResults map[core.QuorumID]uint8) error {
	var multierr error
	for i, cert := range batch.Batch.BlobCertificates {
		blobKey := batch.BlobKeys[i]
		if cert == nil || cert.BlobHeader == nil {
			d.logger.Error("invalid blob certificate in batch")
			err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
			if err != nil {
				multierr = multierror.Append(multierr, fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Failed)
			}
			continue
		}

		failed := false
		for _, q := range cert.BlobHeader.QuorumNumbers {
			if res, ok := quorumResults[q]; !ok || res == 0 {
				d.logger.Error("quorum result not found", "quorumID", q, "blobKey", blobKey.Hex())
				failed = true
				break
			}
		}

		if failed {
			err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.InsufficientSignatures)
			if err != nil {
				multierr = multierror.Append(multierr, fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
			}
			if metadata, ok := batch.Metadata[blobKey]; ok {
				d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.InsufficientSignatures)
			}
			continue
		}

		err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Certified)
		if err != nil {
			multierr = multierror.Append(multierr, fmt.Errorf("failed to update blob status for blob %s to certified: %w", blobKey.Hex(), err))
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			requestedAt := time.Unix(0, int64(metadata.RequestedAt))
			d.metrics.reportE2EDispersalLatency(time.Since(requestedAt))
			d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Certified)
		}
	}

	return multierr
}

func (d *Dispatcher) failBatch(ctx context.Context, batch *batchData) error {
	var multierr error
	for _, blobKey := range batch.BlobKeys {
		err := d.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
		if err != nil {
			multierr = multierror.Append(multierr, fmt.Errorf("failed to update blob status for blob %s to failed: %w", blobKey.Hex(), err))
		}
		if metadata, ok := batch.Metadata[blobKey]; ok {
			d.metrics.reportCompletedBlob(int(metadata.BlobSize), v2.Failed)
		}
	}

	return multierr
}
