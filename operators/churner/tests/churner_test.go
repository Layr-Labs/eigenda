package test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	indexermock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/node/plugin"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Simple operator info struct that only contains what the test needs
type operatorInfo struct {
	ECDSAPrivateKey string
	ECDSAKeyFile    string
	ECDSAPassword   string
	BLSKeyPath      string
	BLSPassword     string
}

var (
	localstackPort                 = "4570"
	rpcURL                         = "http://localhost:8545"
	quorumIds                      = []uint32{0, 1}
	operatorAddr                   = ""
	churnerPrivateKeyHex           = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	operatorToChurnInPrivateKeyHex = "0000000000000000000000000000000000000000000000000000000000000020"
	numRetries                     = 0

	logger = test.GetLogger()
)

// TestSetup contains all the test infrastructure needed for churner tests
type TestSetup struct {
	AnvilContainer      *testbed.AnvilContainer
	LocalstackContainer *testbed.LocalStackContainer
	Contracts           *testbed.DeploymentResult
	Operators           []operatorInfo
	PrivateKeys         *testbed.PrivateKeyMaps
}

func setupTest(t *testing.T) *TestSetup {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping churner test in short mode")
	}

	ctx := t.Context()
	numOperators := 4

	// Start localstack container
	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start localstack container")

	// Start anvil container
	anvilContainer, err := testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start anvil container")

	// Load private keys using testbed
	privateKeys, err := testbed.LoadPrivateKeys(testbed.LoadPrivateKeysInput{
		NumOperators: numOperators,
		NumRelays:    0,
	})
	require.NoError(t, err, "failed to load private keys")

	// Get deployer key from Anvil's default accounts
	deployerKey, _ := testbed.GetAnvilDefaultKeys()

	// Deploy contracts
	logger.Info("Deploying contracts")
	deploymentResult, err := testbed.DeployEigenDAContracts(testbed.DeploymentConfig{
		AnvilRPCURL:      "http://localhost:8545",
		DeployerKey:      deployerKey,
		NumOperators:     numOperators,
		NumRelays:        0,
		MaxOperatorCount: 3, // Set max to 3 so the 4th operator can churn
		Stakes: []testbed.Stakes{
			{Total: 100e18, Distribution: []float32{1, 4, 6, 10}},
			{Total: 100e18, Distribution: []float32{1, 3, 8, 9}},
		},
		PrivateKeys: privateKeys,
		Logger:      logger,
	})
	require.NoError(t, err, "failed to deploy contracts")

	// Create operator info using pre-existing encrypted key files
	operators := make([]operatorInfo, numOperators)
	for i := 0; i < numOperators; i++ {
		operatorKey := fmt.Sprintf("opr%d", i)
		operators[i] = operatorInfo{
			ECDSAPrivateKey: privateKeys.EcdsaMap[operatorKey].PrivateKey,
			ECDSAKeyFile:    privateKeys.EcdsaMap[operatorKey].KeyFile,
			ECDSAPassword:   privateKeys.EcdsaMap[operatorKey].Password,
			BLSKeyPath:      privateKeys.BlsMap[operatorKey].KeyFile,
			BLSPassword:     privateKeys.BlsMap[operatorKey].Password,
		}
	}

	t.Cleanup(func() {
		logger.Info("Stopping containers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = anvilContainer.Terminate(ctx)
		_ = localstackContainer.Terminate(ctx)
	})

	return &TestSetup{
		AnvilContainer:      anvilContainer,
		LocalstackContainer: localstackContainer,
		Contracts:           deploymentResult,
		Operators:           operators,
		PrivateKeys:         privateKeys,
	}
}

func TestChurner(t *testing.T) {
	ctx := t.Context()
	testSetup := setupTest(t)

	// Create mock indexer
	mockIndexer := &indexermock.MockIndexedChainState{}

	// Start churner server directly using the churner package
	grpcPort := "32002"
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	require.NoError(t, err, "failed to listen on port")
	defer func() {
		if err := listener.Close(); err != nil {
			t.Logf("failed to close listener: %v", err)
		}
	}()

	// Create churner config directly (no CLI parsing needed)
	churnerConfig := &churner.Config{
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{"http://localhost:8545"},
			PrivateKeyString: churnerPrivateKeyHex,
		},
		LoggerConfig: common.LoggerConfig{
			Format: common.TextLogFormat,
			HandlerOpts: logging.SLoggerOptions{
				Level:   slog.LevelDebug,
				NoColor: true,
			},
		},
		MetricsConfig: churner.MetricsConfig{
			HTTPPort:      "9095",
			EnableMetrics: true,
		},
		OperatorStateRetrieverAddr: testSetup.Contracts.EigenDA.OperatorStateRetriever,
		EigenDAServiceManagerAddr:  testSetup.Contracts.EigenDA.ServiceManager,
		EigenDADirectory:           testSetup.Contracts.EigenDA.EigenDADirectory,
		ChurnApprovalInterval:      15 * time.Minute,
		PerPublicKeyRateLimit:      1 * time.Second,
	}

	// Create geth client
	gethClient, err := geth.NewMultiHomingClient(churnerConfig.EthClientConfig, gethcommon.Address{}, logger)
	require.NoError(t, err, "failed to create geth client")

	// Create writer
	churnerTx, err := coreeth.NewWriter(
		logger,
		gethClient,
		churnerConfig.OperatorStateRetrieverAddr,
		churnerConfig.EigenDAServiceManagerAddr)
	require.NoError(t, err, "failed to create writer")

	// Create churner with mock indexer
	churnerMetrics := churner.NewMetrics(churnerConfig.MetricsConfig.HTTPPort, logger)
	cn, err := churner.NewChurner(churnerConfig, mockIndexer, churnerTx, logger, churnerMetrics)
	require.NoError(t, err, "failed to create churner")

	// Create churner server
	churnerServer := churner.NewServer(churnerConfig, cn, logger, churnerMetrics)
	err = churnerServer.Start(churnerConfig.MetricsConfig)
	require.NoError(t, err, "failed to start churner server metrics")

	// Create and start gRPC server
	gs := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 300))
	pb.RegisterChurnerServer(gs, churnerServer)
	healthcheck.RegisterHealthServer(pb.Churner_ServiceDesc.ServiceName, gs)

	// Start serving in goroutine
	go func() {
		if err := gs.Serve(listener); err != nil {
			t.Logf("gRPC server stopped: %v", err)
		}
	}()
	defer gs.Stop()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Create gRPC client to connect to the churner
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial churner")
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	}()

	churnerClient := pb.NewChurnerClient(conn)

	quorumIDsUint8 := make([]uint8, len(quorumIds))
	for i, id := range quorumIds {
		quorumIDsUint8[i] = uint8(id)
	}
	var lowestStakeOperatorAddr gethcommon.Address
	var lowestStakeOperatorPubKey *core.G1Point
	var tx *coreeth.Writer
	var operatorPrivateKey *ecdsa.PrivateKey
	var signer blssigner.Signer
	var g1PointBytes []byte
	var g2PointBytes []byte

	for i, op := range testSetup.Operators {
		socket := fmt.Sprintf("localhost:%d:%d", 32000+i, 32100+i) // Simple port assignment

		// Create BLS signer from key file
		opSigner, err := blssigner.NewSigner(blssignerTypes.SignerConfig{
			Path:       op.BLSKeyPath,
			Password:   op.BLSPassword,
			SignerType: blssignerTypes.Local,
		})
		require.NoError(t, err)

		opG1PointHex := opSigner.GetPublicKeyG1()
		opG1PointBytes, err := hex.DecodeString(opG1PointHex)
		require.NoError(t, err)
		opG1Point := new(core.G1Point)
		opG1Point, err = opG1Point.Deserialize(opG1PointBytes)
		require.NoError(t, err)
		opG2PointHex := opSigner.GetPublicKeyG2()
		opG2PointBytes, err := hex.DecodeString(opG2PointHex)
		require.NoError(t, err)
		opG2Point := new(core.G2Point)
		opG2Point, err = opG2Point.Deserialize(opG2PointBytes)
		require.NoError(t, err)
		sk, privateKey, err := plugin.GetECDSAPrivateKey(op.ECDSAKeyFile, op.ECDSAPassword)
		require.NoError(t, err)

		if i == 0 {
			// This is the lowest stake operator that will be eventually churned
			lowestStakeOperatorAddr = sk.Address
			lowestStakeOperatorPubKey = opG1Point
		}
		salt := [32]byte{}
		copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String())))
		expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())
		// Use the hex private key from plugin.GetECDSAPrivateKey for the transactor
		tx = mustCreateTransactorFromScratch(
			t, *privateKey,
			testSetup.Contracts.EigenDA.OperatorStateRetriever,
			testSetup.Contracts.EigenDA.ServiceManager,
			logger)
		if i >= 3 { // MaxOperatorCount is 3, so the 4th operator (index 3) will churn
			// This operator will churn others
			operatorAddr = sk.Address.Hex()
			signer = opSigner
			operatorPrivateKey = sk.PrivateKey
			g1PointBytes = opG1Point.Serialize()
			g2PointBytes = opG2Point.Serialize()
			break
		}
		err = tx.RegisterOperator(ctx, opSigner, socket, quorumIDsUint8, sk.PrivateKey, salt, expiry)
		require.NoError(t, err)
	}
	require.Greater(t, len(lowestStakeOperatorAddr), 0)

	salt := crypto.Keccak256([]byte(operatorToChurnInPrivateKeyHex), []byte("ChurnRequest"))
	request := &pb.ChurnRequest{
		OperatorAddress:            operatorAddr,
		OperatorToRegisterPubkeyG1: g1PointBytes,
		OperatorToRegisterPubkeyG2: g2PointBytes,
		Salt:                       salt,
		QuorumIds:                  quorumIds,
	}

	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		[]byte(request.GetOperatorAddress()),
		request.GetOperatorToRegisterPubkeyG1(),
		request.GetOperatorToRegisterPubkeyG2(),
		request.GetSalt(),
	)
	copy(requestHash[:], requestHashBytes)

	signature, err := signer.Sign(ctx, requestHash[:])
	require.NoError(t, err)
	request.OperatorRequestSignature = signature

	// Set up mock expectation for the lowest stake operator
	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: lowestStakeOperatorPubKey,
	}, nil)

	// Call churner via gRPC instead of direct server call
	reply, err := churnerClient.Churn(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, reply)
	require.NotNil(t, reply.GetSignatureWithSaltAndExpiry().GetSalt())
	require.NotNil(t, reply.GetSignatureWithSaltAndExpiry().GetExpiry())
	require.NotNil(t, reply.GetSignatureWithSaltAndExpiry().GetSignature())
	require.Equal(t, 65, len(reply.GetSignatureWithSaltAndExpiry().GetSignature()))
	require.Len(t, reply.GetOperatorsToChurn(), 2)
	actualQuorums := make([]uint32, 0)
	for _, param := range reply.GetOperatorsToChurn() {
		actualQuorums = append(actualQuorums, param.GetQuorumId())
		require.Equal(t, lowestStakeOperatorAddr, gethcommon.BytesToAddress(param.GetOperator()))
		require.Equal(t, lowestStakeOperatorPubKey.Serialize(), param.GetPubkey())
	}
	require.ElementsMatch(t, quorumIds, actualQuorums)

	salt32 := [32]byte{}
	copy(salt32[:], salt)
	expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())
	err = tx.RegisterOperatorWithChurn(ctx, signer, "localhost:8080", quorumIDsUint8, operatorPrivateKey, salt32, expiry, reply)
	require.NoError(t, err)
}

func mustCreateTransactorFromScratch(
	t *testing.T,
	privateKey string,
	operatorStateRetriever string,
	serviceManager string,
	logger logging.Logger,
) *coreeth.Writer {
	t.Helper()

	ethClientCfg := geth.EthClientConfig{
		RPCURLs:          []string{rpcURL},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       numRetries,
	}

	gethClient, err := geth.NewMultiHomingClient(ethClientCfg, gethcommon.Address{}, logger)
	require.NoError(t, err, "failed to create eth client")

	writer, err := coreeth.NewWriter(logger, gethClient, operatorStateRetriever, serviceManager)
	require.NoError(t, err, "failed to create eth writer")

	return writer
}
