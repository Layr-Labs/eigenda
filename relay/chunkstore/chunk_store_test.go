package chunkstore

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	s3common "github.com/Layr-Labs/eigenda/common/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/stretchr/testify/require"
)

var (
	logger = test.GetLogger()
)

const (
	localstackPort = "4577"
	localstackHost = "http://0.0.0.0:4577"
	bucket         = "eigen-test"
)

func setupLocalStackTest(t *testing.T) s3.Client {
	t.Helper()

	ctx := t.Context()

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb"},
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start LocalStack container")

	t.Cleanup(func() {
		logger.Info("Stopping LocalStack container")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = localstackContainer.Terminate(ctx)
	})

	config := aws.DefaultClientConfig()
	config.EndpointURL = localstackHost
	config.Region = "us-east-1"

	err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
	require.NoError(t, err, "failed to set AWS_ACCESS_KEY_ID")
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
	require.NoError(t, err, "failed to set AWS_SECRET_ACCESS_KEY")

	client, err := s3.NewClient(ctx, *config, logger)
	require.NoError(t, err, "failed to create S3 client")

	err = client.CreateBucket(ctx, bucket)
	require.NoError(t, err, "failed to create S3 bucket")

	return client
}

func getProofs(t *testing.T, count int) []*encoding.Proof {
	t.Helper()

	proofs := make([]*encoding.Proof, count)

	// Note from Cody: I'd rather use randomized proofs here, but I'm not sure how to generate them.
	// Using random data breaks since the deserialization logic rejects invalid proofs.
	var x, y fp.Element
	_, err := x.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	require.NoError(t, err, "failed to set X element for proof")
	_, err = y.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	require.NoError(t, err, "failed to set Y element for proof")

	for i := 0; i < count; i++ {
		proof := encoding.Proof{
			X: x,
			Y: y,
		}
		proofs[i] = &proof

	}

	return proofs
}

func runRandomProofsTest(t *testing.T, client s3.Client) {
	t.Helper()
	ctx := t.Context()

	fragmentSize := rand.Intn(1024) + 100 // ignored since we aren't writing coefficients

	writer := NewChunkWriter(logger, client, bucket, fragmentSize)
	reader := NewChunkReader(logger, client, bucket)

	expectedValues := make(map[corev2.BlobKey][]*encoding.Proof)

	// Write data
	for i := 0; i < 100; i++ {
		key := corev2.BlobKey(random.RandomBytes(32))

		proofs := getProofs(t, rand.Intn(100)+100)
		expectedValues[key] = proofs

		err := writer.PutFrameProofs(ctx, key, proofs)
		require.NoError(t, err, "failed to put frame proofs for blob key %x", key)
	}

	// Read data
	for key, expectedProofs := range expectedValues {
		binaryProofs, err := reader.GetBinaryChunkProofs(ctx, key)
		require.NoError(t, err, "failed to get binary chunk proofs for blob key %x", key)
		proofs := encoding.DeserializeSplitFrameProofs(binaryProofs)
		require.Equal(t, expectedProofs, proofs, "proof mismatch for blob key %x", key)
	}
}

func TestRandomProofs(t *testing.T) {
	random.InitializeRandom()

	t.Run("mock_client", func(t *testing.T) {
		client := s3common.NewMockS3Client()
		runRandomProofsTest(t, client)
	})

	t.Run("localstack_client", func(t *testing.T) {
		client := setupLocalStackTest(t)
		runRandomProofsTest(t, client)
	})
}

func generateRandomFrameCoeffs(
	t *testing.T,
	encoder *rs.Encoder,
	size int,
	params encoding.EncodingParams) []rs.FrameCoeffs {

	frames, _, err := encoder.EncodeBytes(codec.ConvertByPaddingEmptyByte(random.RandomBytes(size)), params)
	require.NoError(t, err, "failed to encode bytes into frame coefficients")
	return frames
}

func runRandomCoefficientsTest(t *testing.T, client s3.Client) {
	t.Helper()
	ctx := t.Context()

	chunkSize := uint64(rand.Intn(1024) + 100)
	fragmentSize := int(chunkSize / 2)
	params := encoding.ParamsFromSysPar(3, 1, chunkSize)
	cfg := encoding.DefaultConfig()
	encoder := rs.NewEncoder(logger, cfg)

	writer := NewChunkWriter(logger, client, bucket, fragmentSize)
	reader := NewChunkReader(logger, client, bucket)

	expectedValues := make(map[corev2.BlobKey][]rs.FrameCoeffs)
	metadataMap := make(map[corev2.BlobKey]*encoding.FragmentInfo)

	// Write data
	for i := 0; i < 100; i++ {
		key := corev2.BlobKey(random.RandomBytes(32))

		coefficients := generateRandomFrameCoeffs(t, encoder, int(chunkSize), params)
		expectedValues[key] = coefficients

		metadata, err := writer.PutFrameCoefficients(ctx, key, coefficients)
		require.NoError(t, err, "failed to put frame coefficients for blob key %x", key)
		metadataMap[key] = metadata
	}

	// Read data
	for key, expectedCoefficients := range expectedValues {
		elementCount, binaryCoefficients, err :=
			reader.GetBinaryChunkCoefficients(ctx, key, metadataMap[key])
		require.NoError(t, err, "failed to get binary chunk coefficients for blob key %x", key)
		coefficients := rs.DeserializeSplitFrameCoeffs(elementCount, binaryCoefficients)
		require.NoError(t, err, "failed to deserialize frame coefficients for blob key %x", key)
		require.Equal(t, len(expectedCoefficients), len(coefficients), "coefficient count mismatch for blob key %x", key)
		for i := 0; i < len(expectedCoefficients); i++ {
			require.Equal(t, expectedCoefficients[i], coefficients[i],
				"coefficient mismatch at index %d for blob key %x", i, key)
		}
	}
}

func TestRandomCoefficients(t *testing.T) {
	random.InitializeRandom()

	t.Run("mock_client", func(t *testing.T) {
		client := s3common.NewMockS3Client()
		runRandomCoefficientsTest(t, client)
	})

	t.Run("localstack_client", func(t *testing.T) {
		client := setupLocalStackTest(t)
		runRandomCoefficientsTest(t, client)
	})
}

func TestCheckProofCoefficientsExist(t *testing.T) {
	random.InitializeRandom()
	client := s3common.NewMockS3Client()

	chunkSize := uint64(rand.Intn(1024) + 100)
	fragmentSize := int(chunkSize / 2)

	params := encoding.ParamsFromSysPar(3, 1, chunkSize)
	cfg := encoding.DefaultConfig()
	encoder := rs.NewEncoder(logger, cfg)

	writer := NewChunkWriter(logger, client, bucket, fragmentSize)
	ctx := t.Context()
	for i := 0; i < 100; i++ {
		key := corev2.BlobKey(random.RandomBytes(32))

		proofs := getProofs(t, rand.Intn(100)+100)
		err := writer.PutFrameProofs(ctx, key, proofs)
		require.NoError(t, err, "failed to put frame proofs for blob key %x", key)
		require.True(t, writer.ProofExists(ctx, key), "proof should exist for blob key %x", key)

		coefficients := generateRandomFrameCoeffs(t, encoder, int(chunkSize), params)
		metadata, err := writer.PutFrameCoefficients(ctx, key, coefficients)
		require.NoError(t, err, "failed to put frame coefficients for blob key %x", key)
		exist, fragmentInfo := writer.CoefficientsExists(ctx, key)
		require.True(t, exist, "coefficients should exist for blob key %x", key)
		require.Equal(t, metadata, fragmentInfo, "fragment info mismatch for blob key %x", key)
	}
}
