package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/node/plugin"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
}

var (
	testConfig   *deploy.Config
	templateName string
	testName     string
)

func TestMain(m *testing.M) {
	flag.Parse()
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {
	rootPath := "../../../"

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			panic(err)
		}
	}

	testConfig = deploy.NewTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = false

	if testing.Short() {
		fmt.Println("Skipping plugin integration test in short mode")
		os.Exit(0)
		return
	}

	fmt.Println("Starting anvil")
	testConfig.StartAnvil()

	fmt.Println("Deploying experiment")
	testConfig.DeployExperiment()
}

func teardown() {
	if testConfig != nil {
		fmt.Println("Stopping anvil")
		testConfig.StopAnvil()
	}
}

func TestPluginOptIn(t *testing.T) {
	operator := testConfig.Operators[0]
	assert.NotEmpty(t, operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), ids)

	testConfig.RunNodePluginBinary("opt-in", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), ids)
}

func TestPluginOptInAndOptOut(t *testing.T) {
	operator := testConfig.Operators[0]
	assert.NotEmpty(t, operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	testConfig.RunNodePluginBinary("opt-in", operator)
	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), ids)

	testConfig.RunNodePluginBinary("opt-out", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), ids)
}

func TestPluginOptInAndQuorumUpdate(t *testing.T) {
	operator := testConfig.Operators[0]
	assert.Equal(t, "0,1", operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	testConfig.RunNodePluginBinary("opt-in", operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(registeredQuorumIds))
	assert.Equal(t, uint8(0), registeredQuorumIds[0])

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), ids)
}

func TestPluginInvalidOperation(t *testing.T) {
	operator := testConfig.Operators[0]
	assert.Equal(t, "0,1", operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), ids)

	testConfig.RunNodePluginBinary("invalid", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(context.Background(), core.QuorumID(0))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), ids)
}

func getOperatorId(t *testing.T, operator deploy.OperatorVars) [32]byte {
	_, privateKey, err := plugin.GetECDSAPrivateKey(operator.NODE_ECDSA_KEY_FILE, operator.NODE_ECDSA_KEY_PASSWORD)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	assert.NoError(t, err)

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{operator.NODE_CHAIN_RPC},
		PrivateKeyString: *privateKey,
	}

	client, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	transactor, err := eth.NewWriter(logger, client, operator.NODE_BLS_OPERATOR_STATE_RETRIVER, operator.NODE_EIGENDA_SERVICE_MANAGER)
	assert.NoError(t, err)
	assert.NotNil(t, transactor)

	kp, err := bls.ReadPrivateKeyFromFile(operator.NODE_BLS_KEY_FILE, operator.NODE_BLS_KEY_PASSWORD)
	assert.NoError(t, err)
	assert.NotNil(t, kp)

	g1point := &core.G1Point{
		G1Affine: kp.PubKey.G1Affine,
	}
	keyPair := &core.KeyPair{
		PrivKey: kp.PrivKey,
		PubKey:  g1point,
	}

	return keyPair.GetPubKeyG1().GetOperatorID()
}

func getTransactor(t *testing.T, operator deploy.OperatorVars) *eth.Writer {
	hexPk := strings.TrimPrefix(testConfig.Pks.EcdsaMap[testConfig.Deployers[0].Name].PrivateKey, "0x")
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	assert.NoError(t, err)

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{operator.NODE_CHAIN_RPC},
		PrivateKeyString: hexPk,
		NumConfirmations: 0,
	}

	client, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	transactor, err := eth.NewWriter(logger, client, operator.NODE_BLS_OPERATOR_STATE_RETRIVER, operator.NODE_EIGENDA_SERVICE_MANAGER)
	assert.NoError(t, err)
	assert.NotNil(t, transactor)

	return transactor
}
