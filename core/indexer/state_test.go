package indexer_test

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/indexer/inmem"
	"github.com/Layr-Labs/eigenda/indexer/leveldb"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	indexedstate "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/inabox/deploy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	headerStoreType string
)

var (
	quorums []core.QuorumID = []core.QuorumID{0}
)

// Get the location of the test folder from the flag
func init() {
	flag.StringVar(&headerStoreType, "headerStore", "leveldb",
		"The header store implementation to be used (inmem, leveldb)")
}

func mustRegisterOperators(env *deploy.Config, logger logging.Logger) {

	for _, op := range env.Operators {
		tx := mustMakeOperatorTransactor(env, op, logger)

		keyPair, err := core.MakeKeyPairFromString(op.NODE_TEST_PRIVATE_BLS)
		Expect(err).To(BeNil())

		socket := fmt.Sprintf("%v:%v", op.NODE_HOSTNAME, op.NODE_DISPERSAL_PORT)

		salt := [32]byte{}
		_, err = rand.Read(salt[:])
		Expect(err).To(BeNil())

		expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())
		privKey, err := crypto.HexToECDSA(op.NODE_PRIVATE_KEY)
		Expect(err).To(BeNil())

		err = tx.RegisterOperator(context.Background(), keyPair, socket, quorums, privKey, salt, expiry)
		Expect(err).To(BeNil())
	}
}

func mustMakeOperatorTransactor(env *deploy.Config, op deploy.OperatorVars, logger logging.Logger) core.Writer {

	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	Expect(ok).To(BeTrue())

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: op.NODE_PRIVATE_KEY,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	c, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	Expect(err).ToNot(HaveOccurred())

	tx, err := eth.NewWriter(logger, c, op.NODE_BLS_OPERATOR_STATE_RETRIVER, op.NODE_EIGENDA_SERVICE_MANAGER)
	Expect(err).To(BeNil())
	return tx

}

func mustMakeTestClients(env *deploy.Config, privateKey string, logger logging.Logger) (common.EthClient, common.RPCEthClient) {

	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	Expect(ok).To(BeTrue())

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	client, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	if err != nil {
		panic(err)
	}

	rpcClient, err := rpc.Dial(deployer.RPC)
	if err != nil {
		panic(err)
	}

	return client, rpcClient

}

func mustMakeChainState(env *deploy.Config, store indexer.HeaderStore, logger logging.Logger) *indexedstate.IndexedChainState {
	client, rpcClient := mustMakeTestClients(env, env.Batcher[0].BATCHER_PRIVATE_KEY, logger)

	tx, err := eth.NewWriter(logger, client, env.EigenDA.OperatorStateRetreiver, env.EigenDA.ServiceManager)
	Expect(err).ToNot(HaveOccurred())

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
	Expect(err).ToNot(HaveOccurred())

	chainState, err := indexedstate.NewIndexedChainState(cs, indexer)
	if err != nil {
		panic(err)
	}
	return chainState
}

var _ = Describe("Indexer", func() {

	Context("when indexing a chain state", func() {

		It("should index the chain state", func() {

			if testName == "" {
				Skip("No test path provided")
			}

			logger := logging.NewNoopLogger()
			ctx, cancel := context.WithCancel(context.Background())
			_ = cancel

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

			Expect(err).ToNot(HaveOccurred())

			chainState := mustMakeChainState(testConfig, store, logger)
			err = chainState.Indexer.Index(ctx)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(1 * time.Second)

			mustRegisterOperators(testConfig, logger)

			time.Sleep(1 * time.Second)
			lastHeader, err := chainState.Indexer.GetLatestHeader(false)
			Expect(err).ToNot(HaveOccurred())
			obj, err := chainState.Indexer.GetObject(lastHeader, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(obj).NotTo(BeNil())

			pubKeys, ok := obj.(*indexedstate.OperatorPubKeys)
			Expect(ok).To(BeTrue())
			Expect(pubKeys.Operators).To(HaveLen(len(testConfig.Operators)))

			obj, err = chainState.Indexer.GetObject(lastHeader, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(obj).NotTo(BeNil())

			sockets, ok := obj.(indexedstate.OperatorSockets)
			Expect(ok).To(BeTrue())
			Expect(sockets).To(HaveLen(len(testConfig.Operators)))

			header, err := chainState.Indexer.GetLatestHeader(false)
			Expect(err).ToNot(HaveOccurred())
			state, err := chainState.GetIndexedOperatorState(ctx, uint(header.Number), quorums)
			Expect(err).ToNot(HaveOccurred())

			Expect(state.IndexedOperators).To(HaveLen(len(testConfig.Operators)))

			// TODO: add further tests

		})

	})
})
