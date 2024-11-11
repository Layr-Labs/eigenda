package encoder_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
)

var (
	logger             = logging.NewNoopLogger()
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	deployLocalStack   bool
	localStackPort     = "4571"
	blobStore          *blobstore.BlobStore
	chunkStoreWriter   chunkstore.ChunkWriter
	chunkStoreReader   chunkstore.ChunkReader
	UUID               = uuid.New()
	s3BucketName       = "test-eigenda"
	mockCommitment     = encoding.BlobCommitments{}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
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
	cfg := aws.DefaultClientConfig()
	cfg.AccessKey = "localstack"
	cfg.SecretAccessKey = "localstack"
	cfg.EndpointURL = fmt.Sprintf("http://0.0.0.0:%s", localStackPort)
	cfg.Region = "us-east-1"

	// Initialize S3 client
	s3Client, err := s3.NewClient(context.Background(), *cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create s3 client: " + err.Error())
	}

	// Create S3 buckets
	err = s3Client.CreateBucket(context.Background(), s3BucketName)
	if err != nil {
		teardown()
		panic("failed to create s3 bucket: " + err.Error())
	}

	// Initialize blob store
	blobStore = blobstore.NewBlobStore(s3BucketName, s3Client, logger)

	// Initialize chunk store writer
	chunkStoreWriter = chunkstore.NewChunkWriter(logger, s3Client, s3BucketName, 512*1024)

	// Initialize chunk store reader
	chunkStoreReader = chunkstore.NewChunkReader(logger, s3Client, s3BucketName)

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
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}
