package controller_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var (
	logger = testutils.GetLogger()

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource

	deployLocalStack bool
	localStackPort   = "4571"

	s3Client          s3.Client
	dynamoClient      dynamodb.Client
	blobStore         *blobstore.BlobStore
	blobMetadataStore *blobstore.BlobMetadataStore

	UUID              = uuid.New()
	s3BucketName      = "test-eigenda-blobstore"
	metadataTableName = fmt.Sprintf("test-BlobMetadata-%v", UUID)

	mockCommitment = encoding.BlobCommitments{}

	heartbeatChan      = make(chan time.Time, 10) // Stores last 10 heartbeats
	heartbeatsReceived []time.Time
	mu                 sync.Mutex
	doneListening      = make(chan struct{})
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {

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

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	_, err := test_utils.CreateTable(context.Background(), cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}

	blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName)

	s3Client, err = s3.NewClient(context.Background(), cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create s3 client: " + err.Error())
	}
	err = s3Client.CreateBucket(context.Background(), s3BucketName)
	if err != nil {
		teardown()
		panic("failed to create s3 bucket: " + err.Error())
	}
	blobStore = blobstore.NewBlobStore(s3BucketName, s3Client, logger)

	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err = lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	if err != nil {
		teardown()
		panic("failed to create mock commitment: " + err.Error())
	}
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	if err != nil {
		teardown()
		panic("failed to create mock commitment: " + err.Error())
	}
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	if err != nil {
		teardown()
		panic("failed to create mock commitment: " + err.Error())
	}
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	if err != nil {
		teardown()
		panic("failed to create mock commitment: " + err.Error())
	}

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	mockCommitment = encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           16,
	}
}

func teardown() {
	mu.Lock()
	defer mu.Unlock()

	if len(heartbeatsReceived) == 0 {
		logger.Error("Expected heartbeats, but none were received")
	}

	close(heartbeatChan) // Ensure the goroutine exits properly

	select {
	case <-doneListening:
	default:
		close(doneListening)
	}

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func newBlob(t *testing.T, quorumNumbers []core.QuorumID) (corev2.BlobKey, *corev2.BlobHeader) {
	accountBytes := make([]byte, 32)
	_, err := rand.Read(accountBytes)
	require.NoError(t, err)
	accountID := gethcommon.HexToAddress(hex.EncodeToString(accountBytes))
	timestamp, err := rand.Int(rand.Reader, big.NewInt(256))
	require.NoError(t, err)
	cumulativePayment, err := rand.Int(rand.Reader, big.NewInt(1024))
	require.NoError(t, err)
	sig := make([]byte, 32)
	_, err = rand.Read(sig)
	require.NoError(t, err)
	bh := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   quorumNumbers,
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			Timestamp:         timestamp.Int64(),
			CumulativePayment: cumulativePayment,
		},
	}
	bk, err := bh.BlobKey()
	require.NoError(t, err)
	return bk, bh
}
