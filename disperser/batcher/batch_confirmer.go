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
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
)

type BatchConfirmerConfig struct {
	PullInterval                 time.Duration
	DispersalTimeout             time.Duration
	DispersalStatusCheckInterval time.Duration
	AttestationTimeout           time.Duration
	SRSOrder                     int
	NumConnections               int
	MaxNumRetriesPerBlob         uint
}

type BatchConfirmer struct {
	BatchConfirmerConfig

	BlobStore      disperser.BlobStore
	MinibatchStore MinibatchStore
	Dispatcher     disperser.Dispatcher
	EncoderClient  disperser.EncoderClient

	ChainState         core.IndexedChainState
	EncodingStreamer   *EncodingStreamer
	Aggregator         core.SignatureAggregator
	Transactor         core.Transactor
	TransactionManager TxnManager
	Minibatcher        *Minibatcher

	ethClient common.EthClient
	logger    logging.Logger
}

func NewBatchConfirmer(
	config BatchConfirmerConfig,
	blobStore disperser.BlobStore,
	minibatchStore MinibatchStore,
	dispatcher disperser.Dispatcher,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	encodingStreamer *EncodingStreamer,
	aggregator core.SignatureAggregator,
	ethClient common.EthClient,
	transactor core.Transactor,
	txnManager TxnManager,
	minibatcher *Minibatcher,
	logger logging.Logger,
) (*BatchConfirmer, error) {
	return &BatchConfirmer{
		BatchConfirmerConfig: config,

		BlobStore:        blobStore,
		MinibatchStore:   minibatchStore,
		Dispatcher:       dispatcher,
		EncodingStreamer: encodingStreamer,

		ChainState:         chainState,
		Aggregator:         aggregator,
		Transactor:         transactor,
		TransactionManager: txnManager,
		Minibatcher:        minibatcher,

		ethClient: ethClient,
		logger:    logger.With("component", "BatchConfirmer"),
	}, nil
}

func (b *BatchConfirmer) Start(ctx context.Context) error {
	err := b.ChainState.Start(ctx)
	if err != nil {
		return err
	}
	// Wait for few seconds for indexer to index blockchain
	// This won't be needed when we switch to using Graph node
	time.Sleep(indexerWarmupDelay)

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

	go func() {
		ticker := time.NewTicker(b.PullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := b.HandleSingleBatch(ctx); err != nil {
					b.logger.Error("failed to process a batch", "err", err)
				}
			}
		}
	}()

	return nil
}

// updateConfirmationInfo updates the confirmation info for each blob in the batch and returns failed blobs to retry.
func (b *BatchConfirmer) updateConfirmationInfo(
	ctx context.Context,
	batchData confirmationMetadata,
	txnReceipt *types.Receipt,
) ([]*disperser.BlobMetadata, error) {
	if txnReceipt.BlockNumber == nil {
		return nil, errors.New("error getting transaction receipt block number")
	}
	if batchData.batchID == uuid.Nil {
		return nil, errors.New("failed to process confirmed batch: batch ID from transaction manager metadata is nil")
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
		return nil, fmt.Errorf("error getting batch header hash: %w", err)
	}
	batchID, err := b.getBatchID(ctx, txnReceipt)
	if err != nil {
		return nil, fmt.Errorf("error fetching batch ID: %w", err)
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
			merkleProof, err := batchData.merkleTree.GenerateProofWithIndex(uint64(blobIndex), 0)
			if err != nil {
				b.logger.Error("failed to generate blob header inclusion proof", "err", err)
				blobsToRetry = append(blobsToRetry, batchData.blobs[blobIndex])
				continue
			}
			proof = serializeProof(merkleProof)
		}

		confirmationInfo := &disperser.ConfirmationInfo{
			BatchHeaderHash:      headerHash,
			BlobIndex:            uint32(blobIndex),
			SignatoryRecordHash:  core.ComputeSignatoryRecordHash(uint32(batchData.batchHeader.ReferenceBlockNumber), batchData.aggSig.NonSigners),
			ReferenceBlockNumber: uint32(batchData.batchHeader.ReferenceBlockNumber),
			BatchRoot:            batchData.batchHeader.BatchRoot[:],
			BlobInclusionProof:   proof,
			BlobCommitment:       &batchData.blobHeaders[blobIndex].BlobCommitments,
			// This is onchain, external batch ID, which is different from the internal representation of batch UUID
			BatchID:                 uint32(batchID),
			ConfirmationTxnHash:     txnReceipt.TxHash,
			ConfirmationBlockNumber: uint32(txnReceipt.BlockNumber.Uint64()),
			Fee:                     []byte{0}, // No fee
			QuorumResults:           batchData.aggSig.QuorumResults,
			BlobQuorumInfos:         batchData.blobHeaders[blobIndex].QuorumInfos,
		}

		if status == disperser.Confirmed {
			// TODO: add metrics for confirmed blobs
			if _, updateConfirmationInfoErr = b.BlobStore.MarkBlobConfirmed(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.logger.Info("blob confirmed", "blobKey", metadata.GetBlobKey())
				// b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Confirmed)
			}
		} else if status == disperser.InsufficientSignatures {
			if _, updateConfirmationInfoErr = b.BlobStore.MarkBlobInsufficientSignatures(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
				b.logger.Debug("blob marked as insufficient signatures", "blobKey", metadata.GetBlobKey())
				// b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.InsufficientSignatures)
			}
		} else {
			updateConfirmationInfoErr = fmt.Errorf("trying to update confirmation info for blob in status other than confirmed or insufficient signatures: %s", status.String())
		}
		if updateConfirmationInfoErr != nil {
			b.logger.Error("error updating blob confirmed metadata", "err", updateConfirmationInfoErr)
			blobsToRetry = append(blobsToRetry, batchData.blobs[blobIndex])
		}
	}

	return blobsToRetry, nil
}

func (b *BatchConfirmer) ProcessConfirmedBatch(ctx context.Context, receiptOrErr *ReceiptOrErr) error {
	if receiptOrErr.Metadata == nil {
		return errors.New("failed to process confirmed batch: no metadata from transaction manager response")
	}
	confirmationMetadata := receiptOrErr.Metadata.(confirmationMetadata)
	blobs := confirmationMetadata.blobs
	if len(blobs) == 0 {
		return errors.New("failed to process confirmed batch: no blobs from transaction manager metadata")
	}
	if confirmationMetadata.batchID == uuid.Nil {
		return errors.New("failed to process confirmed batch: batch ID from transaction manager metadata is nil")
	}
	if receiptOrErr.Err != nil {
		_ = b.handleFailure(ctx, confirmationMetadata.batchID, blobs, FailConfirmBatch)
		return fmt.Errorf("failed to confirm batch onchain: %w", receiptOrErr.Err)
	}
	if confirmationMetadata.aggSig == nil {
		_ = b.handleFailure(ctx, confirmationMetadata.batchID, blobs, FailNoAggregatedSignature)
		return errors.New("failed to process confirmed batch: aggSig from transaction manager metadata is nil")
	}
	b.logger.Info("received ConfirmBatch transaction receipt", "blockNumber", receiptOrErr.Receipt.BlockNumber, "txnHash", receiptOrErr.Receipt.TxHash.Hex())

	// Mark the blobs as complete
	stageTimer := time.Now()
	blobsToRetry, err := b.updateConfirmationInfo(ctx, confirmationMetadata, receiptOrErr.Receipt)
	if err != nil {
		_ = b.handleFailure(ctx, confirmationMetadata.batchID, blobs, FailUpdateConfirmationInfo)
		return fmt.Errorf("failed to update confirmation info: %w", err)
	}
	if len(blobsToRetry) > 0 {
		b.logger.Error("failed to update confirmation info", "failed", len(blobsToRetry), "total", len(blobs))
		_ = b.handleFailure(ctx, confirmationMetadata.batchID, blobsToRetry, FailUpdateConfirmationInfo)
	} else {
		err = b.MinibatchStore.UpdateBatchStatus(ctx, confirmationMetadata.batchID, BatchStatusAttested)
		if err != nil {
			b.logger.Error("error updating batch status", "err", err)
		}
	}
	batchSize := int64(0)
	for _, blobMeta := range blobs {
		batchSize += int64(blobMeta.RequestMetadata.BlobSize)
	}
	b.logger.Debug("Update confirmation info took", "duration", time.Since(stageTimer).String(), "batchSize", batchSize)

	return nil
}

func (b *BatchConfirmer) handleFailure(ctx context.Context, batchID uuid.UUID, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	numPermanentFailures := 0
	for _, metadata := range blobMetadatas {
		retry, err := b.BlobStore.HandleBlobFailure(ctx, metadata, b.MaxNumRetriesPerBlob)
		if err != nil {
			b.logger.Error("HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}

		if retry {
			continue
		}

		numPermanentFailures++
	}

	err := b.MinibatchStore.UpdateBatchStatus(ctx, batchID, BatchStatusFailed)
	if err != nil {
		b.logger.Error("error updating batch status", "err", err)
	}

	// Return the error(s)
	return result.ErrorOrNil()
}

func (b *BatchConfirmer) HandleSingleBatch(ctx context.Context) error {
	currentBlock, err := b.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("error getting current block number: %w", err)
	}

	err = b.EncodingStreamer.UpdateReferenceBlock(uint(currentBlock))
	if err != nil {
		return fmt.Errorf("error updating reference block number: %w", err)
	}

	// Get the pending batch
	b.logger.Debug("Getting latest formed batch...")
	stateUpdateTicker := time.NewTicker(b.DispersalStatusCheckInterval)
	var batch *BatchRecord
	for batch == nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stateUpdateTicker.C:
			batch, err = b.MinibatchStore.GetLatestFormedBatch(ctx)
			if err != nil {
				b.logger.Error("error getting latest formed batch", "err", err)
			}
		}
	}

	// Make sure all minibatches in the batch have been dispersed
	batchDispersed := false
	stateUpdateTicker.Reset(b.DispersalStatusCheckInterval)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, b.DispersalTimeout)
	defer cancel()
	for !batchDispersed {
		select {
		case <-ctxWithTimeout.Done():
			return ctxWithTimeout.Err()
		case <-stateUpdateTicker.C:
			batchDispersed, err = b.MinibatchStore.BatchDispersed(ctx, batch.ID, batch.NumMinibatches)
			if err != nil {
				b.logger.Error("error checking if batch is dispersed", "err", err)
			}
		}
	}

	if !batchDispersed {
		return errors.New("batch not dispersed")
	}

	// Try getting batch state from minibatcher cache
	// TODO(ian-shim): If not found, get it from the minibatch store
	batchState := b.Minibatcher.PopBatchState(batch.ID)
	if batchState == nil {
		return fmt.Errorf("no batch state found for batch %s", batch.ID)
	}

	// Construct batch header
	b.logger.Debug("Constructing batch header...")
	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: batch.ReferenceBlockNumber,
		BatchRoot:            [32]byte{},
	}
	blobHeaderHashes := make([][32]byte, 0)
	for _, blobHeader := range batchState.BlobHeaders {
		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			return fmt.Errorf("error getting blob header hash: %w", err)
		}
		blobHeaderHashes = append(blobHeaderHashes, blobHeaderHash)
	}
	merkleTree, err := batchHeader.SetBatchRootFromBlobHeaderHashes(blobHeaderHashes)
	if err != nil {
		return fmt.Errorf("error setting batch root from blob header hashes: %w", err)
	}
	batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
	if err != nil {
		return fmt.Errorf("error getting batch header hash: %w", err)
	}

	// Make AttestBatch call
	b.logger.Debug("Attesting batch...")
	quorumIDsMap := make(map[core.QuorumID]struct{})
	for _, blobHeader := range batchState.BlobHeaders {
		for _, q := range blobHeader.QuorumInfos {
			if _, ok := quorumIDsMap[q.QuorumID]; !ok {
				quorumIDsMap[q.QuorumID] = struct{}{}
			}
		}
	}
	quorumIDs := make([]core.QuorumID, len(quorumIDsMap))
	i := 0
	for q := range quorumIDsMap {
		quorumIDs[i] = q
		i++
	}

	replyChan, err := b.Dispatcher.AttestBatch(ctx, batchState.OperatorState, blobHeaderHashes, batchHeader)
	if err != nil {
		return fmt.Errorf("error making attesting batch request: %w", err)
	}

	// Aggregate the signatures
	b.logger.Debug("Aggregating signatures...")
	quorumAttestation, err := b.Aggregator.ReceiveSignatures(ctx, batchState.OperatorState, batchHeaderHash, replyChan)
	if err != nil {
		_ = b.handleFailure(ctx, batch.ID, batchState.BlobMetadata, FailAggregateSignatures)
		return fmt.Errorf("error receiving and validating signatures: %w", err)
	}
	operatorCount := make(map[core.QuorumID]int)
	signerCount := make(map[core.QuorumID]int)
	for quorumID, opState := range batchState.OperatorState.Operators {
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
	b.logger.Debug("received signatures", "signerCount", signerCount, "operatorCount", operatorCount)
	for _, quorumResult := range quorumAttestation.QuorumResults {
		b.logger.Info("aggregated quorum result", "quorumID", quorumResult.QuorumID, "percentSigned", quorumResult.PercentSigned)
	}

	numPassed, passedQuorums := numBlobsAttestedByQuorum(quorumAttestation.QuorumResults, batchState.BlobHeaders)
	// TODO(mooselumph): Determine whether to confirm the batch based on the number of successes
	if numPassed == 0 {
		_ = b.handleFailure(ctx, batch.ID, batchState.BlobMetadata, FailNoSignatures)
		return errors.New("no blobs received sufficient signatures")
	}

	nonEmptyQuorums := []core.QuorumID{}
	for quorumID := range passedQuorums {
		b.logger.Info("Quorums successfully attested", "quorumID", quorumID)
		nonEmptyQuorums = append(nonEmptyQuorums, quorumID)
	}

	// Aggregate the signatures across only the non-empty quorums. Excluding empty quorums reduces the gas cost.
	aggSig, err := b.Aggregator.AggregateSignatures(ctx, b.ChainState, batchHeader.ReferenceBlockNumber, quorumAttestation, nonEmptyQuorums)
	if err != nil {
		_ = b.handleFailure(ctx, batch.ID, batchState.BlobMetadata, FailAggregateSignatures)
		return fmt.Errorf("error aggregating signatures: %w", err)
	}

	b.logger.Debug("Confirming batch...")

	txn, err := b.Transactor.BuildConfirmBatchTxn(ctx, batchHeader, aggSig.QuorumResults, aggSig)
	if err != nil {
		_ = b.handleFailure(ctx, batch.ID, batchState.BlobMetadata, FailConfirmBatch)
		return fmt.Errorf("error building confirmBatch transaction: %w", err)
	}
	err = b.TransactionManager.ProcessTransaction(ctx, NewTxnRequest(txn, "confirmBatch", big.NewInt(0), confirmationMetadata{
		batchID:     batch.ID,
		batchHeader: batchHeader,
		blobs:       batchState.BlobMetadata,
		blobHeaders: batchState.BlobHeaders,
		merkleTree:  merkleTree,
		aggSig:      aggSig,
	}))
	if err != nil {
		_ = b.handleFailure(ctx, batch.ID, batchState.BlobMetadata, FailConfirmBatch)
		return fmt.Errorf("error sending confirmBatch transaction: %w", err)
	}

	return nil
}

func (b *BatchConfirmer) parseBatchIDFromReceipt(txReceipt *types.Receipt) (uint32, error) {
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

func (b *BatchConfirmer) getBatchID(ctx context.Context, txReceipt *types.Receipt) (uint32, error) {
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
