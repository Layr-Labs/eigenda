package batcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/go-merkletree"
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
	EncodingTimeout    time.Duration
	AttestationTimeout time.Duration
	ChainReadTimeout   time.Duration
	ChainWriteTimeout  time.Duration
}

type Config struct {
	PullInterval             time.Duration
	FinalizerInterval        time.Duration
	EncoderSocket            string
	SRSOrder                 int
	NumConnections           int
	EncodingRequestQueueSize int
	// BatchSizeMBLimit is the maximum size of a batch in MB
	BatchSizeMBLimit     uint
	MaxNumRetriesPerBlob uint

	TargetNumChunks          uint
	MaxBlobsToFetchFromStore int
}

type Batcher struct {
	Config
	TimeoutConfig

	Queue         disperser.BlobStore
	Dispatcher    disperser.Dispatcher
	Confirmer     disperser.BatchConfirmer
	EncoderClient disperser.EncoderClient

	ChainState            core.IndexedChainState
	AssignmentCoordinator core.AssignmentCoordinator
	Aggregator            core.SignatureAggregator
	EncodingStreamer      *EncodingStreamer
	Metrics               *Metrics

	ethClient common.EthClient
	finalizer Finalizer
	logger    common.Logger
}

func NewBatcher(
	config Config,
	timeoutConfig TimeoutConfig,
	queue disperser.BlobStore,
	dispatcher disperser.Dispatcher,
	confirmer disperser.BatchConfirmer,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	encoderClient disperser.EncoderClient,
	aggregator core.SignatureAggregator,
	ethClient common.EthClient,
	finalizer Finalizer,
	logger common.Logger,
	metrics *Metrics,
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
	}
	encodingWorkerPool := workerpool.New(config.NumConnections)
	encodingStreamer, err := NewEncodingStreamer(streamerConfig, queue, chainState, encoderClient, assignmentCoordinator, batchTrigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
	if err != nil {
		return nil, err
	}

	return &Batcher{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		Queue:         queue,
		Dispatcher:    dispatcher,
		Confirmer:     confirmer,
		EncoderClient: encoderClient,

		ChainState:            chainState,
		AssignmentCoordinator: assignmentCoordinator,
		Aggregator:            aggregator,
		EncodingStreamer:      encodingStreamer,
		Metrics:               metrics,

		ethClient: ethClient,
		finalizer: finalizer,
		logger:    logger,
	}, nil
}

func (b *Batcher) Start(ctx context.Context) error {
	err := b.ChainState.Start(ctx)
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

func (b *Batcher) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	for _, metadata := range blobMetadatas {
		err := b.Queue.HandleBlobFailure(ctx, metadata, b.MaxNumRetriesPerBlob)
		if err != nil {
			b.logger.Error("HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}
		b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
	}
	b.Metrics.UpdateBatchError(reason, len(blobMetadatas))

	// Return the error(s)
	return result.ErrorOrNil()
}
func (b *Batcher) HandleSingleBatch(ctx context.Context) error {
	log := b.logger
	// start a timer
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		b.Metrics.ObserveLatency("total", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	stageTimer := time.Now()
	batch, err := b.EncodingStreamer.CreateBatch()
	if err != nil {
		return err
	}
	log.Trace("[batcher] CreateBatch took", "duration", time.Since(stageTimer))

	// Dispatch encoded batch
	log.Trace("[batcher] Dispatching encoded batch...")
	stageTimer = time.Now()
	update := b.Dispatcher.DisperseBatch(ctx, batch.State, batch.EncodedBlobs, batch.BatchHeader)
	log.Trace("[batcher] DisperseBatch took", "duration", time.Since(stageTimer))

	// Get the batch header hash
	log.Trace("[batcher] Getting batch header hash...")
	headerHash, err := batch.BatchHeader.GetBatchHeaderHash()
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchHeaderHash)
		return fmt.Errorf("HandleSingleBatch: error getting batch header hash: %w", err)
	}

	// Aggregate the signatures
	log.Trace("[batcher] Aggregating signatures...")

	// construct quorumParams
	quorumIDs := make([]core.QuorumID, 0, len(batch.State.AggKeys))
	for quorumID := range batch.State.Operators {
		quorumIDs = append(quorumIDs, quorumID)
	}
	fmt.Println("quorumIDs", quorumIDs)

	stageTimer = time.Now()
	aggSig, err := b.Aggregator.AggregateSignatures(batch.State, quorumIDs, headerHash, update)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailAggregateSignatures)
		return fmt.Errorf("HandleSingleBatch: error aggregating signatures: %w", err)
	}
	log.Trace("[batcher] AggregateSignatures took", "duration", time.Since(stageTimer))
	b.Metrics.ObserveLatency("AggregateSignatures", float64(time.Since(stageTimer).Milliseconds()))
	b.Metrics.UpdateAttestation(len(batch.State.IndexedOperators), len(aggSig.NonSigners))

	passed, numPassed := getBlobQuorumPassStatus(aggSig.QuorumResults, batch.BlobHeaders)
	// TODO(mooselumph): Determine whether to confirm the batch based on the number of successes
	if numPassed == 0 {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailNoSignatures)
		return fmt.Errorf("HandleSingleBatch: no blobs received sufficient signatures")
	}

	// Confirm the batch
	log.Trace("[batcher] Confirming batch...")
	stageTimer = time.Now()
	txnReceipt, err := b.Confirmer.ConfirmBatch(ctx, batch.BatchHeader, aggSig.QuorumResults, aggSig)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailConfirmBatch)
		return fmt.Errorf("HandleSingleBatch: error confirming batch: %w", err)
	}
	log.Trace("[batcher] ConfirmBatch took", "duration", time.Since(stageTimer))
	log.Info("[batcher] Batch confirmed at block", "blockNumber", txnReceipt.BlockNumber, "txnHash", txnReceipt.TxHash.Hex())
	b.Metrics.ObserveLatency("ConfirmBatch", float64(time.Since(stageTimer).Milliseconds()))
	b.Metrics.GasUsed.Set(float64(txnReceipt.GasUsed))

	batchID, err := b.getBatchID(ctx, txnReceipt)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailGetBatchID)
		return fmt.Errorf("HandleSingleBatch: error fetching batch ID: %w", err)
	}

	// Mark the blobs as complete
	log.Trace("[batcher] Marking blobs as complete...")
	stageTimer = time.Now()
	blobsToRetry := make([]*disperser.BlobMetadata, 0)
	var updateConfirmationInfoErr error
	for blobIndex, metadata := range batch.BlobMetadata {
		// Mark the blob failed if it didn't get enough signatures.
		status := disperser.Confirmed
		if !passed[blobIndex] {
			status = disperser.InsufficientSignatures
		}

		var blobHeader *core.BlobHeader
		var proof []byte
		if status == disperser.Confirmed {
			// generate inclusion proof
			if blobIndex >= len(batch.BlobHeaders) {
				return fmt.Errorf("HandleSingleBatch: error confirming blobs: blob header at index %d not found in batch", blobIndex)
			}
			blobHeader = batch.BlobHeaders[blobIndex]

			blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
			if err != nil {
				return fmt.Errorf("HandleSingleBatch: failed to get blob header hash: %w", err)
			}
			merkleProof, err := batch.MerkleTree.GenerateProof(blobHeaderHash[:], 0)
			if err != nil {
				return fmt.Errorf("HandleSingleBatch: failed to generate blob header inclusion proof: %w", err)
			}
			proof = serializeProof(merkleProof)
		}

		confirmationInfo := &disperser.ConfirmationInfo{
			BatchHeaderHash:         headerHash,
			BlobIndex:               uint32(blobIndex),
			SignatoryRecordHash:     core.ComputeSignatoryRecordHash(uint32(batch.BatchHeader.ReferenceBlockNumber), aggSig.NonSigners),
			ReferenceBlockNumber:    uint32(batch.BatchHeader.ReferenceBlockNumber),
			BatchRoot:               batch.BatchHeader.BatchRoot[:],
			BlobInclusionProof:      proof,
			BlobCommitment:          &batch.BlobHeaders[blobIndex].BlobCommitments,
			BatchID:                 uint32(batchID),
			ConfirmationTxnHash:     txnReceipt.TxHash,
			ConfirmationBlockNumber: uint32(txnReceipt.BlockNumber.Uint64()),
			Fee:                     []byte{0}, // No fee
			QuorumResults:           aggSig.QuorumResults,
			BlobQuorumInfos:         batch.BlobHeaders[blobIndex].QuorumInfos,
		}

		if status == disperser.Confirmed {
			if _, updateConfirmationInfoErr = b.Queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Confirmed)
				// remove encoded blob from storage so we don't disperse it again
				b.EncodingStreamer.RemoveEncodedBlob(metadata)
			}
		} else if status == disperser.InsufficientSignatures {
			if _, updateConfirmationInfoErr = b.Queue.MarkBlobInsufficientSignatures(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.InsufficientSignatures)
				// remove encoded blob from storage so we don't disperse it again
				b.EncodingStreamer.RemoveEncodedBlob(metadata)
			}
		} else {
			updateConfirmationInfoErr = fmt.Errorf("HandleSingleBatch: trying to update confirmation info for blob in status other than confirmed or insufficient signatures: %s", status.String())
		}
		if updateConfirmationInfoErr != nil {
			log.Error("HandleSingleBatch: error updating blob confirmed metadata", "err", updateConfirmationInfoErr)
			blobsToRetry = append(blobsToRetry, batch.BlobMetadata[blobIndex])
		}
		requestTime := time.Unix(0, int64(metadata.RequestMetadata.RequestedAt))
		b.Metrics.ObserveLatency("E2E", float64(time.Since(requestTime).Milliseconds()))
	}

	if len(blobsToRetry) > 0 {
		_ = b.handleFailure(ctx, blobsToRetry, FailUpdateConfirmationInfo)
		if len(blobsToRetry) == len(batch.BlobMetadata) {
			return fmt.Errorf("HandleSingleBatch: failed to update blob confirmed metadata for all blobs in batch: %w", updateConfirmationInfoErr)
		}
	}

	log.Trace("[batcher] Update confirmation info took", "duration", time.Since(stageTimer))
	b.Metrics.ObserveLatency("UpdateConfirmationInfo", float64(time.Since(stageTimer).Milliseconds()))
	batchSize := int64(0)
	for _, blobMeta := range batch.BlobMetadata {
		batchSize += int64(blobMeta.RequestMetadata.BlobSize)
	}
	b.Metrics.IncrementBatchCount(batchSize)
	return nil
}

func serializeProof(proof *merkletree.Proof) []byte {
	proofBytes := make([]byte, 0)
	for _, hash := range proof.Hashes {
		proofBytes = append(proofBytes, hash[:]...)
	}
	return proofBytes
}

func (b *Batcher) parseBatchIDFromReceipt(ctx context.Context, txReceipt *types.Receipt) (uint32, error) {
	if len(txReceipt.Logs) == 0 {
		return 0, fmt.Errorf("failed to get transaction receipt with logs")
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
				return 0, err
			}
			eventAbi, err := smAbi.EventByID(common.BatchConfirmedEventSigHash)
			if err != nil {
				return 0, err
			}
			unpackedData, err := eventAbi.Inputs.Unpack(log.Data)
			if err != nil {
				return 0, err
			}

			// There should be exactly two inputs in the data field, batchId and fee.
			// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L20
			if len(unpackedData) != 2 {
				return 0, fmt.Errorf("BatchConfirmed log should contain exactly 2 inputs. Found %d", len(unpackedData))
			}
			return unpackedData[0].(uint32), nil
		}
	}
	return 0, fmt.Errorf("failed to find BatchConfirmed log from the transaction")
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

	batchID, err = b.parseBatchIDFromReceipt(ctx, txReceipt)
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

		batchID, err = b.parseBatchIDFromReceipt(ctx, txReceipt)
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

// Determine failure status for each blob based on stake signed per quorum. We fail a blob if it received
// insufficient signatures for any quorum
func getBlobQuorumPassStatus(signedQuorums map[core.QuorumID]*core.QuorumResult, headers []*core.BlobHeader) ([]bool, int) {
	numPassed := 0
	passed := make([]bool, len(headers))
	for ind, blob := range headers {
		thisPassed := true
		for _, quorum := range blob.QuorumInfos {
			if signedQuorums[quorum.QuorumID].PercentSigned < quorum.QuorumThreshold {
				thisPassed = false
				break
			}
		}
		passed[ind] = thisPassed
		if thisPassed {
			numPassed++
		}
	}

	return passed, numPassed
}
