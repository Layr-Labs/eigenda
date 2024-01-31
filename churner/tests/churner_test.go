package test

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	commock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	dacore "github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	indexermock "github.com/Layr-Labs/eigenda/core/thegraph/mock"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
}

var (
	keyPair                        *dacore.KeyPair
	testConfig                     *deploy.Config
	templateName                   string
	testName                       string
	logger                         = &commock.Logger{}
	mockIndexer                    = &indexermock.MockIndexedChainState{}
	rpcURL                         = "http://localhost:8545"
	quorumIds                      = []uint32{0}
	operatorAddr                   = gethcommon.HexToAddress("0x0000000000000000000000000000000000000001")
	churnerPrivateKeyHex           = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	operatorToChurnInPrivateKeyHex = "0000000000000000000000000000000000000000000000000000000000000020"
)

func TestMain(m *testing.M) {
	flag.Parse()
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {
	rootPath := "../../"

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
		fmt.Println("Skipping churner test in short mode")
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
		testConfig.StopGraphNode()
	}
}

func TestChurner(t *testing.T) {
	ctx := context.Background()

	// get the operator to churn in's transactor and bls public key registered
	op := testConfig.Operators[0]
	operatorTransactor, err := createTransactorFromScratch(
		op.NODE_PRIVATE_KEY,
		testConfig.EigenDA.OperatorStateRetreiver,
		testConfig.EigenDA.ServiceManager,
		logger,
	)
	assert.NoError(t, err)

	keyPair, err := dacore.GenRandomBlsKeys()
	assert.NoError(t, err)

	quorumIds_ := make([]uint8, len(quorumIds))
	for i, q := range quorumIds {
		quorumIds_[i] = uint8(q)
	}

	operatorSalt := [32]byte{}
	_, err = rand.Read(operatorSalt[:])
	assert.NoError(t, err)

	expiry := big.NewInt(1000)
	privKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	err = operatorTransactor.RegisterOperator(ctx, keyPair, "socket", quorumIds_, privKey, operatorSalt, expiry)
	assert.NoError(t, err)

	server := newTestServer(t)

	salt := crypto.Keccak256([]byte(operatorToChurnInPrivateKeyHex), []byte("ChurnRequest"))
	request := &pb.ChurnRequest{
		OperatorToRegisterPubkeyG1: keyPair.PubKey.Serialize(),
		OperatorToRegisterPubkeyG2: keyPair.GetPubKeyG2().Serialize(),
		Salt:                       salt,
		QuorumIds:                  quorumIds,
	}

	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		request.OperatorToRegisterPubkeyG1,
		request.OperatorToRegisterPubkeyG2,
		request.Salt,
	)
	copy(requestHash[:], requestHashBytes)

	signature := keyPair.SignMessage(requestHash)
	request.OperatorRequestSignature = signature.Serialize()

	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: keyPair.PubKey,
	}, nil)

	reply, err := server.Churn(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotNil(t, reply.SignatureWithSaltAndExpiry.GetSalt())
	assert.NotNil(t, reply.SignatureWithSaltAndExpiry.GetExpiry())
	assert.NotNil(t, reply.SignatureWithSaltAndExpiry.GetSignature())
	assert.Equal(t, 65, len(reply.SignatureWithSaltAndExpiry.GetSignature()))

	for _, param := range reply.OperatorsToChurn {
		assert.Equal(t, uint32(0), param.GetQuorumId())
		assert.Equal(t, operatorAddr.Bytes(), param.GetOperator())
		assert.Equal(t, keyPair.PubKey.Serialize(), param.GetPubkey())
	}
}

func createTransactorFromScratch(privateKey, operatorStateRetriever, serviceManager string, logger common.Logger) (*eth.Transactor, error) {
	ethClientCfg := geth.EthClientConfig{
		RPCURL:           rpcURL,
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
	}

	gethClient, err := geth.NewClient(ethClientCfg, logger)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	return eth.NewTransactor(logger, gethClient, operatorStateRetriever, serviceManager)
}

func newTestServer(t *testing.T) *churner.Server {

	var err error
	keyPair, err = dacore.GenRandomBlsKeys()
	if err != nil {
		t.Fatalf("Generating random BLS keys Error: %s", err.Error())
	}

	config := &churner.Config{
		EthClientConfig: geth.EthClientConfig{
			RPCURL:           rpcURL,
			PrivateKeyString: churnerPrivateKeyHex,
		},
		LoggerConfig:                  logging.DefaultCLIConfig(),
		BLSOperatorStateRetrieverAddr: testConfig.EigenDA.OperatorStateRetreiver,
		EigenDAServiceManagerAddr:     testConfig.EigenDA.ServiceManager,
	}

	operatorTransactorChurner, err := createTransactorFromScratch(
		churnerPrivateKeyHex,
		testConfig.EigenDA.OperatorStateRetreiver,
		testConfig.EigenDA.ServiceManager,
		logger,
	)
	assert.NoError(t, err)

	metrics := churner.NewMetrics("9001", logger)
	cn, err := churner.NewChurner(config, mockIndexer, operatorTransactorChurner, logger, metrics)
	assert.NoError(t, err)

	return churner.NewServer(config, cn, logger, metrics)
}
