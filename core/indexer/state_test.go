package indexer_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/indexer/inmem"
	"github.com/Layr-Labs/eigenda/indexer/leveldb"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigensdk-go/logging"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

var (
	quorums []core.QuorumID = []core.QuorumID{0}
)

func mustRegisterOperators(t *testing.T, env *deploy.Config, logger logging.Logger) {
	t.Helper()
	for _, op := range env.Operators {
		tx := mustMakeOperatorTransactor(t, env, op, logger)

		signer, err := blssigner.NewSigner(blssignerTypes.SignerConfig{
			PrivateKey: op.NODE_TEST_PRIVATE_BLS,
			SignerType: blssignerTypes.PrivateKey,
		})
		require.NoError(t, err, "failed to create signer")

		socket := fmt.Sprintf("%v:%v", op.NODE_HOSTNAME, op.NODE_DISPERSAL_PORT)

		salt := [32]byte{}
		_, err = rand.Read(salt[:])
		require.NoError(t, err, "failed to generate salt")

		expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())
		privKey, err := crypto.HexToECDSA(op.NODE_PRIVATE_KEY)
		require.NoError(t, err, "failed to parse private key")

		err = tx.RegisterOperator(context.Background(), signer, socket, quorums, privKey, salt, expiry)
		require.NoError(t, err, "failed to register operator")
	}
}

func mustMakeOperatorTransactor(t *testing.T, env *deploy.Config, op deploy.OperatorVars, logger logging.Logger) core.Writer {
	t.Helper()
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	require.True(t, ok, "deployer not found")

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: op.NODE_PRIVATE_KEY,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	c, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	require.NoError(t, err, "failed to create geth client")

	tx, err := eth.NewWriter(logger, c, op.NODE_BLS_OPERATOR_STATE_RETRIVER, op.NODE_EIGENDA_SERVICE_MANAGER)
	require.NoError(t, err, "failed to create writer")
	return tx
}

func mustMakeTestClients(t *testing.T, env *deploy.Config, privateKey string, logger logging.Logger) (common.EthClient, common.RPCEthClient) {
	t.Helper()
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	require.True(t, ok, "deployer not found")

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	client, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	require.NoError(t, err, "failed to create geth client")

	rpcClient, err := rpc.Dial(deployer.RPC)
	require.NoError(t, err, "failed to create RPC client")

	return client, rpcClient
}

func mustMakeChainState(t *testing.T, env *deploy.Config, _ indexer.HeaderStore, logger logging.Logger) *coreindexer.IndexedChainState {
	t.Helper()
	client, rpcClient := mustMakeTestClients(t, env, env.Batcher[0].BATCHER_PRIVATE_KEY, logger)

	tx, err := eth.NewWriter(logger, client, env.EigenDA.OperatorStateRetriever, env.EigenDA.ServiceManager)
	require.NoError(t, err, "failed to create writer")

	var (
		cs            = eth.NewChainState(tx, client)
		indexerConfig = indexer.Config{
			PullInterval: 1 * time.Second,
		}
	)

	indexer, err := coreindexer.CreateNewIndexer(
		&indexerConfig,
		client,
		rpcClient,
		env.EigenDA.ServiceManager,
		logger,
	)
	require.NoError(t, err, "failed to create indexer")

	chainState, err := coreindexer.NewIndexedChainState(cs, indexer)
	require.NoError(t, err, "failed to create indexed chain state")
	return chainState
}

// This test exercises the core indexer, which is not used in production. Since this test is flaky, disable it.
var skip = true

func TestIndexChainState(t *testing.T) {
	if skip {
		t.Skip("Test disabled - core indexer not used in production")
	}

	if testName == "" {
		t.Skip("No test path provided")
	}

	logger := test.GetLogger()
	ctx := t.Context()

	var (
		store indexer.HeaderStore
		err   error
	)
	if headerStoreType == "leveldb" {
		dbPath := filepath.Join(testConfig.Path, "db")
		s, err := leveldb.NewHeaderStore(dbPath)
		if err == nil {
			defer s.Close()
			defer func() { _ = os.RemoveAll(dbPath) }()
			store = s
		}
	} else {
		store = inmem.NewHeaderStore()
	}

	require.NoError(t, err, "failed to create header store")

	chainState := mustMakeChainState(t, testConfig, store, logger)
	err = chainState.Indexer.Index(ctx)
	require.NoError(t, err, "failed to index")

	time.Sleep(1 * time.Second)

	mustRegisterOperators(t, testConfig, logger)

	time.Sleep(1 * time.Second)
	lastHeader, err := chainState.Indexer.GetLatestHeader(false)
	require.NoError(t, err, "failed to get latest header")
	obj, err := chainState.Indexer.GetObject(lastHeader, 0)
	require.NoError(t, err, "failed to get object at index 0")
	require.NotNil(t, obj, "object should not be nil")

	pubKeys, ok := obj.(*coreindexer.OperatorPubKeys)
	require.True(t, ok, "object should be OperatorPubKeys")
	require.Len(t, pubKeys.Operators, len(testConfig.Operators), "unexpected number of operators")

	obj, err = chainState.Indexer.GetObject(lastHeader, 1)
	require.NoError(t, err, "failed to get object at index 1")
	require.NotNil(t, obj, "object should not be nil")

	sockets, ok := obj.(coreindexer.OperatorSockets)
	require.True(t, ok, "object should be OperatorSockets")
	require.Len(t, sockets, len(testConfig.Operators), "unexpected number of sockets")

	header, err := chainState.Indexer.GetLatestHeader(false)
	require.NoError(t, err, "failed to get latest header")
	state, err := chainState.GetIndexedOperatorState(ctx, uint(header.Number), quorums)
	require.NoError(t, err, "failed to get indexed operator state")

	require.Len(t, state.IndexedOperators, len(testConfig.Operators), "unexpected number of indexed operators")

	// TODO: add further tests
}
