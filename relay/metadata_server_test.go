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
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	p "github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	prover             *p.Prover
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
	setup(t)

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

func randomBlobHeader(t *testing.T) *v2.BlobHeader {

	data := tu.RandomBytes(128)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitments(data)
	assert.NoError(t, err)
	assert.NoError(t, err)
	commitmentProto, err := commitments.ToProfobuf()

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

	return blobHeader
}

func TestFetchingIndividualMetadata(t *testing.T) {
	tu.InitializeRandom()

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardown()
	}()

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header := randomBlobHeader(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		totalChunkSizeBytes := uint32(rand.Intn(1024 * 1024 * 1024))
		fragmentSizeBytes := uint32(rand.Intn(1024 * 1024))

		totalChunkSizeMap[blobKey] = totalChunkSizeBytes
		fragmentSizeMap[blobKey] = fragmentSizeBytes

		err = metadataStore.PutBlobCertificate(
			context.Background(),
			&v2.BlobCertificate{
				BlobHeader: header,
			},
			&encoding.FragmentInfo{
				TotalChunkSizeBytes: totalChunkSizeBytes,
				FragmentSizeBytes:   fragmentSizeBytes,
			})
		require.NoError(t, err)
	}

	// Read the metadata back
	for blobKey, totalChunkSizeBytes := range totalChunkSizeMap {
		cert, fragmentInfo, err := metadataStore.GetBlobCertificate(context.Background(), blobKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.NotNil(t, fragmentInfo)

		require.Equal(t, totalChunkSizeBytes, fragmentInfo.TotalChunkSizeBytes)
		require.Equal(t, fragmentSizeMap[blobKey], fragmentInfo.FragmentSizeBytes)
	}
}
