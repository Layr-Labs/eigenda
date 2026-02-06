package chainstate

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	ejectionmanager "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Indexer indexes operator state events from Ethereum contracts.
type Indexer struct {
	config    *IndexerConfig
	store     store.Store
	persister *store.JSONPersister

	ethClient *ethclient.Client

	// Contract bindings
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
	blsApkRegistry      *blsapkregistry.ContractBLSApkRegistry
	ejectionManager     *ejectionmanager.ContractEjectionManager

	logger logging.Logger
}

// NewIndexer creates a new chainstate indexer.
func NewIndexer(
	ctx context.Context,
	config *IndexerConfig,
	ethClient *ethclient.Client,
	logger logging.Logger,
) (*Indexer, error) {
	// Get contract addresses from EigenDADirectory
	var (
		registryCoordinatorAddr gethcommon.Address
		blsApkRegistryAddr      gethcommon.Address
		ejectionManagerAddr     gethcommon.Address
	)

	contractDirectory, err := directory.NewContractDirectory(ctx, logger, ethClient,
		gethcommon.HexToAddress(config.EigenDADirectory))
	if err != nil {
		return nil, fmt.Errorf("new contract directory: %w", err)
	}

	// Registry Coordinator
	registryCoordinatorAddr, err = contractDirectory.GetContractAddress(ctx, directory.RegistryCoordinator)
	if err != nil {
		return nil, fmt.Errorf("get registry coordinator addr: %w", err)
	}

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddr,
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator binding: %w", err)
	}

	// BLS APK Registry
	blsApkRegistryAddr, err = registryCoordinator.BlsApkRegistry(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get BLS APK registry address: %w", err)
	}

	blsApkRegistry, err := blsapkregistry.NewContractBLSApkRegistry(
		blsApkRegistryAddr,
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLS APK registry binding: %w", err)
	}

	// Ejection Manager
	ejectionManagerAddr, err = contractDirectory.GetContractAddress(ctx, directory.EigenDAEjectionManager)
	if err != nil {
		return nil, fmt.Errorf("get ejection manager addr: %w", err)
	}

	ejectionManager, err := ejectionmanager.NewContractEjectionManager(
		ejectionManagerAddr,
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ejection manager binding: %w", err)
	}

	// Create store and persister
	memStore := store.NewMemoryStore()
	persister := store.NewJSONPersister(memStore, config.PersistencePath, logger)

	return &Indexer{
		config:              config,
		store:               memStore,
		persister:           persister,
		ethClient:           ethClient,
		registryCoordinator: registryCoordinator,
		blsApkRegistry:      blsApkRegistry,
		ejectionManager:     ejectionManager,
		logger:              logger.With("component", "ChainStateIndexer"),
	}, nil
}

// Start starts the indexer by loading persisted state and beginning the indexing loop.
func (i *Indexer) Start(ctx context.Context) error {
	// Load persisted state
	if err := i.persister.Load(ctx); err != nil {
		return fmt.Errorf("failed to load persisted state: %w", err)
	}

	// Start periodic persistence
	go i.persister.StartPeriodicSave(ctx, i.config.PersistInterval)

	// Start indexing loop
	go i.indexLoop(ctx)

	i.logger.Info("Indexer started successfully")
	return nil
}

// indexLoop continuously polls for new blocks and indexes them.
func (i *Indexer) indexLoop(ctx context.Context) {
	ticker := time.NewTicker(i.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := i.indexNewBlocks(ctx); err != nil {
				i.logger.Error("Failed to index blocks", "error", err)
			}
		case <-ctx.Done():
			i.logger.Info("Index loop stopped")
			return
		}
	}
}

// indexNewBlocks indexes all new blocks since the last indexed block.
func (i *Indexer) indexNewBlocks(ctx context.Context) error {
	lastIndexed, err := i.store.GetLastIndexedBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last indexed block: %w", err)
	}

	if lastIndexed == 0 {
		lastIndexed = i.config.StartBlockNumber
		if lastIndexed == 0 {
			// If no start block configured, start from current block
			currentBlock, err := i.ethClient.BlockNumber(ctx)
			if err != nil {
				return fmt.Errorf("failed to get current block: %w", err)
			}
			lastIndexed = currentBlock
			i.logger.Info("Starting indexing from current block", "block", lastIndexed)
		}
	}

	latestBlock, err := i.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	if lastIndexed >= latestBlock {
		return nil // Already caught up
	}

	// Process in batches
	fromBlock := lastIndexed + 1
	toBlock := min(fromBlock+i.config.BlockBatchSize-1, latestBlock)

	i.logger.Debug("Indexing block range", "from", fromBlock, "to", toBlock)

	// Index events from each contract
	if err := i.indexRegistryCoordinatorEvents(ctx, fromBlock, toBlock); err != nil {
		return fmt.Errorf("failed to index registry coordinator events: %w", err)
	}
	if err := i.indexBLSApkRegistryEvents(ctx, fromBlock, toBlock); err != nil {
		return fmt.Errorf("failed to index BLS APK registry events: %w", err)
	}
	if err := i.indexEjectionManagerEvents(ctx, fromBlock, toBlock); err != nil {
		return fmt.Errorf("failed to index ejection manager events: %w", err)
	}

	// Update last indexed block
	if err := i.store.SetLastIndexedBlock(ctx, toBlock); err != nil {
		return fmt.Errorf("failed to set last indexed block: %w", err)
	}

	i.logger.Info("Indexed blocks", "from", fromBlock, "to", toBlock)

	return nil
}

// indexRegistryCoordinatorEvents indexes events from the RegistryCoordinator contract.
func (i *Indexer) indexRegistryCoordinatorEvents(ctx context.Context, from, to uint64) error {
	filterOpts := &bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}

	// Index OperatorRegistered events
	regIter, err := i.registryCoordinator.FilterOperatorRegistered(filterOpts, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorRegistered events: %w", err)
	}
	defer regIter.Close()

	for regIter.Next() {
		event := regIter.Event
		operator := &types.Operator{
			ID:                      event.OperatorId,
			Address:                 event.Operator,
			RegisteredAtBlockNumber: event.Raw.BlockNumber,
			RegisteredTxHash:        event.Raw.TxHash,
		}

		if err := i.store.SaveOperator(ctx, operator); err != nil {
			return fmt.Errorf("failed to save operator: %w", err)
		}

		i.logger.Debug("Indexed operator registration", "operator_id", fmt.Sprintf("%x", event.OperatorId), "block", event.Raw.BlockNumber)
	}

	if err := regIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorRegistered events: %w", err)
	}

	// Index OperatorDeregistered events
	deregIter, err := i.registryCoordinator.FilterOperatorDeregistered(filterOpts, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorDeregistered events: %w", err)
	}
	defer deregIter.Close()

	for deregIter.Next() {
		event := deregIter.Event
		if err := i.store.DeregisterOperator(ctx, event.OperatorId, event.Raw.BlockNumber, event.Raw.TxHash); err != nil {
			i.logger.Warn("Failed to deregister operator (may not exist yet)", "operator_id", fmt.Sprintf("%x", event.OperatorId), "error", err)
			continue
		}

		i.logger.Debug("Indexed operator deregistration", "operator_id", fmt.Sprintf("%x", event.OperatorId), "block", event.Raw.BlockNumber)
	}

	if err := deregIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorDeregistered events: %w", err)
	}

	// Index OperatorSocketUpdate events
	socketIter, err := i.registryCoordinator.FilterOperatorSocketUpdate(filterOpts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorSocketUpdate events: %w", err)
	}
	defer socketIter.Close()

	for socketIter.Next() {
		event := socketIter.Event
		update := &types.OperatorSocketUpdate{
			OperatorID:  event.OperatorId,
			Socket:      event.Socket,
			BlockNumber: event.Raw.BlockNumber,
			TxHash:      event.Raw.TxHash,
			UpdatedAt:   time.Now(),
		}

		if err := i.store.SaveSocketUpdate(ctx, update); err != nil {
			return fmt.Errorf("failed to save socket update: %w", err)
		}

		// Also update the operator's socket
		if err := i.store.UpdateOperatorSocket(ctx, event.OperatorId, event.Socket, event.Raw.BlockNumber); err != nil {
			i.logger.Warn("Failed to update operator socket (may not exist yet)", "operator_id", fmt.Sprintf("%x", event.OperatorId), "error", err)
		}

		i.logger.Debug("Indexed socket update", "operator_id", fmt.Sprintf("%x", event.OperatorId), "socket", event.Socket, "block", event.Raw.BlockNumber)
	}

	if err := socketIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorSocketUpdate events: %w", err)
	}

	return nil
}

// indexBLSApkRegistryEvents indexes events from the BLSApkRegistry contract.
func (i *Indexer) indexBLSApkRegistryEvents(ctx context.Context, from, to uint64) error {
	filterOpts := &bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}

	// Index NewPubkeyRegistration events
	pubkeyIter, err := i.blsApkRegistry.FilterNewPubkeyRegistration(filterOpts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter NewPubkeyRegistration events: %w", err)
	}
	defer pubkeyIter.Close()

	for pubkeyIter.Next() {
		event := pubkeyIter.Event

		// Convert BN254 points to core.G1Point and core.G2Point
		g1Point := core.NewG1Point(event.PubkeyG1.X, event.PubkeyG1.Y)

		// For G2Point, we need to manually construct the bn254.G2Affine
		// G2 X and Y are E2 extension field elements with two components each
		var g2Affine bn254.G2Affine
		g2Affine.X.SetString(event.PubkeyG2.X[0].String(), event.PubkeyG2.X[1].String())
		g2Affine.Y.SetString(event.PubkeyG2.Y[0].String(), event.PubkeyG2.Y[1].String())
		g2Point := &core.G2Point{G2Affine: &g2Affine}

		// Convert operator address to operator ID using the contract
		operatorID, err := i.registryCoordinator.GetOperatorId(nil, event.Operator)
		if err != nil {
			i.logger.Warn("Failed to get operator ID for pubkey registration", "operator", event.Operator, "error", err)
			continue
		}

		// Get the operator and update their BLS keys
		op, err := i.store.GetOperator(ctx, operatorID)
		if err != nil {
			i.logger.Warn("Operator not found when registering pubkey", "operator_id", fmt.Sprintf("%x", operatorID), "error", err)
			continue
		}

		op.BLSPubKeyG1 = g1Point
		op.BLSPubKeyG2 = g2Point

		if err := i.store.SaveOperator(ctx, op); err != nil {
			return fmt.Errorf("failed to update operator pubkey: %w", err)
		}

		i.logger.Debug("Indexed BLS pubkey registration", "operator_id", fmt.Sprintf("%x", event.Operator), "block", event.Raw.BlockNumber)
	}

	if err := pubkeyIter.Error(); err != nil {
		return fmt.Errorf("error iterating NewPubkeyRegistration events: %w", err)
	}

	// Index OperatorAddedToQuorums events
	addedIter, err := i.blsApkRegistry.FilterOperatorAddedToQuorums(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorAddedToQuorums events: %w", err)
	}
	defer addedIter.Close()

	for addedIter.Next() {
		event := addedIter.Event

		// Get the operator and update their quorum memberships
		op, err := i.store.GetOperator(ctx, event.OperatorId)
		if err != nil {
			i.logger.Warn("Operator not found when adding to quorums", "operator_id", fmt.Sprintf("%x", event.OperatorId), "error", err)
			continue
		}

		// Add new quorums (avoiding duplicates)
		for _, newQuorum := range event.QuorumNumbers {
			found := false
			for _, existingQuorum := range op.QuorumIDs {
				if existingQuorum == core.QuorumID(newQuorum) {
					found = true
					break
				}
			}
			if !found {
				op.QuorumIDs = append(op.QuorumIDs, core.QuorumID(newQuorum))
			}
		}

		if err := i.store.SaveOperator(ctx, op); err != nil {
			return fmt.Errorf("failed to update operator quorums: %w", err)
		}

		i.logger.Debug("Indexed operator added to quorums", "operator_id", fmt.Sprintf("%x", event.OperatorId), "quorums", event.QuorumNumbers, "block", event.Raw.BlockNumber)
	}

	if err := addedIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorAddedToQuorums events: %w", err)
	}

	// Index OperatorRemovedFromQuorums events
	removedIter, err := i.blsApkRegistry.FilterOperatorRemovedFromQuorums(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorRemovedFromQuorums events: %w", err)
	}
	defer removedIter.Close()

	for removedIter.Next() {
		event := removedIter.Event

		// Get the operator and update their quorum memberships
		op, err := i.store.GetOperator(ctx, event.OperatorId)
		if err != nil {
			i.logger.Warn("Operator not found when removing from quorums", "operator_id", fmt.Sprintf("%x", event.OperatorId), "error", err)
			continue
		}

		// Remove quorums
		var newQuorumIDs []core.QuorumID
		for _, existingQuorum := range op.QuorumIDs {
			shouldRemove := false
			for _, removedQuorum := range event.QuorumNumbers {
				if existingQuorum == core.QuorumID(removedQuorum) {
					shouldRemove = true
					break
				}
			}
			if !shouldRemove {
				newQuorumIDs = append(newQuorumIDs, existingQuorum)
			}
		}
		op.QuorumIDs = newQuorumIDs

		if err := i.store.SaveOperator(ctx, op); err != nil {
			return fmt.Errorf("failed to update operator quorums: %w", err)
		}

		i.logger.Debug("Indexed operator removed from quorums", "operator_id", fmt.Sprintf("%x", event.OperatorId), "quorums", event.QuorumNumbers, "block", event.Raw.BlockNumber)
	}

	if err := removedIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorRemovedFromQuorums events: %w", err)
	}

	return nil
}

// indexEjectionManagerEvents indexes events from the EjectionManager contract.
func (i *Indexer) indexEjectionManagerEvents(ctx context.Context, from, to uint64) error {
	filterOpts := &bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}

	// Index OperatorEjected events
	ejectedIter, err := i.ejectionManager.FilterOperatorEjected(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter OperatorEjected events: %w", err)
	}
	defer ejectedIter.Close()

	for ejectedIter.Next() {
		event := ejectedIter.Event

		// OperatorEjected event has a single QuorumNumber, not QuorumNumbers array
		quorumIDs := []core.QuorumID{core.QuorumID(event.QuorumNumber)}

		ejection := &types.OperatorEjection{
			OperatorID:  event.OperatorId,
			QuorumIDs:   quorumIDs,
			BlockNumber: event.Raw.BlockNumber,
			TxHash:      event.Raw.TxHash,
			EjectedAt:   time.Now(),
		}

		if err := i.store.SaveEjection(ctx, ejection); err != nil {
			return fmt.Errorf("failed to save ejection: %w", err)
		}

		i.logger.Debug("Indexed operator ejection", "operator_id", fmt.Sprintf("%x", event.OperatorId), "quorum", event.QuorumNumber, "block", event.Raw.BlockNumber)
	}

	if err := ejectedIter.Error(); err != nil {
		return fmt.Errorf("error iterating OperatorEjected events: %w", err)
	}

	return nil
}

// GetStore returns the underlying store (useful for API server).
func (i *Indexer) GetStore() store.Store {
	return i.store
}
