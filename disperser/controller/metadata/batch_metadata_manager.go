package metadata

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TODO future Cody: add a mock instance and get the unit tests passing

type BatchMetadataManager interface {
	GetMetadata() *BatchMetadata
}

var _ BatchMetadataManager = (*batchMetadataManager)(nil)

// The BatchMetadataManager responsible for providing BatchMetadata for use in the creation of new batches.
// This utility periodically downloads recent onchain date, so that new batches are created with recent
// batch metadata.
type batchMetadataManager struct {
	ctx    context.Context
	logger logging.Logger

	// The underlying eth client used to interact with the blockchain.
	ethClient bind.ContractBackend

	// Used to get operator state. The IndexedChainState utility fetches state both from onchain sources and from
	// the indexer. When we eventually move all data onchain, we can ditch the indexer and just call directly
	// into the contract bindings in this file.
	indexedChainState core.IndexedChainState

	// A handle for communicating with the registry coordinator contract.
	registryCoordinator *regcoordinator.ContractRegistryCoordinator

	// When choosing a new reference block number (RBN), select the block that is this many blocks in the past.
	// This is a hedge against forking.
	referenceBlockOffset uint64

	// The time between updates to the metadata.
	updatePeriod time.Duration

	// The most recent batch metadata.
	metadata atomic.Pointer[BatchMetadata]
}

// Create a new BatchMetadataManager with the specified quorums.
//
// updatePeriod is the period at which the manager will update its metadata. Actual update timing may vary
// depending on the amount of time it takes to successfully get new data.
func NewBatchMetadataManager(
	ctx context.Context,
	logger logging.Logger,
	ethClient bind.ContractBackend,
	indexedChainState core.IndexedChainState,
	registryCoordinatorAddress gethcommon.Address,
	updatePeriod time.Duration,
	referenceBlockOffset uint64,
) (BatchMetadataManager, error) {

	registryCoordinator, err := regcoordinator.NewContractRegistryCoordinator(registryCoordinatorAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator client: %w", err)
	}

	manager := &batchMetadataManager{
		ctx:                  ctx,
		logger:               logger,
		ethClient:            ethClient,
		metadata:             atomic.Pointer[BatchMetadata]{},
		indexedChainState:    indexedChainState,
		registryCoordinator:  registryCoordinator,
		referenceBlockOffset: referenceBlockOffset,
		updatePeriod:         updatePeriod,
	}

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

// Fetch the next reference block number (RBN) to use.
func (m *batchMetadataManager) getNextReferenceBlockNumber() (uint64, error) {
	// Get the latest block header to determine the current reference block number.
	latestHeader, err := m.ethClient.HeaderByNumber(m.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block header: %w", err)
	}
	latestBlockNumber := latestHeader.Number.Uint64()

	if latestBlockNumber < m.referenceBlockOffset {
		return 0, fmt.Errorf("latest block number is less than RBN offset: %d < %d",
			latestBlockNumber, m.referenceBlockOffset)
	}

	return latestBlockNumber - m.referenceBlockOffset, nil
}

// get a list of all quorums that are registered for a particular reference block number.
func (m *batchMetadataManager) getQuorums(referenceBlockNumber uint64) ([]core.QuorumID, error) {

	// Quorums are assigned starting at 0, and then sequentially without gaps. If we
	// know the number of quorums, we can generate a list of quorum IDs.

	quorumCount, err := m.registryCoordinator.QuorumCount(&bind.CallOpts{
		Context:     m.ctx,
		BlockNumber: new(big.Int).SetUint64(referenceBlockNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum count: %w", err)
	}

	quorums := make([]core.QuorumID, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		quorums[i] = i
	}

	return quorums, nil
}

// updateMetadata fetches the latest batch metadata from the blockchain and updates m.operatorState.
// This method is called periodically to ensure that metadata reflects a recent(ish) reference block.
func (m *batchMetadataManager) updateMetadata() error {
	referenceBlockNumber, err := m.getNextReferenceBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get next reference block number: %w", err)
	}
	if referenceBlockNumber < m.metadata.Load().referenceBlockNumber {
		// Only update if the new RBN is greater than the most recent one.
		return nil
	}

	quorums, err := m.getQuorums(referenceBlockNumber)
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

	for m.ctx.Err() == nil {
		<-ticker.C

		err := m.updateMetadata()
		if err != nil {
			m.logger.Errorf("failed to update metadata: %v", err)
		}
	}
}
