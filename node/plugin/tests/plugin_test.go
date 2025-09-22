package test

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/node/plugin"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	logger = test.GetLogger()

	// Shared test resources
	anvilContainer *testbed.AnvilContainer
	deployResult   *testbed.DeploymentResult
	privateKeys    *testbed.PrivateKeyMaps
	testOperator   OperatorConfig
)

// OperatorConfig holds configuration for a test operator
type OperatorConfig struct {
	NODE_HOSTNAME                    string
	NODE_DISPERSAL_PORT              string
	NODE_RETRIEVAL_PORT              string
	NODE_V2_DISPERSAL_PORT           string
	NODE_V2_RETRIEVAL_PORT           string
	NODE_ECDSA_KEY_FILE              string
	NODE_BLS_KEY_FILE                string
	NODE_ECDSA_KEY_PASSWORD          string
	NODE_BLS_KEY_PASSWORD            string
	NODE_QUORUM_ID_LIST              string
	NODE_CHAIN_RPC                   string
	NODE_EIGENDA_DIRECTORY           string
	NODE_BLS_OPERATOR_STATE_RETRIVER string
	NODE_EIGENDA_SERVICE_MANAGER     string
	NODE_CHURNER_URL                 string
}

// TestMain sets up the test environment once for all tests
func TestMain(m *testing.M) {
	// Parse flags first to initialize the testing framework
	flag.Parse()

	if testing.Short() {
		logger.Info("Skipping plugin integration tests in short mode")
		os.Exit(0)
	}

	setupAndRun(m)
}

func setupAndRun(m *testing.M) {
	ctx := context.Background()

	var err error
	anvilContainer, err = testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		Logger:         logger,
	})
	if err != nil {
		logger.Fatal("Failed to start anvil container:", err)
	}

	logger.Info("Loading private keys")
	privateKeys, err = testbed.LoadPrivateKeys(testbed.LoadPrivateKeysInput{
		NumOperators: 1,
		NumRelays:    0,
	})
	if err != nil {
		logger.Fatal("Failed to load private keys:", err)
	}

	logger.Info("Deploying contracts")
	// Get deployer key from testbed
	deployerKey, _ := testbed.GetAnvilDefaultKeys()

	// Deploy contracts using testbed
	deployConfig := testbed.DeploymentConfig{
		AnvilRPCURL:      "http://localhost:8545",
		DeployerKey:      deployerKey,
		NumOperators:     1,
		NumRelays:        0,
		MaxOperatorCount: 10,
		PrivateKeys:      privateKeys,
		Logger:           logger,
	}

	deployResult, err = testbed.DeployEigenDAContracts(deployConfig)
	if err != nil {
		logger.Fatal("Failed to deploy contracts:", err)
	}

	logger.Info("Setting up test operators")
	setupTestOperators()

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

func setupTestOperators() {
	// Create operator configurations using testbed keys
	opName := "opr0"
	operator := OperatorConfig{
		NODE_HOSTNAME:                    "localhost",
		NODE_DISPERSAL_PORT:              "32003",
		NODE_RETRIEVAL_PORT:              "32004",
		NODE_V2_DISPERSAL_PORT:           "32005",
		NODE_V2_RETRIEVAL_PORT:           "32006",
		NODE_ECDSA_KEY_FILE:              privateKeys.EcdsaMap[opName].KeyFile,
		NODE_BLS_KEY_FILE:                privateKeys.BlsMap[opName].KeyFile,
		NODE_ECDSA_KEY_PASSWORD:          privateKeys.EcdsaMap[opName].Password,
		NODE_BLS_KEY_PASSWORD:            privateKeys.BlsMap[opName].Password,
		NODE_QUORUM_ID_LIST:              "0,1",
		NODE_CHAIN_RPC:                   "http://localhost:8545",
		NODE_EIGENDA_DIRECTORY:           deployResult.EigenDA.EigenDADirectory,
		NODE_BLS_OPERATOR_STATE_RETRIVER: deployResult.EigenDA.OperatorStateRetriever,
		NODE_EIGENDA_SERVICE_MANAGER:     deployResult.EigenDA.ServiceManager,
		NODE_CHURNER_URL:                 "",
	}
	testOperator = operator
}

func TestPluginOptIn(t *testing.T) {
	ctx := t.Context()

	require.NotEmpty(t, testOperator.NODE_QUORUM_ID_LIST)

	runNodePlugin(t, "opt-out", testOperator)

	tx := getTransactor(t, testOperator)
	operatorID := getOperatorId(t, testOperator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)

	runNodePlugin(t, "opt-in", testOperator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 2, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(1), ids)
}

func TestPluginOptInAndOptOut(t *testing.T) {
	ctx := t.Context()

	require.NotEmpty(t, testOperator.NODE_QUORUM_ID_LIST)

	runNodePlugin(t, "opt-out", testOperator)

	tx := getTransactor(t, testOperator)
	operatorID := getOperatorId(t, testOperator)

	runNodePlugin(t, "opt-in", testOperator)
	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 2, len(registeredQuorumIds))
	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(1), ids)

	runNodePlugin(t, "opt-out", testOperator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)
}

func TestPluginOptInAndQuorumUpdate(t *testing.T) {
	ctx := t.Context()

	require.Equal(t, "0,1", testOperator.NODE_QUORUM_ID_LIST)

	runNodePlugin(t, "opt-out", testOperator)

	tx := getTransactor(t, testOperator)
	operatorID := getOperatorId(t, testOperator)

	runNodePlugin(t, "opt-in", testOperator)

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

	require.Equal(t, "0,1", testOperator.NODE_QUORUM_ID_LIST)

	runNodePlugin(t, "opt-out", testOperator)

	tx := getTransactor(t, testOperator)
	operatorID := getOperatorId(t, testOperator)

	registeredQuorumIds, err := tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err := tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)

	runNodePlugin(t, "invalid", testOperator)

	registeredQuorumIds, err = tx.GetRegisteredQuorumIdsForOperator(ctx, operatorID)
	require.NoError(t, err)
	require.Equal(t, 0, len(registeredQuorumIds))

	ids, err = tx.GetNumberOfRegisteredOperatorForQuorum(ctx, core.QuorumID(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), ids)
}

func getOperatorId(t *testing.T, operator OperatorConfig) [32]byte {
	t.Helper()

	_, privateKey, err := plugin.GetECDSAPrivateKey(operator.NODE_ECDSA_KEY_FILE, operator.NODE_ECDSA_KEY_PASSWORD)
	require.NoError(t, err)
	require.NotNil(t, privateKey)
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

func getTransactor(t *testing.T, operator OperatorConfig) *eth.Writer {
	t.Helper()

	// Use deployer key from testbed
	deployerKey, _ := testbed.GetAnvilDefaultKeys()
	hexPk := strings.TrimPrefix(deployerKey, "0x")
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

// runNodePlugin runs the node plugin directly using go run
func runNodePlugin(t *testing.T, operation string, operator OperatorConfig) {
	t.Helper()

	ctx := t.Context()
	socket := string(core.MakeOperatorSocket(
		operator.NODE_HOSTNAME,
		operator.NODE_DISPERSAL_PORT,
		operator.NODE_RETRIEVAL_PORT,
		operator.NODE_V2_DISPERSAL_PORT,
		operator.NODE_V2_RETRIEVAL_PORT,
	))

	// Get the path to the node plugin cmd directory relative to this file
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	rootDir := filepath.Join(testDir, "..", "..", "..")
	pluginCmdPath := filepath.Join(rootDir, "node", "plugin", "cmd")

	// Run the plugin directly with go run
	cmd := exec.CommandContext(ctx, "go", "run", pluginCmdPath,
		"--operation", operation,
		"--ecdsa-key-file", operator.NODE_ECDSA_KEY_FILE,
		"--bls-key-file", operator.NODE_BLS_KEY_FILE,
		"--ecdsa-key-password", operator.NODE_ECDSA_KEY_PASSWORD,
		"--bls-key-password", operator.NODE_BLS_KEY_PASSWORD,
		"--socket", socket,
		"--quorum-id-list", operator.NODE_QUORUM_ID_LIST,
		"--chain-rpc", operator.NODE_CHAIN_RPC,
		"--eigenda-directory", operator.NODE_EIGENDA_DIRECTORY,
		"--bls-operator-state-retriever", operator.NODE_BLS_OPERATOR_STATE_RETRIVER,
		"--eigenda-service-manager", operator.NODE_EIGENDA_SERVICE_MANAGER,
		"--churner-url", operator.NODE_CHURNER_URL,
		"--num-confirmations", "0",
	)

	logger.Info("Running node plugin", "operation", operation)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Fatalf("failed to run node plugin: %v, output: %s", err, string(output))
	}

	logger.Info("Node plugin executed successfully", "output", string(output))
}
