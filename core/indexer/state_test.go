package indexer_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/config"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/indexer/inmem"
	"github.com/Layr-Labs/eigenda/indexer/leveldb"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	indexedstate "github.com/Layr-Labs/eigenda/core/indexer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gethcommon "github.com/ethereum/go-ethereum/common"
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

func mustRegisterOperators(env *config.ConfigLock, logger common.Logger) {

	for _, op := range env.Envs.Operators {
		tx := mustMakeOperatorTransactor(&env.Config, op, logger)

		keyPair, err := core.MakeKeyPairFromString(op.NODE_TEST_PRIVATE_BLS)
		Expect(err).To(BeNil())

		err = tx.RegisterBLSPublicKey(context.Background(), keyPair)
		Expect(err).To(BeNil())

		socket := fmt.Sprintf("%v:%v", op.NODE_HOSTNAME, op.NODE_DISPERSAL_PORT)

		err = tx.RegisterOperator(context.Background(), keyPair.GetPubKeyG1(), socket, quorums)
		Expect(err).To(BeNil())
	}
}

func mustMakeOperatorTransactor(env *config.Config, op config.OperatorVars, logger common.Logger) core.Transactor {

	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	Expect(ok).To(BeTrue())

	pk, _ := strings.CutPrefix(op.NODE_PRIVATE_KEY, "0x")
	config := geth.EthClientConfig{
		RPCURL:           deployer.RPC,
		PrivateKeyString: pk,
	}

	c, err := geth.NewClient(config, logger)
	Expect(err).ToNot(HaveOccurred())

	tx, err := eth.NewTransactor(logger, c, op.NODE_BLS_OPERATOR_STATE_RETRIVER, op.NODE_EIGENDA_SERVICE_MANAGER)
	Expect(err).To(BeNil())
	return tx

}

func mustMakeTestClients(env *config.Config, privateKey string, logger common.Logger) (common.EthClient, common.RPCEthClient) {

	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	Expect(ok).To(BeTrue())

	config := geth.EthClientConfig{
		RPCURL:           deployer.RPC,
		PrivateKeyString: privateKey,
	}

	client, err := geth.NewClient(config, logger)
	if err != nil {
		panic(err)
	}

	rpcClient, err := rpc.Dial(deployer.RPC)
	if err != nil {
		panic(err)
	}

	return client, rpcClient

}

func mustMakeChainState(lock *config.ConfigLock, store indexer.HeaderStore, logger common.Logger) *indexedstate.IndexedChainState {
	client, rpcClient := mustMakeTestClients(&lock.Config, lock.Envs.Batcher.BATCHER_PRIVATE_KEY, logger)

	tx, err := eth.NewTransactor(logger, client, lock.Config.EigenDA.OperatorStateRetreiver, lock.Config.EigenDA.ServiceManager)
	Expect(err).ToNot(HaveOccurred())
	cs := eth.NewChainState(tx, client)

	addr := gethcommon.HexToAddress(lock.Config.EigenDA.ServiceManager)

	indexerConfig := &indexer.Config{
		PullInterval: 1 * time.Second,
	}
	chainState, err := indexedstate.NewIndexedChainState(indexerConfig, addr, cs, store, client, rpcClient, logger)
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

			logger, err := logging.GetLogger(logging.Config{
				StdLevel:  "debug",
				FileLevel: "debug",
			})
			Expect(err).ToNot(HaveOccurred())

			ctx, cancel := context.WithCancel(context.Background())
			_ = cancel

			var (
				store indexer.HeaderStore
			)
			if headerStoreType == "leveldb" {
				dbPath := filepath.Join(lock.Path, "db")
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

			chainState := mustMakeChainState(lock, store, logger)
			err = chainState.Indexer.Index(ctx)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(1 * time.Second)

			mustRegisterOperators(lock, logger)

			time.Sleep(1 * time.Second)

			obj, _, err := chainState.Indexer.HeaderStore.GetLatestObject(chainState.Indexer.Handlers[0].Acc, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(obj).NotTo(BeNil())

			pubKeys, ok := obj.(*indexedstate.OperatorPubKeys)
			Expect(ok).To(BeTrue())
			Expect(pubKeys.Operators).To(HaveLen(len(lock.Operators)))

			obj, header, err := chainState.Indexer.HeaderStore.GetLatestObject(chainState.Indexer.Handlers[1].Acc, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(obj).NotTo(BeNil())

			sockets, ok := obj.(indexedstate.OperatorSockets)
			Expect(ok).To(BeTrue())
			Expect(sockets).To(HaveLen(len(lock.Operators)))

			state, err := chainState.GetIndexedOperatorState(ctx, uint(header.Number), quorums)
			Expect(err).ToNot(HaveOccurred())

			Expect(state.IndexedOperators).To(HaveLen(len(lock.Operators)))

			// TODO: add further tests

		})

	})
})
