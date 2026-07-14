package chainstate

import (
	"context"
	"math/big"
	"sort"
	"testing"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	ejectionmanager "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// newTestIndexer builds an Indexer with only the pieces the event handlers
// touch (store and logger). The contract bindings and eth client stay nil:
// handler unit tests apply decoded events directly, so no RPC is involved.
func newTestIndexer(t *testing.T) *Indexer {
	t.Helper()
	return &Indexer{
		config: DefaultIndexerConfig(),
		store:  store.NewMemoryStore(),
		logger: test.GetLogger(),
	}
}

func registeredEvent(
	id core.OperatorID,
	addr gethcommon.Address,
	block uint64,
	logIndex uint,
) *regcoordinator.ContractEigenDARegistryCoordinatorOperatorRegistered {
	return &regcoordinator.ContractEigenDARegistryCoordinatorOperatorRegistered{
		Operator:   addr,
		OperatorId: id,
		Raw:        gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

func deregisteredEvent(
	id core.OperatorID,
	block uint64,
	logIndex uint,
) *regcoordinator.ContractEigenDARegistryCoordinatorOperatorDeregistered {
	return &regcoordinator.ContractEigenDARegistryCoordinatorOperatorDeregistered{
		OperatorId: id,
		Raw:        gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

func socketEvent(
	id core.OperatorID,
	socket string,
	block uint64,
	logIndex uint,
) *regcoordinator.ContractEigenDARegistryCoordinatorOperatorSocketUpdate {
	return &regcoordinator.ContractEigenDARegistryCoordinatorOperatorSocketUpdate{
		OperatorId: id,
		Socket:     socket,
		Raw:        gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

func addedToQuorumsEvent(
	id core.OperatorID,
	quorums []byte,
	block uint64,
	logIndex uint,
) *blsapkregistry.ContractBLSApkRegistryOperatorAddedToQuorums {
	return &blsapkregistry.ContractBLSApkRegistryOperatorAddedToQuorums{
		OperatorId:    id,
		QuorumNumbers: quorums,
		Raw:           gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

func removedFromQuorumsEvent(
	id core.OperatorID,
	quorums []byte,
	block uint64,
	logIndex uint,
) *blsapkregistry.ContractBLSApkRegistryOperatorRemovedFromQuorums {
	return &blsapkregistry.ContractBLSApkRegistryOperatorRemovedFromQuorums{
		OperatorId:    id,
		QuorumNumbers: quorums,
		Raw:           gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

// pubkeyEvent builds a NewPubkeyRegistration event from a generated key pair,
// serializing the G2 coordinates in the contract's [A1, A0] component order
// (the reverse of gnark-crypto's), as the real event data arrives.
func pubkeyEvent(
	t *testing.T,
	keyPair *core.KeyPair,
	addr gethcommon.Address,
	block uint64,
	logIndex uint,
) *blsapkregistry.ContractBLSApkRegistryNewPubkeyRegistration {
	t.Helper()
	g1 := keyPair.GetPubKeyG1()
	g2 := keyPair.GetPubKeyG2()
	return &blsapkregistry.ContractBLSApkRegistryNewPubkeyRegistration{
		Operator: addr,
		PubkeyG1: blsapkregistry.BN254G1Point{
			X: g1.X.BigInt(new(big.Int)),
			Y: g1.Y.BigInt(new(big.Int)),
		},
		PubkeyG2: blsapkregistry.BN254G2Point{
			X: [2]*big.Int{g2.X.A1.BigInt(new(big.Int)), g2.X.A0.BigInt(new(big.Int))},
			Y: [2]*big.Int{g2.Y.A1.BigInt(new(big.Int)), g2.Y.A0.BigInt(new(big.Int))},
		},
		Raw: gethtypes.Log{BlockNumber: block, Index: logIndex},
	}
}

func TestDeregThenReregWithinBatch(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)
	id := core.OperatorID{1}
	addr := gethcommon.HexToAddress("0x1234")

	// Prior state: registered at block 10.
	require.NoError(t, indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 10, 0)))

	// Both events land in one indexed range. Chronological order: dereg@100
	// then re-reg@150. The final state must be registered.
	events := []chainEvent{
		{
			blockNumber: 100, logIndex: 0,
			apply: func(ctx context.Context) error {
				return indexer.handleOperatorDeregistered(ctx, deregisteredEvent(id, 100, 0))
			},
		},
		{
			blockNumber: 150, logIndex: 0,
			apply: func(ctx context.Context) error {
				return indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 150, 0))
			},
		},
	}
	applySorted(t, ctx, events)

	op, err := indexer.store.GetOperator(ctx, id)
	require.NoError(t, err)
	require.Nil(t, op.DeregisteredAtBlockNumber, "re-registration must clear the deregistration marker")
	require.Equal(t, uint64(150), op.RegisteredAtBlockNumber)
}

func TestQuorumRemoveThenReAddWithinBatch(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)
	id := core.OperatorID{1}
	addr := gethcommon.HexToAddress("0x1234")

	require.NoError(t, indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 10, 0)))
	require.NoError(t, indexer.handleOperatorAddedToQuorums(ctx, addedToQuorumsEvent(id, []byte{0, 1}, 10, 1)))

	events := []chainEvent{
		{
			blockNumber: 100, logIndex: 0,
			apply: func(ctx context.Context) error {
				return indexer.handleOperatorRemovedFromQuorums(ctx, removedFromQuorumsEvent(id, []byte{0}, 100, 0))
			},
		},
		{
			blockNumber: 150, logIndex: 0,
			apply: func(ctx context.Context) error {
				return indexer.handleOperatorAddedToQuorums(ctx, addedToQuorumsEvent(id, []byte{0}, 150, 0))
			},
		},
	}
	applySorted(t, ctx, events)

	op, err := indexer.store.GetOperator(ctx, id)
	require.NoError(t, err)
	require.ElementsMatch(t, []core.QuorumID{0, 1}, op.QuorumIDs,
		"re-added quorum must survive a remove earlier in the same batch")
}

// TestFirstRegistrationIntraTxOrder replays a registration transaction's
// events in the order the contracts actually emit them: NewPubkeyRegistration
// and OperatorSocketUpdate BEFORE OperatorRegistered (see
// RegistryCoordinator._registerOperator and _getOrCreateOperatorId). The
// skeleton records created by the earlier events must merge into one complete
// operator.
func TestFirstRegistrationIntraTxOrder(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)
	addr := gethcommon.HexToAddress("0x1234")

	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	id := keyPair.GetPubKeyG1().GetOperatorID()

	require.NoError(t, indexer.handleNewPubkeyRegistration(ctx, pubkeyEvent(t, keyPair, addr, 100, 5)))
	require.NoError(t, indexer.handleOperatorSocketUpdate(ctx, socketEvent(id, "host:1", 100, 7)))
	require.NoError(t, indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 100, 9)))
	require.NoError(t, indexer.handleOperatorAddedToQuorums(ctx, addedToQuorumsEvent(id, []byte{0}, 100, 11)))

	op, err := indexer.store.GetOperator(ctx, id)
	require.NoError(t, err)
	require.Equal(t, addr, op.Address)
	require.Equal(t, "host:1", op.Socket)
	require.Equal(t, uint64(100), op.RegisteredAtBlockNumber)
	require.NotNil(t, op.BLSPubKeyG1)
	require.NotNil(t, op.BLSPubKeyG2)
	require.Equal(t, []core.QuorumID{0}, op.QuorumIDs)
	require.Nil(t, op.DeregisteredAtBlockNumber)

	// The locally computed operator ID must match the contract's definition
	// (keccak of the G1 pubkey), and the G2 component swap must yield the key
	// actually paired with G1.
	require.True(t, op.BLSPubKeyG1.G1Affine.Equal(keyPair.GetPubKeyG1().G1Affine))
	ok, err := op.BLSPubKeyG1.VerifyEquivalence(op.BLSPubKeyG2)
	require.NoError(t, err)
	require.True(t, ok, "stored G2 key must be the pair of the G1 key (component swap correctness)")
}

func TestReregistrationPreservesBLSKeys(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)
	addr := gethcommon.HexToAddress("0x1234")

	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	id := keyPair.GetPubKeyG1().GetOperatorID()

	// First registration (pubkey event fires once, ever).
	require.NoError(t, indexer.handleNewPubkeyRegistration(ctx, pubkeyEvent(t, keyPair, addr, 100, 0)))
	require.NoError(t, indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 100, 1)))

	// Deregister, then re-register in a later range. No second pubkey event.
	require.NoError(t, indexer.handleOperatorDeregistered(ctx, deregisteredEvent(id, 200, 0)))
	require.NoError(t, indexer.handleOperatorRegistered(ctx, registeredEvent(id, addr, 300, 0)))

	op, err := indexer.store.GetOperator(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, op.BLSPubKeyG1, "re-registration must not wipe BLS keys")
	require.NotNil(t, op.BLSPubKeyG2)
	require.Nil(t, op.DeregisteredAtBlockNumber)
	require.Equal(t, uint64(300), op.RegisteredAtBlockNumber)
}

func TestDeregistrationForUnknownOperatorFailsBatch(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)

	err := indexer.handleOperatorDeregistered(ctx, deregisteredEvent(core.OperatorID{1}, 100, 0))
	require.Error(t, err, "unknown operator signals an indexing gap; the batch must fail for retry")

	err = indexer.handleOperatorAddedToQuorums(ctx, addedToQuorumsEvent(core.OperatorID{1}, []byte{0}, 100, 0))
	require.Error(t, err)

	err = indexer.handleOperatorRemovedFromQuorums(ctx, removedFromQuorumsEvent(core.OperatorID{1}, []byte{0}, 100, 0))
	require.Error(t, err)
}

func TestChainEventSortOrder(t *testing.T) {
	// Same comparator as indexBlockRange: block number, then log index.
	events := []chainEvent{
		{blockNumber: 150, logIndex: 0},
		{blockNumber: 100, logIndex: 9},
		{blockNumber: 100, logIndex: 5},
		{blockNumber: 100, logIndex: 7},
	}
	sortChainEvents(events)

	var got [][2]uint64
	for _, e := range events {
		got = append(got, [2]uint64{e.blockNumber, uint64(e.logIndex)})
	}
	require.Equal(t, [][2]uint64{{100, 5}, {100, 7}, {100, 9}, {150, 0}}, got)
}

// applySorted sorts the events with the same comparator indexBlockRange uses
// and applies them, failing the test on any handler error.
func applySorted(t *testing.T, ctx context.Context, events []chainEvent) {
	t.Helper()
	sortChainEvents(events)
	for _, e := range events {
		require.NoError(t, e.apply(ctx))
	}
}

func sortChainEvents(events []chainEvent) {
	sort.Slice(events, func(a, b int) bool {
		if events[a].blockNumber != events[b].blockNumber {
			return events[a].blockNumber < events[b].blockNumber
		}
		return events[a].logIndex < events[b].logIndex
	})
}

func TestMultiQuorumEjectionRecordsAllQuorums(t *testing.T) {
	ctx := context.Background()
	indexer := newTestIndexer(t)
	id := core.OperatorID{1}
	tx := gethcommon.Hash{0xab}

	// ejectOperators emits one OperatorEjected per quorum with the same tx
	// hash and block, differing only in log index.
	for i, quorum := range []uint8{0, 1} {
		event := &ejectionmanager.ContractEjectionManagerOperatorEjected{
			OperatorId:   id,
			QuorumNumber: quorum,
			Raw:          gethtypes.Log{BlockNumber: 100, Index: uint(i), TxHash: tx},
		}
		require.NoError(t, indexer.handleOperatorEjected(ctx, event))
	}

	ejections, err := indexer.store.ListEjections(ctx, &id, 0, 0)
	require.NoError(t, err)
	require.Len(t, ejections, 2, "one record per per-quorum ejection event")
}
