package test

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/node/plugin"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
}

var (
	templateName string
	testName     string

	logger = testutils.GetLogger()

	// Shared test resources
	anvilContainer *testbed.AnvilContainer
	testConfig     *deploy.Config
)

// TestMain sets up the test environment once for all tests
func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		logger.Info("Skipping plugin integration tests in short mode")
		os.Exit(0)
	}

	setupAndRun(m)
}

func setupAndRun(m *testing.M) {
	ctx := context.Background()
	rootPath := "../../../"

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			logger.Fatal("Failed to create test directory:", err)
		}
	}

	testConfig = deploy.NewTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = false

	logger.Info("Starting anvil")
	var err error
	anvilContainer, err = testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true, // This will bind container port 8545 to host port 8545
		Logger:         logger,
	})
	if err != nil {
		logger.Fatal("Failed to start anvil container:", err)
	}

	logger.Info("Deploying experiment")
	if err := testConfig.DeployExperiment(); err != nil {
		logger.Fatal("Failed to deploy experiment:", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func cleanup() {
	if anvilContainer != nil {
		logger.Info("Stopping anvil")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = anvilContainer.Terminate(ctx)
	}
}

func TestPluginOptIn(t *testing.T) {
	ctx := t.Context()

	operator := testConfig.Operators[0]
	require.NotEmpty(t, operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)

	testConfig.RunNodePluginBinary("opt-in", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 2, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(1), ids)
}

func TestPluginOptInAndOptOut(t *testing.T) {
	ctx := t.Context()

	operator := testConfig.Operators[0]
	require.NotEmpty(t, operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	testConfig.RunNodePluginBinary("opt-in", operator)
	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 2, len(registeredQuorumIds))
	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(1), ids)

	testConfig.RunNodePluginBinary("opt-out", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)
}

func TestPluginOptInAndQuorumUpdate(t *testing.T) {
	ctx := t.Context()

	operator := testConfig.Operators[0]
	require.Equal(t, "0,1", operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	testConfig.RunNodePluginBinary("opt-in", operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 2, len(registeredQuorumIds))
	require.Equal(t, uint8(0), registeredQuorumIds[0])

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(1), ids)
}

func TestPluginInvalidOperation(t *testing.T) {
	ctx := t.Context()

	operator := testConfig.Operators[0]
	require.Equal(t, "0,1", operator.NODE_QUORUM_ID_LIST)

	testConfig.RunNodePluginBinary("opt-out", operator)

	tx := getTransactor(t, operator)
	operatorID := getOperatorId(t, operator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)

	testConfig.RunNodePluginBinary("invalid", operator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)
}

func getOperatorId(t *testing.T, operator deploy.OperatorVars) [32]byte {
	t.Helper()

	_, privateKey, err := plugin.GetECDSAPrivateKey(operator.NODE_ECDSA_KEY_FILE, operator.NODE_ECDSA_KEY_PASSWORD)
	require.NoError(t, err)
	require.NotNil(t, privateKey)
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	require.NoError(t, err)

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{operator.NODE_CHAIN_RPC},
		PrivateKeyString: *privateKey,
	}

	client, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	require.NoError(t, err)
	require.NotNil(t, client)

	transactor, err := eth.NewWriter(
		logger, client, operator.NODE_BLS_OPERATOR_STATE_RETRIVER, operator.NODE_EIGENDA_SERVICE_MANAGER)
	require.NoError(t, err)
	require.NotNil(t, transactor)

	kp, err := bls.ReadPrivateKeyFromFile(operator.NODE_BLS_KEY_FILE, operator.NODE_BLS_KEY_PASSWORD)
	require.NoError(t, err)
	require.NotNil(t, kp)
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
	t.Helper()

	hexPk := strings.TrimPrefix(testConfig.Pks.EcdsaMap[testConfig.Deployers[0].Name].PrivateKey, "0x")
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	require.NoError(t, err)

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{operator.NODE_CHAIN_RPC},
		PrivateKeyString: hexPk,
		NumConfirmations: 0,
	}

	client, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	require.NoError(t, err)
	require.NotNil(t, client)

	transactor, err := eth.NewWriter(
		logger, client, operator.NODE_BLS_OPERATOR_STATE_RETRIVER, operator.NODE_EIGENDA_SERVICE_MANAGER)
	require.NoError(t, err)
	require.NotNil(t, transactor)

	return transactor
}
