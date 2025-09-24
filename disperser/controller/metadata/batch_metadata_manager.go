package metadata

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// An object responsible for acquiring and providing batch metadata (i.e. operator state and reference block number)
// for the creation of new batches.
type BatchMetadataManager interface {

	// GetMetadata returns the metadata required to create a new batch. Although the data will be updated periodically,
	// this utility makes no guarantees about the freshness of the data returned by this method. Keeping up to date
	// with the most recent onchain data is done on a best effort basis.
	GetMetadata() *BatchMetadata

	// Release resources associated with this manager.
	Close()
}

var _ BatchMetadataManager = (*batchMetadataManager)(nil)

// A standard implementation of the BatchMetadataManager interface. Does all metadata fetching in a background
// goroutine, guaranteeing that GetMetadata() never blocks.
type batchMetadataManager struct {
	ctx    context.Context
	logger logging.Logger

	// Used to get operator state. The IndexedChainState utility fetches state both from onchain sources and from
	// the indexer. When we eventually move all data onchain, we can ditch the indexer and just call directly
	// into the contract bindings in this file.
	indexedChainState core.IndexedChainState

	// A utility for fetching the list of registered quorums for a given reference block number.
	quorumScanner eth.QuorumScanner

	// Used to look up the reference block number (RBN) to use for batch creation.
	referenceBlockProvider ReferenceBlockProvider

	// The time between updates to the metadata.
	updatePeriod time.Duration

	// The most recent batch metadata.
	metadata atomic.Pointer[BatchMetadata]

	alive atomic.Bool
}

// Create a new BatchMetadataManager.
//
// This constructor does an initial blocking metadata fetch, so that any call to GetBlobMetadata() after this
// constructor returns can immediately return valid metadata. It also starts a background goroutine that periodically
// updates the metadata at a rate defined by updatePeriod. Actual update timing may vary depending on the amount of
// time it takes to successfully get new data.
func NewBatchMetadataManager(
	ctx context.Context,
	logger logging.Logger,
	contractBackend bind.ContractBackend,
	indexedChainState core.IndexedChainState,
	registryCoordinatorAddress gethcommon.Address,
	updatePeriod time.Duration,
	referenceBlockOffset uint64,
) (BatchMetadataManager, error) {

	rbnProvider := NewReferenceBlockProvider(logger, contractBackend, referenceBlockOffset)

	quorumScanner, err := eth.NewQuorumScanner(contractBackend, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create quorum scanner: %w", err)
	}

	manager := &batchMetadataManager{
		ctx:                    ctx,
		logger:                 logger,
		metadata:               atomic.Pointer[BatchMetadata]{},
		indexedChainState:      indexedChainState,
		quorumScanner:          quorumScanner,
		referenceBlockProvider: rbnProvider,
		updatePeriod:           updatePeriod,
	}
	manager.alive.Store(true)

	// Make sure we have valid metadata before the constructor returns.
	err = manager.updateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to update initial metadata: %w", err)
	}

	go manager.updateLoop()

	return manager, nil
}

// GetMetadata returns the most recent batch metadata. This method is thread safe.
func (m *batchMetadataManager) GetMetadata() *BatchMetadata {
	return m.metadata.Load()
}

// Close releases resources associated with this manager.
func (m *batchMetadataManager) Close() {
	m.alive.Store(false)
}

// updateMetadata fetches the latest batch metadata from the blockchain and updates m.operatorState.
// This method is called periodically to ensure that metadata reflects a recent(ish) reference block.
func (m *batchMetadataManager) updateMetadata() error {
	referenceBlockNumber, err := m.referenceBlockProvider.GetReferenceBlockNumber(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to get next reference block number: %w", err)
	}

	previousMetadata := m.metadata.Load()
	if previousMetadata != nil {
		// reference block provider prevents RBN from going backwards
		enforce.GreaterThanOrEqual(referenceBlockNumber, previousMetadata.referenceBlockNumber,
			"reference block number went backwards")

		if referenceBlockNumber == previousMetadata.referenceBlockNumber {
			// Only update if the new RBN is greater than the most recent one.
			m.logger.Debugf("reference block number %d is the same as the previous one, skipping update",
				referenceBlockNumber)
			return nil
		}
	}

	quorums, err := m.quorumScanner.GetQuorums(m.ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get quorums for block %d: %w", referenceBlockNumber, err)
	}

	operatorState, err := m.indexedChainState.GetIndexedOperatorState(m.ctx, uint(referenceBlockNumber), quorums)
	if err != nil {
		return fmt.Errorf("failed to get operator state for block %d: %w", referenceBlockNumber, err)
	}

	m.logger.Debugf("Fetched operator state for block %d, there are %d operators in %d quorums",
		referenceBlockNumber, len(operatorState.IndexedOperators), len(quorums))

	metadata := NewBatchMetadata(referenceBlockNumber, operatorState)
	m.metadata.Store(metadata)

	return nil
}

// periodically updates the batch metadata.
func (m *batchMetadataManager) updateLoop() {
	ticker := time.NewTicker(m.updatePeriod)
	defer ticker.Stop()

	for m.ctx.Err() == nil && m.alive.Load() {
		<-ticker.C

		err := m.updateMetadata()
		if err != nil {
			m.logger.Errorf("failed to update metadata: %v", err)
		}
	}
}
