package relay

import (
	"context"
	"fmt"
	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbcommonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	p "github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	prover             *p.Prover
	bucketName         = fmt.Sprintf("test-bucket-%v", UUID)
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
)

func setup(t *testing.T) {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "..")
	changeDirectory(filepath.Join(rootPath, "inabox"))

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
		require.NoError(t, err)
	}

	// Only set up the prover once, it's expensive
	if prover == nil {
		config := &kzg.KzgConfig{
			G1Path:          "./resources/kzg/g1.point.300000",
			G2Path:          "./resources/kzg/g2.point.300000",
			CacheDir:        "./resources/kzg/SRSTables",
			SRSOrder:        8192,
			SRSNumberToLoad: 8192,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		}
		var err error
		prover, err = p.NewProver(config, true)
		require.NoError(t, err)
	}
}

func changeDirectory(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("Failed to change directories. Error: %s", err)
	}

	newDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get working directory. Error: %s", err)
	}
	log.Printf("Current Working Directory: %s\n", newDir)
}

func teardown() {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func buildMetadataStore(t *testing.T) *blobstore.BlobMetadataStore {

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
	require.NoError(t, err)
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
	require.NoError(t, err)

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     localstackHost,
	}

	_, err = test_utils.CreateTable(
		context.Background(),
		cfg,
		metadataTableName,
		blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	require.NoError(t, err)

	dynamoClient, err := dynamodb.NewClient(cfg, logger)
	require.NoError(t, err)

	return blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		metadataTableName)
}

func buildBlobStore(t *testing.T, logger logging.Logger) *blobstore.BlobStore {
	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     localstackHost,
	}

	client, err := s3.NewClient(context.Background(), cfg, logger)
	require.NoError(t, err)

	err = client.CreateBucket(context.Background(), bucketName)
	require.NoError(t, err)

	return blobstore.NewBlobStore(bucketName, client, logger)
}

func buildChunkStore(
	t *testing.T,
	logger logging.Logger,
	shards []uint32) (chunkstore.ChunkReader, chunkstore.ChunkWriter) {

	cfg := aws.ClientConfig{
		Region:               "us-east-1",
		AccessKey:            "localstack",
		SecretAccessKey:      "localstack",
		EndpointURL:          localstackHost,
		FragmentWriteTimeout: time.Duration(10) * time.Second,
		FragmentReadTimeout:  time.Duration(10) * time.Second,
	}

	client, err := s3.NewClient(context.Background(), cfg, logger)
	require.NoError(t, err)

	err = client.CreateBucket(context.Background(), bucketName)
	require.NoError(t, err)

	// intentionally use very small fragment size
	chunkWriter := chunkstore.NewChunkWriter(logger, client, bucketName, 32)
	chunkReader := chunkstore.NewChunkReader(logger, client, bucketName, shards)

	return chunkReader, chunkWriter
}

func randomBlob(t *testing.T) (*v2.BlobHeader, []byte) {

	data := tu.RandomBytes(128)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitments(data)
	require.NoError(t, err)
	require.NoError(t, err)
	commitmentProto, err := commitments.ToProfobuf()
	require.NoError(t, err)

	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         tu.RandomString(10),
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := v2.NewBlobHeader(blobHeaderProto)
	require.NoError(t, err)

	return blobHeader, data
}

func randomBlobChunks(t *testing.T) (*v2.BlobHeader, []byte, []*encoding.Frame) {
	header, data := randomBlob(t)

	params := encoding.ParamsFromMins(16, 4)
	_, frames, err := prover.EncodeAndProve(data, params)
	require.NoError(t, err)

	return header, data, frames
}

func disassembleFrames(frames []*encoding.Frame) ([]*rs.Frame, []*encoding.Proof) {
	rsFrames := make([]*rs.Frame, len(frames))
	proofs := make([]*encoding.Proof, len(frames))

	for i, frame := range frames {
		rsFrames[i] = &rs.Frame{
			Coeffs: frame.Coeffs,
		}
		proofs[i] = &frame.Proof
	}

	return rsFrames, proofs
}