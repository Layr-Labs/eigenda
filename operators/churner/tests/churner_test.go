package test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	indexermock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
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
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
}

var (
	templateName string
	testName     string

	localstackPort                 = "4570"
	mockIndexer                    = &indexermock.MockIndexedChainState{}
	rpcURL                         = "http://localhost:8545"
	quorumIds                      = []uint32{0, 1}
	operatorAddr                   = ""
	churnerPrivateKeyHex           = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	operatorToChurnInPrivateKeyHex = "0000000000000000000000000000000000000000000000000000000000000020"
	numRetries                     = 0

	logger = test.GetLogger()
)

func setupTest(t *testing.T) (*testbed.AnvilContainer, *testbed.LocalStackContainer, *deploy.Config) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping churner test in short mode")
	}

	flag.Parse()
	ctx := t.Context()
	rootPath := "../../../"

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		require.NoError(t, err, "failed to create test directory")
	}

	testConfig := deploy.NewTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = false

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start localstack container")

	anvilContainer, err := testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true, // This will bind container port 8545 to host port 8545
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start anvil container")

	logger.Info("Deploying experiment")
	err = testConfig.DeployExperiment()
	require.NoError(t, err, "failed to deploy experiment")

	t.Cleanup(func() {
		logger.Info("Stopping containers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = anvilContainer.Terminate(ctx)
		_ = localstackContainer.Terminate(ctx)
	})

	return anvilContainer, localstackContainer, testConfig
}

func TestChurner(t *testing.T) {
	ctx := t.Context()
	_, _, testConfig := setupTest(t)

	server := newTestServer(t, testConfig)
	quorumIDsUint8 := make([]uint8, len(quorumIds))
	for i, id := range quorumIds {
		quorumIDsUint8[i] = uint8(id)
	}
	var lowestStakeOperatorAddr gethcommon.Address
	var lowestStakeOperatorPubKey *core.G1Point
	var tx *eth.Writer
	var operatorPrivateKey *ecdsa.PrivateKey
	var signer blssigner.Signer
	var g1PointBytes []byte
	var g2PointBytes []byte
	for i, op := range testConfig.Operators {
		socket := fmt.Sprintf("%s:%s:%s", op.NODE_HOSTNAME, op.NODE_DISPERSAL_PORT, op.NODE_RETRIEVAL_PORT)
		opSigner, err := blssigner.NewSigner(blssignerTypes.SignerConfig{
			Path:       op.NODE_BLS_KEY_FILE,
			Password:   op.NODE_BLS_KEY_PASSWORD,
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
		sk, privateKey, err := plugin.GetECDSAPrivateKey(op.NODE_ECDSA_KEY_FILE, op.NODE_ECDSA_KEY_PASSWORD)
		require.NoError(t, err)
		if i == 0 {
			// This is the lowest stake operator that will be eventually churned
			lowestStakeOperatorAddr = sk.Address
			lowestStakeOperatorPubKey = opG1Point
		}
		salt := [32]byte{}
		copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String())))
		expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())
		tx = mustCreateTransactorFromScratch(
			t, *privateKey, testConfig.EigenDA.OperatorStateRetriever, testConfig.EigenDA.ServiceManager, logger)
		if i >= testConfig.Services.Counts.NumMaxOperatorCount {
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

	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: lowestStakeOperatorPubKey,
	}, nil)

	reply, err := server.Churn(ctx, request)
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
) *eth.Writer {
	t.Helper()

	ethClientCfg := geth.EthClientConfig{
		RPCURLs:          []string{rpcURL},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       numRetries,
	}

	gethClient, err := geth.NewMultiHomingClient(ethClientCfg, gethcommon.Address{}, logger)
	require.NoError(t, err, "failed to create eth client")

	writer, err := eth.NewWriter(logger, gethClient, operatorStateRetriever, serviceManager)
	require.NoError(t, err, "failed to create eth writer")

	return writer
}

func newTestServer(t *testing.T, testConfig *deploy.Config) *churner.Server {
	t.Helper()

	config := &churner.Config{
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{rpcURL},
			PrivateKeyString: churnerPrivateKeyHex,
			NumRetries:       numRetries,
		},
		LoggerConfig:               *common.DefaultLoggerConfig(),
		OperatorStateRetrieverAddr: testConfig.EigenDA.OperatorStateRetriever,
		EigenDAServiceManagerAddr:  testConfig.EigenDA.ServiceManager,
		EigenDADirectory:           testConfig.EigenDA.EigenDADirectory,
		ChurnApprovalInterval:      15 * time.Minute,
	}

	operatorTransactorChurner := mustCreateTransactorFromScratch(
		t,
		churnerPrivateKeyHex,
		testConfig.EigenDA.OperatorStateRetriever,
		testConfig.EigenDA.ServiceManager,
		logger,
	)

	metrics := churner.NewMetrics("9001", logger)
	cn, err := churner.NewChurner(config, mockIndexer, operatorTransactorChurner, logger, metrics)
	require.NoError(t, err)

	return churner.NewServer(config, cn, logger, metrics)
}
