package relay

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
)

func setupLocalstack() error {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
		if err != nil && err.Error() == "container already exists" {
			teardownLocalstack()
			return err
		}
	}
	return nil
}

func teardownLocalstack() {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func buildMetadataStore(t *testing.T) *blobstore.BlobMetadataStore {
	setupLocalstack()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	config := aws.DefaultClientConfig()
	config.EndpointURL = localstackHost
	config.Region = "us-east-1"

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

	dynamoClient, err := dynamodb.NewClient(cfg, logger)
	require.NoError(t, err)

	return blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		metadataTableName)
}

func randomBlobHeader() *v2.BlobHeader {

	blobCommitments := encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: fp.Element{rand.Uint64(), rand.Uint64(), rand.Uint64(), rand.Uint64()},
			Y: fp.Element{rand.Uint64(), rand.Uint64(), rand.Uint64(), rand.Uint64()},
		},
	}

	return &v2.BlobHeader{
		BlobCommitments: blobCommitments,
	}
}

func TestFetchingIndividualMetadata(t *testing.T) {
	tu.InitializeRandom()

	metadataStore := buildMetadataStore(t)
	defer func() {
		teardownLocalstack()
	}()

	totalChunkSizeMap := make(map[v2.BlobKey]uint32)
	fragmentSizeMap := make(map[v2.BlobKey]uint32)

	// Write some metadata
	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header := randomBlobHeader()
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

}
