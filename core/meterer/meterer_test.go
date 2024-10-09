package meterer_test

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	dynamoClient       *commondynamodb.Client
	clientConfig       commonaws.ClientConfig
	privateKey1        *ecdsa.PrivateKey
	privateKey2        *ecdsa.PrivateKey
	signer             *auth.EIP712Signer
	mt                 *meterer.Meterer

	deployLocalStack  bool
	localStackPort    = "4566"
	paymentChainState = &mock.MockOnchainPaymentState{}
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

// // Mock data initialization method
// func InitializeMockPayments(pcs *meterer.OnchainPaymentState, privateKey1 *ecdsa.PrivateKey, privateKey2 *ecdsa.PrivateKey) {
// 	// Initialize mock active reservations
// 	now := uint64(time.Now().Unix())
// 	pcs.ActiveReservations = map[string]core.ActiveReservation{
// 		crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(): {DataRate: 100, StartTimestamp: now + 1200, EndTimestamp: now + 1800, QuorumSplit: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}},
// 		crypto.PubkeyToAddress(privateKey2.PublicKey).Hex(): {DataRate: 200, StartTimestamp: now - 120, EndTimestamp: now + 180, QuorumSplit: []byte{30, 70}, QuorumNumbers: []uint8{0, 1}},
// 	}
// 	pcs.OnDemandPayments = map[string]core.OnDemandPayment{
// 		crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(): {CumulativePayment: 1500},
// 		crypto.PubkeyToAddress(privateKey2.PublicKey).Hex(): {CumulativePayment: 1000},
// 	}
// }

func setup(_ *testing.M) {

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		teardown()
		panic("failed to create logger")
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	dynamoClient, err = commondynamodb.NewClient(clientConfig, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client")
	}

	chainID := big.NewInt(17000)
	verifyingContract := gethcommon.HexToAddress("0x1234000000000000000000000000000000000000")
	signer = auth.NewEIP712Signer(chainID, verifyingContract)

	privateKey1, err = crypto.GenerateKey()
	privateKey2, err = crypto.GenerateKey()

	logger = logging.NewNoopLogger()
	config := meterer.Config{
		PricePerChargeable:   1,
		MinChargeableSize:    1,
		GlobalBytesPerSecond: 1000,
		ReservationWindow:    60,
		ChainReadTimeout:     3 * time.Second,
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:4566"),
	}

	store, err := meterer.NewOffchainStore(
		clientConfig,
		"reservations",
		"ondemand",
		"global",
		logger,
	)
	if err != nil {
		teardown()
		panic("failed to create offchain store")
	}

	// add some default sensible configs
	mt, err = meterer.NewMeterer(
		config,
		paymentChainState,
		store,
		logging.NewNoopLogger(),
		// metrics.NewNoopMetrics(),
	)

	if err != nil {
		teardown()
		panic("failed to create meterer")
	}
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}