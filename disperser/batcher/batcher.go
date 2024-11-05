package batcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/workerpool"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/go-merkletree/v2"
)

const (
	QuantizationFactor = uint(1)
	indexerWarmupDelay = 2 * time.Second
)

type BatchPlan struct {
	IncludedBlobs []*disperser.BlobMetadata
	Quorums       map[core.QuorumID]QuorumInfo
	State         *core.IndexedOperatorState
}

type QuorumInfo struct {
	Assignments        map[core.OperatorID]core.Assignment
	Info               core.AssignmentInfo
	QuantizationFactor uint
}

type TimeoutConfig struct {
	EncodingTimeout     time.Duration
	AttestationTimeout  time.Duration
	ChainReadTimeout    time.Duration
	ChainWriteTimeout   time.Duration
	ChainStateTimeout   time.Duration
	TxnBroadcastTimeout time.Duration
}

type Config struct {
	PullInterval             time.Duration
	FinalizerInterval        time.Duration
	FinalizerPoolSize        int
	EncoderSocket            string
	SRSOrder                 int
	NumConnections           int
	EncodingRequestQueueSize int
	// BatchSizeMBLimit is the maximum size of a batch in MB
	BatchSizeMBLimit     uint
	MaxNumRetriesPerBlob uint

	FinalizationBlockDelay uint

	TargetNumChunks          uint
	MaxBlobsToFetchFromStore int
}

type Batcher struct {
	Config
	TimeoutConfig

	Queue         disperser.BlobStore
	Dispatcher    disperser.Dispatcher
	EncoderClient disperser.EncoderClient

	ChainState            core.IndexedChainState
	AssignmentCoordinator core.AssignmentCoordinator
	Aggregator            core.SignatureAggregator
	EncodingStreamer      *EncodingStreamer
	Transactor            core.Writer
	TransactionManager    TxnManager
	Metrics               *Metrics
	HeartbeatChan         chan time.Time

	ethClient common.EthClient
	finalizer Finalizer
	logger    logging.Logger
}

func NewBatcher(
	config Config,
	timeoutConfig TimeoutConfig,
	queue disperser.BlobStore,
	dispatcher disperser.Dispatcher,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	encoderClient disperser.EncoderClient,
	aggregator core.SignatureAggregator,
	ethClient common.EthClient,
	finalizer Finalizer,
	transactor core.Writer,
	txnManager TxnManager,
	logger logging.Logger,
	metrics *Metrics,
	heartbeatChan chan time.Time,
) (*Batcher, error) {
	batchTrigger := NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		uint64(config.BatchSizeMBLimit)*1024*1024, // convert to bytes
	)
	streamerConfig := StreamerConfig{
		SRSOrder:                 config.SRSOrder,
		EncodingRequestTimeout:   config.PullInterval,
		EncodingQueueLimit:       config.EncodingRequestQueueSize,
		TargetNumChunks:          config.TargetNumChunks,
		MaxBlobsToFetchFromStore: config.MaxBlobsToFetchFromStore,
		FinalizationBlockDelay:   config.FinalizationBlockDelay,
		ChainStateTimeout:        timeoutConfig.ChainStateTimeout,
	}
	encodingWorkerPool := workerpool.New(config.NumConnections)
	encodingStreamer, err := NewEncodingStreamer(streamerConfig, queue, chainState, encoderClient, assignmentCoordinator, batchTrigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, metrics, logger)
	if err != nil {
		return nil, err
	}

	return &Batcher{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		Queue:         queue,
		Dispatcher:    dispatcher,
		EncoderClient: encoderClient,

		ChainState:            chainState,
		AssignmentCoordinator: assignmentCoordinator,
		Aggregator:            aggregator,
		EncodingStreamer:      encodingStreamer,
		Transactor:            transactor,
		TransactionManager:    txnManager,
		Metrics:               metrics,

		ethClient:     ethClient,
		finalizer:     finalizer,
		logger:        logger.With("component", "Batcher"),
		HeartbeatChan: heartbeatChan,
	}, nil
}

func (b *Batcher) RecoverState(ctx context.Context) error {
	b.logger.Info("Recovering state...")
	start := time.Now()
	metas, err := b.Queue.GetBlobMetadataByStatus(ctx, disperser.Dispersing)
	if err != nil {
		return fmt.Errorf("failed to get blobs in dispersing state: %w", err)
	}
	expired := 0
	processing := 0
	for _, meta := range metas {
		if meta.Expiry == 0 || meta.Expiry < uint64(time.Now().Unix()) {
			err = b.Queue.MarkBlobFailed(ctx, meta.GetBlobKey())
			if err != nil {
				return fmt.Errorf("failed to mark blob (%s) as failed: %w", meta.GetBlobKey(), err)
			}
			expired += 1
		} else {
			err = b.Queue.MarkBlobProcessing(ctx, meta.GetBlobKey())
			if err != nil {
				return fmt.Errorf("failed to mark blob (%s) as processing: %w", meta.GetBlobKey(), err)
			}
			processing += 1
		}
	}
	b.logger.Info("Recovering state took", "duration", time.Since(start), "numBlobs", len(metas), "expired", expired, "processing", processing)
	return nil
}

func (b *Batcher) Start(ctx context.Context) error {
	err := b.RecoverState(ctx)
	if err != nil {
		return fmt.Errorf("failed to recover state: %w", err)
	}
	err = b.ChainState.Start(ctx)
	if err != nil {
		return err
	}
	// Wait for few seconds for indexer to index blockchain
	// This won't be needed when we switch to using Graph node
	time.Sleep(indexerWarmupDelay)
	err = b.EncodingStreamer.Start(ctx)
	if err != nil {
		return err
	}
	batchTrigger := b.EncodingStreamer.EncodedSizeNotifier

	go func() {
		receiptChan := b.TransactionManager.ReceiptChan()
		for {
			select {
			case <-ctx.Done():
				return
			case receiptOrErr := <-receiptChan:
				b.logger.Info("received response from transaction manager", "receipt", receiptOrErr.Receipt, "err", receiptOrErr.Err)
				err := b.ProcessConfirmedBatch(ctx, receiptOrErr)
				if err != nil {
					b.logger.Error("failed to process confirmed batch", "err", err)
				}
			}
		}
	}()
	b.TransactionManager.Start(ctx)

	b.finalizer.Start(ctx)

	go func() {
		ticker := time.NewTicker(b.PullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := b.HandleSingleBatch(ctx); err != nil {
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Warn("no encoded results to make a batch with")
					} else {
						b.logger.Error("failed to process a batch", "err", err)
					}
				}
			case <-batchTrigger.Notify:
				ticker.Stop()

				if err := b.HandleSingleBatch(ctx); err != nil {
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Warn("no encoded results to make a batch with")
					} else {
						b.logger.Error("failed to process a batch", "err", err)
					}
				}
				ticker.Reset(b.PullInterval)
			}
		}
	}()

	return nil
}

// updateConfirmationInfo updates the confirmation info for each blob in the batch and returns failed blobs to retry.
func (b *Batcher) updateConfirmationInfo(
	ctx context.Context,
	batchData confirmationMetadata,
	txnReceipt *types.Receipt,
) ([]*disperser.BlobMetadata, error) {
	if txnReceipt.BlockNumber == nil {
		return nil, errors.New("HandleSingleBatch: error getting transaction receipt block number")
	}
	if len(batchData.blobs) == 0 {
		return nil, errors.New("failed to process confirmed batch: no blobs from transaction manager metadata")
	}
	if batchData.batchHeader == nil {
		return nil, errors.New("failed to process confirmed batch: batch header from transaction manager metadata is nil")
	}
	if len(batchData.blobHeaders) == 0 {
		return nil, errors.New("failed to process confirmed batch: no blob headers from transaction manager metadata")
	}
	if batchData.merkleTree == nil {
		return nil, errors.New("failed to process confirmed batch: merkle tree from transaction manager metadata is nil")
	}
	if batchData.aggSig == nil {
		return nil, errors.New("failed to process confirmed batch: aggSig from transaction manager metadata is nil")
	}
	headerHash, err := batchData.batchHeader.GetBatchHeaderHash()
	if err != nil {
		return nil, fmt.Errorf("HandleSingleBatch: error getting batch header hash: %w", err)
	}
	batchID, err := b.getBatchID(ctx, txnReceipt)
	if err != nil {
		return nil, fmt.Errorf("HandleSingleBatch: error fetching batch ID: %w", err)
	}

	blobsToRetry := make([]*disperser.BlobMetadata, 0)
	var updateConfirmationInfoErr error

	for blobIndex, metadata := range batchData.blobs {
		// Mark the blob failed if it didn't get enough signatures.
		status := disperser.InsufficientSignatures

		var proof []byte
		if isBlobAttested(batchData.aggSig.QuorumResults, batchData.blobHeaders[blobIndex]) {
			status = disperser.Confirmed
			// generate inclusion proof
			if blobIndex >= len(batchData.blobHeaders) {
				b.logger.Error("HandleSingleBatch: error confirming blobs: blob header not found in batch", "index", blobIndex)
				blobsToRetry = append(blobsToRetry, batchData.blobs[blobIndex])
				continue
			}

			merkleProof, err := batchData.merkleTree.GenerateProofWithIndex(uint64(blobIndex), 0)
			if err != nil {
				b.logger.Error("HandleSingleBatch: failed to generate blob header inclusion proof", "err", err)
				blobsToRetry = append(blobsToRetry, batchData.blobs[blobIndex])
				continue
			}
			proof = serializeProof(merkleProof)
		}

		confirmationInfo := &disperser.ConfirmationInfo{
			BatchHeaderHash:         headerHash,
			BlobIndex:               uint32(blobIndex),
			SignatoryRecordHash:     core.ComputeSignatoryRecordHash(uint32(batchData.batchHeader.ReferenceBlockNumber), batchData.aggSig.NonSigners),
			ReferenceBlockNumber:    uint32(batchData.batchHeader.ReferenceBlockNumber),
			BatchRoot:               batchData.batchHeader.BatchRoot[:],
			BlobInclusionProof:      proof,
			BlobCommitment:          &batchData.blobHeaders[blobIndex].BlobCommitments,
			BatchID:                 uint32(batchID),
			ConfirmationTxnHash:     txnReceipt.TxHash,
			ConfirmationBlockNumber: uint32(txnReceipt.BlockNumber.Uint64()),
			Fee:                     []byte{0}, // No fee
			QuorumResults:           batchData.aggSig.QuorumResults,
			BlobQuorumInfos:         batchData.blobHeaders[blobIndex].QuorumInfos,
		}

		if status == disperser.Confirmed {
			if _, updateConfirmationInfoErr = b.Queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Confirmed)
			}
		} else if status == disperser.InsufficientSignatures {
			if _, updateConfirmationInfoErr = b.Queue.MarkBlobInsufficientSignatures(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.InsufficientSignatures)
			}
		} else {
			updateConfirmationInfoErr = fmt.Errorf("HandleSingleBatch: trying to update confirmation info for blob in status other than confirmed or insufficient signatures: %s", status.String())
		}
		if updateConfirmationInfoErr != nil {
			b.logger.Error("HandleSingleBatch: error updating blob confirmed metadata", "err", updateConfirmationInfoErr)
			blobsToRetry = append(blobsToRetry, batchData.blobs[blobIndex])
		}
		requestTime := time.Unix(0, int64(metadata.RequestMetadata.RequestedAt))
		b.Metrics.ObserveLatency("E2E", float64(time.Since(requestTime).Milliseconds()))
		b.Metrics.ObserveBlobAge("confirmed", float64(time.Since(requestTime).Milliseconds()))
		for _, quorumInfo := range batchData.blobHeaders[blobIndex].QuorumInfos {
			b.Metrics.IncrementBlobSize("confirmed", quorumInfo.QuorumID, int(metadata.RequestMetadata.BlobSize))
		}
	}

	return blobsToRetry, nil
}

func (b *Batcher) ProcessConfirmedBatch(ctx context.Context, receiptOrErr *ReceiptOrErr) error {
	if receiptOrErr.Metadata == nil {
		return errors.New("failed to process confirmed batch: no metadata from transaction manager response")
	}
	confirmationMetadata := receiptOrErr.Metadata.(confirmationMetadata)
	blobs := confirmationMetadata.blobs
	if len(blobs) == 0 {
		return errors.New("failed to process confirmed batch: no blobs from transaction manager metadata")
	}
	if receiptOrErr.Err != nil {
		_ = b.handleFailure(ctx, blobs, FailConfirmBatch)
		return fmt.Errorf("failed to confirm batch onchain: %w", receiptOrErr.Err)
	}
	if confirmationMetadata.aggSig == nil {
		_ = b.handleFailure(ctx, blobs, FailNoAggregatedSignature)
		return errors.New("failed to process confirmed batch: aggSig from transaction manager metadata is nil")
	}
	b.logger.Info("received ConfirmBatch transaction receipt", "blockNumber", receiptOrErr.Receipt.BlockNumber, "txnHash", receiptOrErr.Receipt.TxHash.Hex())

	// Mark the blobs as complete
	stageTimer := time.Now()
	blobsToRetry, err := b.updateConfirmationInfo(ctx, confirmationMetadata, receiptOrErr.Receipt)
	if err != nil {
		_ = b.handleFailure(ctx, blobs, FailUpdateConfirmationInfo)
		return fmt.Errorf("failed to update confirmation info: %w", err)
	}
	if len(blobsToRetry) > 0 {
		b.logger.Error("failed to update confirmation info", "failed", len(blobsToRetry), "total", len(blobs))
		_ = b.handleFailure(ctx, blobsToRetry, FailUpdateConfirmationInfo)
	}
	b.logger.Debug("Update confirmation info took", "duration", time.Since(stageTimer).String())
	b.Metrics.ObserveLatency("UpdateConfirmationInfo", float64(time.Since(stageTimer).Milliseconds()))
	batchSize := int64(0)
	for _, blobMeta := range blobs {
		batchSize += int64(blobMeta.RequestMetadata.BlobSize)
	}
	b.Metrics.IncrementBatchCount(batchSize)

	return nil
}

func (b *Batcher) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	numPermanentFailures := 0
	for _, metadata := range blobMetadatas {
		b.EncodingStreamer.RemoveEncodedBlob(metadata)
		retry, err := b.Queue.HandleBlobFailure(ctx, metadata, b.MaxNumRetriesPerBlob)
		if err != nil {
			b.logger.Error("HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}

		if retry {
			continue
		}

		if reason == FailNoSignatures {
			b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.InsufficientSignatures)
		} else {
			b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
		}
		numPermanentFailures++
	}
	b.Metrics.UpdateBatchError(reason, numPermanentFailures)

	// Return the error(s)
	return result.ErrorOrNil()
}

type confirmationMetadata struct {
	batchID     uuid.UUID
	batchHeader *core.BatchHeader
	blobs       []*disperser.BlobMetadata
	blobHeaders []*core.BlobHeader
	merkleTree  *merkletree.MerkleTree
	aggSig      *core.SignatureAggregation
}

func (b *Batcher) observeBlobAge(stage string, batch *batch) {
	for _, m := range batch.BlobMetadata {
		requestTime := time.Unix(0, int64(m.RequestMetadata.RequestedAt))
		b.Metrics.ObserveBlobAge(stage, float64(time.Since(requestTime).Milliseconds()))
	}
}

func (b *Batcher) observeBlobAgeAndSize(stage string, batch *batch) {
	for i, m := range batch.BlobMetadata {
		requestTime := time.Unix(0, int64(m.RequestMetadata.RequestedAt))
		b.Metrics.ObserveBlobAge(stage, float64(time.Since(requestTime).Milliseconds()))
		for _, quorumInfo := range batch.BlobHeaders[i].QuorumInfos {
			b.Metrics.IncrementBlobSize(stage, quorumInfo.QuorumID, int(m.RequestMetadata.BlobSize))
		}
	}
}

func (b *Batcher) HandleSingleBatch(ctx context.Context) error {
	log := b.logger

	// Signal Liveness to indicate no stall
	b.signalLiveness()

	// start a timer
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		b.Metrics.ObserveLatency("total", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	stageTimer := time.Now()
	batch, err := b.EncodingStreamer.CreateBatch(ctx)
	if err != nil {
		return err
	}
	log.Debug("CreateBatch took", "duration", time.Since(stageTimer))
	b.observeBlobAge("batched", batch)

	// Dispatch encoded batch
	log.Debug("Dispatching encoded batch...")
	stageTimer = time.Now()
	update := b.Dispatcher.DisperseBatch(ctx, batch.State, batch.EncodedBlobs, batch.BatchHeader)
	log.Debug("DisperseBatch took", "duration", time.Since(stageTimer))
	b.observeBlobAge("attestation_requested", batch)
	h, err := batch.State.OperatorState.Hash()
	if err != nil {
		log.Error("HandleSingleBatch: error getting operator state hash", "err", err)
	}
	hStr := make([]string, 0, len(h))
	for q, hash := range h {
		hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
	}
	log.Info("Dispatched encoded batch", "operatorStateHash", hStr)

	// Get the batch header hash
	log.Debug("Getting batch header hash...")
	headerHash, err := batch.BatchHeader.GetBatchHeaderHash()
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchHeaderHash)
		return fmt.Errorf("HandleSingleBatch: error getting batch header hash: %w", err)
	}

	// Aggregate the signatures
	log.Debug("Aggregating signatures...")

	stageTimer = time.Now()
	quorumAttestation, err := b.Aggregator.ReceiveSignatures(ctx, batch.State, headerHash, update)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailAggregateSignatures)
		return fmt.Errorf("HandleSingleBatch: error receiving and validating signatures: %w", err)
	}
	operatorCount := make(map[core.QuorumID]int)
	signerCount := make(map[core.QuorumID]int)
	for quorumID, opState := range batch.State.Operators {
		operatorCount[quorumID] = len(opState)
		if _, ok := signerCount[quorumID]; !ok {
			signerCount[quorumID] = 0
		}
		for opID := range opState {
			if _, ok := quorumAttestation.SignerMap[opID]; ok {
				signerCount[quorumID]++
			}
		}
	}
	b.Metrics.UpdateAttestation(operatorCount, signerCount, quorumAttestation.QuorumResults)
	for _, quorumResult := range quorumAttestation.QuorumResults {
		log.Info("Aggregated quorum result", "quorumID", quorumResult.QuorumID, "percentSigned", quorumResult.PercentSigned)
	}

	b.observeBlobAgeAndSize("attested", batch)

	numPassed, passedQuorums := numBlobsAttestedByQuorum(quorumAttestation.QuorumResults, batch.BlobHeaders)
	// TODO(mooselumph): Determine whether to confirm the batch based on the number of successes
	if numPassed == 0 {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailNoSignatures)
		return errors.New("HandleSingleBatch: no blobs received sufficient signatures")
	}

	nonEmptyQuorums := []core.QuorumID{}
	for quorumID := range passedQuorums {
		log.Info("Quorums successfully attested", "quorumID", quorumID)
		nonEmptyQuorums = append(nonEmptyQuorums, quorumID)
	}

	// Aggregate the signatures across only the non-empty quorums. Excluding empty quorums reduces the gas cost.
	aggSig, err := b.Aggregator.AggregateSignatures(ctx, b.ChainState, batch.BatchHeader.ReferenceBlockNumber, quorumAttestation, nonEmptyQuorums)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailAggregateSignatures)
		return fmt.Errorf("HandleSingleBatch: error aggregating signatures: %w", err)
	}

	log.Debug("AggregateSignatures took", "duration", time.Since(stageTimer))
	b.Metrics.ObserveLatency("AggregateSignatures", float64(time.Since(stageTimer).Milliseconds()))

	// Confirm the batch
	log.Debug("Confirming batch...")

	txn, err := b.Transactor.BuildConfirmBatchTxn(ctx, batch.BatchHeader, aggSig.QuorumResults, aggSig)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailConfirmBatch)
		return fmt.Errorf("HandleSingleBatch: error building confirmBatch transaction: %w", err)
	}
	err = b.TransactionManager.ProcessTransaction(ctx, NewTxnRequest(txn, "confirmBatch", big.NewInt(0), confirmationMetadata{
		batchID:     uuid.Nil,
		batchHeader: batch.BatchHeader,
		blobs:       batch.BlobMetadata,
		blobHeaders: batch.BlobHeaders,
		merkleTree:  batch.MerkleTree,
		aggSig:      aggSig,
	}))
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailConfirmBatch)
		return fmt.Errorf("HandleSingleBatch: error sending confirmBatch transaction: %w", err)
	}

	return nil
}

func serializeProof(proof *merkletree.Proof) []byte {
	proofBytes := make([]byte, 0)
	for _, hash := range proof.Hashes {
		proofBytes = append(proofBytes, hash[:]...)
	}
	return proofBytes
}

func (b *Batcher) parseBatchIDFromReceipt(txReceipt *types.Receipt) (uint32, error) {
	if len(txReceipt.Logs) == 0 {
		return 0, errors.New("failed to get transaction receipt with logs")
	}
	for _, log := range txReceipt.Logs {
		if len(log.Topics) == 0 {
			b.logger.Debug("transaction receipt has no topics")
			continue
		}
		b.logger.Debug("[getBatchIDFromReceipt] ", "sigHash", log.Topics[0].Hex())

		if log.Topics[0] == common.BatchConfirmedEventSigHash {
			smAbi, err := abi.JSON(bytes.NewReader(common.ServiceManagerAbi))
			if err != nil {
				return 0, fmt.Errorf("failed to parse ServiceManager ABI: %w", err)
			}
			eventAbi, err := smAbi.EventByID(common.BatchConfirmedEventSigHash)
			if err != nil {
				return 0, fmt.Errorf("failed to parse BatchConfirmed event ABI: %w", err)
			}
			unpackedData, err := eventAbi.Inputs.Unpack(log.Data)
			if err != nil {
				return 0, fmt.Errorf("failed to unpack BatchConfirmed log data: %w", err)
			}

			// There should be exactly one input in the data field, batchId.
			// Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L17
			if len(unpackedData) != 1 {
				return 0, fmt.Errorf("BatchConfirmed log should contain exactly 1 inputs. Found %d", len(unpackedData))
			}
			return unpackedData[0].(uint32), nil
		}
	}
	return 0, errors.New("failed to find BatchConfirmed log from the transaction")
}

func (b *Batcher) getBatchID(ctx context.Context, txReceipt *types.Receipt) (uint32, error) {
	const (
		maxRetries = 4
		baseDelay  = 1 * time.Second
	)
	var (
		batchID uint32
		err     error
	)

	batchID, err = b.parseBatchIDFromReceipt(txReceipt)
	if err == nil {
		return batchID, nil
	}

	txHash := txReceipt.TxHash
	for i := 0; i < maxRetries; i++ {
		retrySec := math.Pow(2, float64(i))
		b.logger.Warn("failed to get transaction receipt, retrying...", "retryIn", retrySec, "err", err)
		time.Sleep(time.Duration(retrySec) * baseDelay)

		txReceipt, err = b.ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			continue
		}

		batchID, err = b.parseBatchIDFromReceipt(txReceipt)
		if err == nil {
			return batchID, nil
		}
	}

	if err != nil {
		b.logger.Warn("failed to get transaction receipt after retries", "numRetries", maxRetries, "err", err)
		return 0, err
	}

	return batchID, nil
}

// numBlobsAttestedByQuorum returns two values:
// 1. the number of blobs that have been successfully attested by all quorums
// 2. map[QuorumID]struct{} contains quorums that have been successfully attested by the quorum (has at least one blob attested in the quorum)
func numBlobsAttestedByQuorum(signedQuorums map[core.QuorumID]*core.QuorumResult, headers []*core.BlobHeader) (int, map[core.QuorumID]struct{}) {
	numPassed := 0
	quorums := make(map[core.QuorumID]struct{})
	for _, blob := range headers {
		thisPassed := true
		for _, quorum := range blob.QuorumInfos {
			if signedQuorums[quorum.QuorumID].PercentSigned < quorum.ConfirmationThreshold {
				thisPassed = false
			} else {
				quorums[quorum.QuorumID] = struct{}{}
			}
		}
		if thisPassed {
			numPassed++
		}
	}

	return numPassed, quorums
}

func isBlobAttested(signedQuorums map[core.QuorumID]*core.QuorumResult, header *core.BlobHeader) bool {
	for _, quorum := range header.QuorumInfos {
		if _, ok := signedQuorums[quorum.QuorumID]; !ok {
			return false
		}
		if signedQuorums[quorum.QuorumID].PercentSigned < quorum.ConfirmationThreshold {
			return false
		}
	}
	return true
}

func (b *Batcher) signalLiveness() {
	select {
	case b.HeartbeatChan <- time.Now():
		b.logger.Info("Heartbeat signal sent")
	default:
		// This case happens if there's no receiver ready to consume the heartbeat signal.
		// It prevents the goroutine from blocking if the channel is full or not being listened to.
		b.logger.Warn("Heartbeat signal skipped, no receiver on the channel")
	}
}
