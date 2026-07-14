package chainstate

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	ejectionmanager "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// fakeEthClient satisfies IndexerEthClient with a fixed head and records the
// block ranges requested via FilterLogs, returning no logs so that
// indexNewBlocks exercises only the range arithmetic and catch-up loop.
type fakeEthClient struct {
	head   uint64
	ranges [][2]uint64
}

var _ IndexerEthClient = (*fakeEthClient)(nil)

func (f *fakeEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	return f.head, nil
}

func (f *fakeEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]gethtypes.Log, error) {
	f.ranges = append(f.ranges, [2]uint64{q.FromBlock.Uint64(), q.ToBlock.Uint64()})
	return nil, nil
}

func (f *fakeEthClient) CodeAt(
	ctx context.Context,
	contract gethcommon.Address,
	blockNumber *big.Int,
) ([]byte, error) {
	return nil, nil
}

func (f *fakeEthClient) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	return nil, nil
}

func (f *fakeEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*gethtypes.Header, error) {
	return &gethtypes.Header{Number: number}, nil
}

func (f *fakeEthClient) SubscribeFilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
	ch chan<- gethtypes.Log,
) (ethereum.Subscription, error) {
	return nil, nil
}

func (f *fakeEthClient) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	return nil, nil
}

func (f *fakeEthClient) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	return 0, nil
}

func (f *fakeEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (f *fakeEthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (f *fakeEthClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return 0, nil
}

func (f *fakeEthClient) SendTransaction(ctx context.Context, tx *gethtypes.Transaction) error {
	return nil
}

// newRangeTestIndexer wires an Indexer whose contract bindings all point at
// the fake client, so FilterOperatorRegistered etc. call fakeEthClient.
// FilterLogs and return empty iterators.
func newRangeTestIndexer(t *testing.T, client *fakeEthClient, startBlock, batchSize uint64) *Indexer {
	t.Helper()

	cfg := DefaultIndexerConfig()
	cfg.StartBlockNumber = startBlock
	cfg.BlockBatchSize = batchSize

	// The binding constructors only parse the ABI; no RPC happens until a
	// Filter* call, which hits the fake. The contract addresses are irrelevant
	// to the fake, so zero addresses suffice.
	regCoord, err := regcoordinator.NewContractEigenDARegistryCoordinator(gethcommon.Address{}, client)
	require.NoError(t, err)
	blsReg, err := blsapkregistry.NewContractBLSApkRegistry(gethcommon.Address{}, client)
	require.NoError(t, err)
	ejMgr, err := ejectionmanager.NewContractEjectionManager(gethcommon.Address{}, client)
	require.NoError(t, err)

	return &Indexer{
		config:              cfg,
		store:               store.NewMemoryStore(),
		ethClient:           client,
		registryCoordinator: regCoord,
		blsApkRegistry:      blsReg,
		ejectionManager:     ejMgr,
		logger:              test.GetLogger(),
	}
}

func TestFreshStoreIndexesStartBlockItself(t *testing.T) {
	client := &fakeEthClient{head: 2500}
	indexer := newRangeTestIndexer(t, client, 1000, 1000)

	require.NoError(t, indexer.indexNewBlocks(context.Background()))

	require.NotEmpty(t, client.ranges)
	require.Equal(t, uint64(1000), client.ranges[0][0],
		"the configured start block itself must be indexed (inclusive)")

	last, err := indexer.store.GetLastIndexedBlock(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2500), last, "one call must catch up to the observed head")
}

func TestFreshStoreWithZeroStartBlockStartsAtHead(t *testing.T) {
	client := &fakeEthClient{head: 500}
	indexer := newRangeTestIndexer(t, client, 0, 1000)

	require.NoError(t, indexer.indexNewBlocks(context.Background()))

	require.NotEmpty(t, client.ranges)
	require.Equal(t, [2]uint64{500, 500}, client.ranges[0],
		"StartBlockNumber=0 must begin at the current head")
}

func TestResumeFromLastIndexedBlock(t *testing.T) {
	client := &fakeEthClient{head: 300}
	indexer := newRangeTestIndexer(t, client, 1, 1000)
	require.NoError(t, indexer.store.SetLastIndexedBlock(context.Background(), 200))

	require.NoError(t, indexer.indexNewBlocks(context.Background()))

	require.NotEmpty(t, client.ranges)
	require.Equal(t, [2]uint64{201, 300}, client.ranges[0],
		"an existing store must resume at lastIndexed+1, ignoring StartBlockNumber")
}

func TestCatchUpLoopBatches(t *testing.T) {
	client := &fakeEthClient{head: 250}
	indexer := newRangeTestIndexer(t, client, 1, 100)

	require.NoError(t, indexer.indexNewBlocks(context.Background()))

	// Each contract filters the same ranges; assert on the distinct ranges in
	// order of first appearance.
	var distinct [][2]uint64
	for _, r := range client.ranges {
		if len(distinct) == 0 || distinct[len(distinct)-1] != r {
			distinct = append(distinct, r)
		}
	}
	require.Equal(t, [][2]uint64{{1, 100}, {101, 200}, {201, 250}}, distinct,
		"a single call must process batch after batch until the head")

	last, err := indexer.store.GetLastIndexedBlock(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(250), last)
}

func TestAlreadyCaughtUpIsNoOp(t *testing.T) {
	client := &fakeEthClient{head: 100}
	indexer := newRangeTestIndexer(t, client, 1, 1000)
	require.NoError(t, indexer.store.SetLastIndexedBlock(context.Background(), 100))

	require.NoError(t, indexer.indexNewBlocks(context.Background()))
	require.Empty(t, client.ranges, "caught-up indexer must not filter any range")
}
