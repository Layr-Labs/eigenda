package chainstate

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	ejectionmanager "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	stakeregistry "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Indexer indexes operator state events from Ethereum contracts.
type Indexer struct {
	config    *IndexerConfig
	store     store.Store
	persister *store.JSONPersister

	ethClient IndexerEthClient

	// Contract bindings
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
	blsApkRegistry      *blsapkregistry.ContractBLSApkRegistry
	ejectionManager     *ejectionmanager.ContractEjectionManager
	stakeRegistry       *stakeregistry.ContractStakeRegistry

	// wg tracks the background goroutines (index loop and periodic persister)
	// so that Wait can block until the final state save completes on shutdown.
	wg sync.WaitGroup

	logger logging.Logger
}

// IndexerEthClient is the subset of an Ethereum client the indexer needs
// beyond serving as a bind.ContractBackend for the contract bindings. Both
// *ethclient.Client and common.EthClient (e.g. geth.MultiHomingClient)
// satisfy it.
type IndexerEthClient interface {
	bind.ContractBackend
	BlockNumber(ctx context.Context) (uint64, error)
}

// NewIndexer creates a new chainstate indexer.
func NewIndexer(
	ctx context.Context,
	config *IndexerConfig,
	ethClient IndexerEthClient,
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

	// Stake Registry (used to record total quorum stake alongside APK snapshots)
	stakeRegistryAddr, err := contractDirectory.GetContractAddress(ctx, directory.StakeRegistry)
	if err != nil {
		return nil, fmt.Errorf("get stake registry addr: %w", err)
	}

	stakeRegistry, err := stakeregistry.NewContractStakeRegistry(stakeRegistryAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create stake registry binding: %w", err)
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
		stakeRegistry:       stakeRegistry,
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
	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		i.persister.StartPeriodicSave(ctx, i.config.PersistInterval)
	}()

	// Start indexing loop
	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		i.indexLoop(ctx)
	}()

	i.logger.Info("Indexer started successfully")
	return nil
}

// Wait blocks until the indexer's background goroutines have stopped. This
// should be called after the context passed to Start has been cancelled, to
// ensure the persister's final state save completes before the process exits.
func (i *Indexer) Wait() {
	i.wg.Wait()
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

// indexNewBlocks indexes all blocks between the last indexed block and the
// current chain head, in batches of BlockBatchSize. It loops until it has
// caught up to the head observed at call time (any blocks mined meanwhile are
// picked up on the next poll tick), so a backfill proceeds at full speed
// instead of one batch per PollInterval.
func (i *Indexer) indexNewBlocks(ctx context.Context) error {
	latestBlock, err := i.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	lastIndexed, err := i.store.GetLastIndexedBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last indexed block: %w", err)
	}

	// fromBlock is inclusive. On a fresh store (nothing indexed yet) it is the
	// configured start block itself, so events in that block are indexed.
	fromBlock := lastIndexed + 1
	if lastIndexed == 0 {
		if i.config.StartBlockNumber > 0 {
			fromBlock = i.config.StartBlockNumber
		} else {
			fromBlock = latestBlock
			i.logger.Warn(
				"No start block configured; indexing from the current chain head. "+
					"Events before this block (including one-time BLS pubkey registrations) "+
					"will not be indexed. Set StartBlockNumber to the contract deployment "+
					"block to index full history.",
				"block", latestBlock,
			)
		}
	}

	for fromBlock <= latestBlock {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("indexing interrupted: %w", err)
		}

		toBlock := min(fromBlock+i.config.BlockBatchSize-1, latestBlock)

		i.logger.Debug("Indexing block range", "from", fromBlock, "to", toBlock)

		if err := i.indexBlockRange(ctx, fromBlock, toBlock); err != nil {
			return fmt.Errorf("failed to index blocks %d-%d: %w", fromBlock, toBlock, err)
		}

		// Update last indexed block
		if err := i.store.SetLastIndexedBlock(ctx, toBlock); err != nil {
			return fmt.Errorf("failed to set last indexed block: %w", err)
		}

		i.logger.Info("Indexed blocks", "from", fromBlock, "to", toBlock)

		fromBlock = toBlock + 1
	}

	return nil
}

// chainEvent is a decoded contract event paired with its position in the
// chain, so events from all contracts can be replayed in occurrence order.
type chainEvent struct {
	blockNumber uint64
	logIndex    uint
	apply       func(ctx context.Context) error
}

// indexBlockRange indexes all watched contract events in [from, to].
//
// Events are collected from every contract first and then applied in strict
// (block number, log index) order. Applying them per event type instead (all
// registrations, then all deregistrations, ...) corrupts state whenever a
// range contains opposing events for the same operator: a deregistration at
// block N followed by a re-registration at block N+k would be applied in the
// reverse order, leaving a live operator permanently marked deregistered
// (and similarly for quorum remove-then-re-add).
func (i *Indexer) indexBlockRange(ctx context.Context, from, to uint64) error {
	events, affectedQuorums, err := i.collectEvents(ctx, from, to)
	if err != nil {
		return err
	}

	sort.Slice(events, func(a, b int) bool {
		if events[a].blockNumber != events[b].blockNumber {
			return events[a].blockNumber < events[b].blockNumber
		}
		return events[a].logIndex < events[b].logIndex
	})

	for _, event := range events {
		if err := event.apply(ctx); err != nil {
			return err
		}
	}

	// Snapshot each affected quorum's APK once per block, after all membership
	// changes for the range have been applied.
	for blockNum, byQuorum := range affectedQuorums {
		quorums := make([]core.QuorumID, 0, len(byQuorum))
		for quorum := range byQuorum {
			quorums = append(quorums, quorum)
		}
		if err := i.snapshotQuorumAPKs(ctx, quorums, blockNum); err != nil {
			return fmt.Errorf("failed to snapshot quorum APKs: %w", err)
		}
	}

	return nil
}

// collectEvents gathers all watched events in [from, to] from every contract,
// without applying them. It also returns, per block, the set of quorums whose
// aggregate public key changed in that block (so each can be snapshotted once
// per block, however many events share the block).
func (i *Indexer) collectEvents(
	ctx context.Context,
	from, to uint64,
) ([]chainEvent, map[uint64]map[core.QuorumID]struct{}, error) {
	filterOpts := &bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}

	var events []chainEvent

	affectedQuorums := make(map[uint64]map[core.QuorumID]struct{})
	markAffected := func(blockNum uint64, quorums []byte) {
		byQuorum, ok := affectedQuorums[blockNum]
		if !ok {
			byQuorum = make(map[core.QuorumID]struct{})
			affectedQuorums[blockNum] = byQuorum
		}
		for _, quorum := range quorums {
			byQuorum[quorum] = struct{}{}
		}
	}

	// RegistryCoordinator: OperatorRegistered
	regIter, err := i.registryCoordinator.FilterOperatorRegistered(filterOpts, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorRegistered events: %w", err)
	}
	defer func() { _ = regIter.Close() }()
	for regIter.Next() {
		event := regIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorRegistered(ctx, event) },
		})
	}
	if err := regIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorRegistered events: %w", err)
	}

	// RegistryCoordinator: OperatorDeregistered
	deregIter, err := i.registryCoordinator.FilterOperatorDeregistered(filterOpts, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorDeregistered events: %w", err)
	}
	defer func() { _ = deregIter.Close() }()
	for deregIter.Next() {
		event := deregIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorDeregistered(ctx, event) },
		})
	}
	if err := deregIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorDeregistered events: %w", err)
	}

	// RegistryCoordinator: OperatorSocketUpdate
	socketIter, err := i.registryCoordinator.FilterOperatorSocketUpdate(filterOpts, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorSocketUpdate events: %w", err)
	}
	defer func() { _ = socketIter.Close() }()
	for socketIter.Next() {
		event := socketIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorSocketUpdate(ctx, event) },
		})
	}
	if err := socketIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorSocketUpdate events: %w", err)
	}

	// BLSApkRegistry: NewPubkeyRegistration
	pubkeyIter, err := i.blsApkRegistry.FilterNewPubkeyRegistration(filterOpts, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter NewPubkeyRegistration events: %w", err)
	}
	defer func() { _ = pubkeyIter.Close() }()
	for pubkeyIter.Next() {
		event := pubkeyIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleNewPubkeyRegistration(ctx, event) },
		})
	}
	if err := pubkeyIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating NewPubkeyRegistration events: %w", err)
	}

	// BLSApkRegistry: OperatorAddedToQuorums
	addedIter, err := i.blsApkRegistry.FilterOperatorAddedToQuorums(filterOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorAddedToQuorums events: %w", err)
	}
	defer func() { _ = addedIter.Close() }()
	for addedIter.Next() {
		event := addedIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorAddedToQuorums(ctx, event) },
		})
		// The aggregate public key of each affected quorum changed at this block.
		markAffected(event.Raw.BlockNumber, event.QuorumNumbers)
	}
	if err := addedIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorAddedToQuorums events: %w", err)
	}

	// BLSApkRegistry: OperatorRemovedFromQuorums
	removedIter, err := i.blsApkRegistry.FilterOperatorRemovedFromQuorums(filterOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorRemovedFromQuorums events: %w", err)
	}
	defer func() { _ = removedIter.Close() }()
	for removedIter.Next() {
		event := removedIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorRemovedFromQuorums(ctx, event) },
		})
		// The aggregate public key of each affected quorum changed at this block.
		markAffected(event.Raw.BlockNumber, event.QuorumNumbers)
	}
	if err := removedIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorRemovedFromQuorums events: %w", err)
	}

	// EjectionManager: OperatorEjected
	ejectedIter, err := i.ejectionManager.FilterOperatorEjected(filterOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter OperatorEjected events: %w", err)
	}
	defer func() { _ = ejectedIter.Close() }()
	for ejectedIter.Next() {
		event := ejectedIter.Event
		events = append(events, chainEvent{
			blockNumber: event.Raw.BlockNumber,
			logIndex:    event.Raw.Index,
			apply:       func(ctx context.Context) error { return i.handleOperatorEjected(ctx, event) },
		})
	}
	if err := ejectedIter.Error(); err != nil {
		return nil, nil, fmt.Errorf("error iterating OperatorEjected events: %w", err)
	}

	return events, affectedQuorums, nil
}

// handleOperatorRegistered records an operator (re-)registration.
func (i *Indexer) handleOperatorRegistered(
	ctx context.Context,
	event *regcoordinator.ContractEigenDARegistryCoordinatorOperatorRegistered,
) error {
	// A registration may be a RE-registration of an operator that previously
	// deregistered, or follow a NewPubkeyRegistration / OperatorSocketUpdate
	// emitted earlier in the same transaction (the contract emits both before
	// OperatorRegistered), so a record may already exist. Overwriting it with a
	// fresh struct would wipe the BLS keys irrecoverably (NewPubkeyRegistration
	// is emitted at most once per operator ever), so load any existing record
	// and update it in place (mirroring the operator-state subgraph, which only
	// resets the deregistration marker on re-registration).
	operator, err := i.store.GetOperator(ctx, event.OperatorId)
	if errors.Is(err, store.ErrOperatorNotFound) {
		operator = &types.Operator{
			ID: event.OperatorId,
		}
	} else if err != nil {
		return fmt.Errorf("failed to get operator %x: %w", event.OperatorId, err)
	}

	operator.Address = event.Operator
	operator.RegisteredAtBlockNumber = event.Raw.BlockNumber
	operator.RegisteredTxHash = event.Raw.TxHash
	operator.DeregisteredAtBlockNumber = nil
	operator.DeregisteredTxHash = nil

	if err := i.store.SaveOperator(ctx, operator); err != nil {
		return fmt.Errorf("failed to save operator: %w", err)
	}

	i.logger.Debug(
		"Indexed operator registration",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleOperatorDeregistered records an operator deregistration.
func (i *Indexer) handleOperatorDeregistered(
	ctx context.Context,
	event *regcoordinator.ContractEigenDARegistryCoordinatorOperatorDeregistered,
) error {
	// Events are applied in chain order and on-chain a deregistration always
	// follows the operator's registration, so the operator must already exist.
	// A not-found here means a genuine gap (a start block after the operator's
	// registration, a missed earlier range, or a reorg), so fail the batch to
	// retry rather than silently dropping the event.
	if err := i.store.DeregisterOperator(ctx, event.OperatorId, event.Raw.BlockNumber, event.Raw.TxHash); err != nil {
		return fmt.Errorf("failed to deregister operator %x at block %d: %w",
			event.OperatorId, event.Raw.BlockNumber, err)
	}

	i.logger.Debug(
		"Indexed operator deregistration",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleOperatorSocketUpdate records a socket update and applies it to the
// operator's current state.
func (i *Indexer) handleOperatorSocketUpdate(
	ctx context.Context,
	event *regcoordinator.ContractEigenDARegistryCoordinatorOperatorSocketUpdate,
) error {
	update := &types.OperatorSocketUpdate{
		OperatorID:  event.OperatorId,
		Socket:      event.Socket,
		BlockNumber: event.Raw.BlockNumber,
		LogIndex:    event.Raw.Index,
		TxHash:      event.Raw.TxHash,
		UpdatedAt:   time.Now(),
	}

	if err := i.store.SaveSocketUpdate(ctx, update); err != nil {
		return fmt.Errorf("failed to save socket update: %w", err)
	}

	// Also update the operator's socket. Within a registration transaction the
	// contract emits OperatorSocketUpdate BEFORE OperatorRegistered, so the
	// record may not exist yet; create a skeleton that the registration event
	// (later in the same block) fills in. Dropping the update instead would
	// leave the operator without a socket forever, since the event is not
	// re-emitted.
	err := i.store.UpdateOperatorSocket(ctx, event.OperatorId, event.Socket, event.Raw.BlockNumber)
	if errors.Is(err, store.ErrOperatorNotFound) {
		skeleton := &types.Operator{
			ID:     event.OperatorId,
			Socket: event.Socket,
		}
		if err := i.store.SaveOperator(ctx, skeleton); err != nil {
			return fmt.Errorf("failed to save operator for socket update: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to update operator socket: %w", err)
	}

	i.logger.Debug(
		"Indexed socket update",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"socket",
		event.Socket,
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleNewPubkeyRegistration records an operator's BLS public keys.
func (i *Indexer) handleNewPubkeyRegistration(
	ctx context.Context,
	event *blsapkregistry.ContractBLSApkRegistryNewPubkeyRegistration,
) error {
	g1Point := core.NewG1Point(event.PubkeyG1.X, event.PubkeyG1.Y)
	g2Point := core.NewG2Point(event.PubkeyG2.X, event.PubkeyG2.Y)

	// The operator ID is defined as the keccak hash of the G1 pubkey (this is
	// what RegistryCoordinator.getOperatorId returns), so it is computed
	// locally instead of with a per-event contract call.
	operatorID := g1Point.GetOperatorID()

	// The contract emits NewPubkeyRegistration BEFORE OperatorRegistered within
	// the registration transaction, so on an operator's first registration no
	// record exists yet; create a skeleton that the registration event (later
	// in the same block) fills in. The event is emitted at most once per
	// operator ever (the contract forbids re-registering a pubkey), so dropping
	// it would lose the operator's BLS keys irrecoverably.
	op, err := i.store.GetOperator(ctx, operatorID)
	if errors.Is(err, store.ErrOperatorNotFound) {
		op = &types.Operator{
			ID:      operatorID,
			Address: event.Operator,
		}
	} else if err != nil {
		return fmt.Errorf("failed to get operator %x for pubkey registration: %w", operatorID, err)
	}

	op.BLSPubKeyG1 = g1Point
	op.BLSPubKeyG2 = g2Point

	if err := i.store.SaveOperator(ctx, op); err != nil {
		return fmt.Errorf("failed to update operator pubkey: %w", err)
	}

	i.logger.Debug(
		"Indexed BLS pubkey registration",
		"operator_id",
		fmt.Sprintf("%x", operatorID),
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleOperatorAddedToQuorums records quorum-membership additions.
func (i *Indexer) handleOperatorAddedToQuorums(
	ctx context.Context,
	event *blsapkregistry.ContractBLSApkRegistryOperatorAddedToQuorums,
) error {
	// Events are applied in chain order and the contract emits
	// OperatorAddedToQuorums after OperatorRegistered within the registration
	// transaction, so the operator must already exist. A not-found here means a
	// genuine gap (a missed earlier range or a reorg), so fail the batch to
	// retry rather than silently dropping the membership change (which would
	// also skip its APK snapshot).
	op, err := i.store.GetOperator(ctx, event.OperatorId)
	if err != nil {
		return fmt.Errorf("failed to add operator %x to quorums at block %d: %w",
			event.OperatorId, event.Raw.BlockNumber, err)
	}

	// Add new quorums (avoiding duplicates)
	for _, newQuorum := range event.QuorumNumbers {
		found := false
		for _, existingQuorum := range op.QuorumIDs {
			if existingQuorum == newQuorum {
				found = true
				break
			}
		}
		if !found {
			op.QuorumIDs = append(op.QuorumIDs, newQuorum)
		}
	}

	if err := i.store.SaveOperator(ctx, op); err != nil {
		return fmt.Errorf("failed to update operator quorums: %w", err)
	}

	i.logger.Debug(
		"Indexed operator added to quorums",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"quorums",
		event.QuorumNumbers,
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleOperatorRemovedFromQuorums records quorum-membership removals.
func (i *Indexer) handleOperatorRemovedFromQuorums(
	ctx context.Context,
	event *blsapkregistry.ContractBLSApkRegistryOperatorRemovedFromQuorums,
) error {
	// As with adding to quorums, a not-found operator here signals a genuine
	// indexing gap rather than a benign condition, so fail the batch to retry.
	op, err := i.store.GetOperator(ctx, event.OperatorId)
	if err != nil {
		return fmt.Errorf("failed to remove operator %x from quorums at block %d: %w",
			event.OperatorId, event.Raw.BlockNumber, err)
	}

	// Remove quorums
	var newQuorumIDs []core.QuorumID
	for _, existingQuorum := range op.QuorumIDs {
		shouldRemove := false
		for _, removedQuorum := range event.QuorumNumbers {
			if existingQuorum == removedQuorum {
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

	i.logger.Debug(
		"Indexed operator removed from quorums",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"quorums",
		event.QuorumNumbers,
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// handleOperatorEjected records an operator ejection.
func (i *Indexer) handleOperatorEjected(
	ctx context.Context,
	event *ejectionmanager.ContractEjectionManagerOperatorEjected,
) error {
	// OperatorEjected event has a single QuorumNumber, not QuorumNumbers array
	quorumIDs := []core.QuorumID{event.QuorumNumber}

	ejection := &types.OperatorEjection{
		OperatorID:  event.OperatorId,
		QuorumIDs:   quorumIDs,
		BlockNumber: event.Raw.BlockNumber,
		LogIndex:    event.Raw.Index,
		TxHash:      event.Raw.TxHash,
		EjectedAt:   time.Now(),
	}

	if err := i.store.SaveEjection(ctx, ejection); err != nil {
		return fmt.Errorf("failed to save ejection: %w", err)
	}

	i.logger.Debug(
		"Indexed operator ejection",
		"operator_id",
		fmt.Sprintf("%x", event.OperatorId),
		"quorum",
		event.QuorumNumber,
		"block",
		event.Raw.BlockNumber,
	)
	return nil
}

// snapshotQuorumAPKs records an aggregate-public-key snapshot for each of the
// given quorums as of blockNum. The aggregate key is maintained on-chain by the
// BLSApkRegistry, so we read it directly rather than summing operator keys
// ourselves; total quorum stake is read from the StakeRegistry.
//
// To match the operator-state subgraph this indexer replaces, the contract reads
// are made at blockNum (the block of the membership-change event), not at the
// chain head. Reading historical state this way requires the configured RPC to
// be an archive node. The snapshot timestamp is the block's timestamp, again to
// match the subgraph rather than recording wall-clock indexing time.
//
// Any read or save failure is returned so the caller can retry the block range,
// rather than silently leaving a quorum without a snapshot at this block.
func (i *Indexer) snapshotQuorumAPKs(ctx context.Context, quorumNumbers []core.QuorumID, blockNum uint64) error {
	if len(quorumNumbers) == 0 {
		return nil
	}

	blockBig := new(big.Int).SetUint64(blockNum)
	callOpts := &bind.CallOpts{Context: ctx, BlockNumber: blockBig}

	// Fetch the block header once to stamp every snapshot with the block time.
	header, err := i.ethClient.HeaderByNumber(ctx, blockBig)
	if err != nil {
		return fmt.Errorf("failed to get header for block %d: %w", blockNum, err)
	}
	blockTime := time.Unix(int64(header.Time), 0).UTC()

	for _, quorumID := range quorumNumbers {
		apk, err := i.blsApkRegistry.GetApk(callOpts, quorumID)
		if err != nil {
			return fmt.Errorf("failed to read APK for quorum %d at block %d: %w", quorumID, blockNum, err)
		}

		// callOpts.BlockNumber pins this view call to blockNum, so "current"
		// total stake here means the total stake as of that block.
		totalStake, err := i.stakeRegistry.GetCurrentTotalStake(callOpts, quorumID)
		if err != nil {
			return fmt.Errorf("failed to read total stake for quorum %d at block %d: %w", quorumID, blockNum, err)
		}

		snapshot := &types.QuorumAPK{
			QuorumID:    quorumID,
			BlockNumber: blockNum,
			APK:         core.NewG1Point(apk.X, apk.Y),
			TotalStake:  totalStake,
			UpdatedAt:   blockTime,
		}

		if err := i.store.SaveQuorumAPK(ctx, snapshot); err != nil {
			return fmt.Errorf("failed to save APK snapshot for quorum %d at block %d: %w", quorumID, blockNum, err)
		}

		i.logger.Debug("Indexed quorum APK snapshot", "quorum", quorumID, "block", blockNum)
	}

	return nil
}

// GetStore returns the underlying store (useful for API server).
func (i *Indexer) GetStore() store.Store {
	return i.store
}
