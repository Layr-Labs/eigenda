package encoder_test

import (
	"context"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws/mock"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

type testComponents struct {
	encoderServer    *encoder.EncoderServerV2
	blobStore        *blobstore.BlobStore
	chunkStoreWriter chunkstore.ChunkWriter
	chunkStoreReader chunkstore.ChunkReader
	s3Client         *mock.S3Client
	dynamoDBClient   *mock.MockDynamoDBClient
}

func makeTestProver(numPoint uint64) (encoding.Prover, error) {
	// We need the larger SRS for testing the encoder with 8192 chunks
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point.300000",
		G2Path:          "../../inabox/resources/kzg/g2.point.300000",
		G2PowerOf2Path:  "../../inabox/resources/kzg/g2.point.300000.powerOf2",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        300000,
		SRSNumberToLoad: numPoint,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	p, err := prover.NewProver(kzgConfig, false)

	return p, err
}

func TestEncodeBlob(t *testing.T) {
	const (
		testDataSize   = 16 * 1024
		timeoutSeconds = 30
		randSeed       = uint64(42)
	)

	var (
		codingRatio = corev2.ParametersMap[0].CodingRate
		numChunks   = corev2.ParametersMap[0].NumChunks
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	createTestData := func(t *testing.T, size int) []byte {
		t.Helper()
		data := make([]byte, size)
		_, err := rand.New(rand.NewSource(randSeed)).Read(data)
		if !assert.NoError(t, err, "Failed to create test data") {
			t.FailNow()
		}

		return codec.ConvertByPaddingEmptyByte(data)
	}

	c := createTestComponents(t)
	server := c.encoderServer

	// Setup test data
	data := createTestData(t, testDataSize)
	blobSize := uint(len(data))
	blobLength := encoding.GetBlobLength(blobSize)

	// Get chunk length for blob version 0
	chunkLength, err := corev2.GetChunkLength(0, core.NextPowerOf2(uint32(blobLength)))
	if !assert.NoError(t, err, "Failed to get chunk length") {
		t.FailNow()
	}

	t.Logf("Test parameters: blobversion=%d, blobLength=%d, codingRatio=%d, numChunks=%d, chunkLength=%d",
		0, blobLength, codingRatio, numChunks, chunkLength)

	// Create blob header and key
	blobHeader := createTestBlobHeader(t)
	blobKey, err := blobHeader.BlobKey()
	if !assert.NoError(t, err, "Failed to create blob key") {
		t.FailNow()
	}

	// Store test data
	if err := c.blobStore.StoreBlob(ctx, blobKey, data); !assert.NoError(t, err, "Failed to store blob") {
		t.FailNow()
	}

	// Verify storage succeded
	t.Run("Verify Blob Storage", func(t *testing.T) {
		storedData, err := c.blobStore.GetBlob(ctx, blobKey)
		assert.NoError(t, err, "Failed to get stored blob")
		assert.Equal(t, data, storedData, "Stored data doesn't match original")
	})

	// Create and execute encoding request
	req := &pb.EncodeBlobRequest{
		BlobKey: blobKey[:],
		EncodingParams: &pb.EncodingParams{
			ChunkLength: uint64(chunkLength),
			NumChunks:   uint64(numChunks),
		},
	}

	expectedUploadCalls := 1
	expectedFragmentedUploadObjectCalls := 0
	assert.Equal(t, c.s3Client.Called["UploadObject"], expectedUploadCalls)
	assert.Equal(t, c.s3Client.Called["FragmentedUploadObject"], expectedFragmentedUploadObjectCalls)
	resp, err := server.EncodeBlob(ctx, req)
	if !assert.NoError(t, err, "EncodeBlob failed") {
		t.FailNow()
	}
	expectedUploadCalls++
	expectedFragmentedUploadObjectCalls++
	assert.Equal(t, c.s3Client.Called["UploadObject"], expectedUploadCalls)
	assert.Equal(t, c.s3Client.Called["FragmentedUploadObject"], expectedFragmentedUploadObjectCalls)

	// Verify encoding results
	t.Run("Verify Encoding Results", func(t *testing.T) {
		assert.NotNil(t, resp, "Response should not be nil")
		assert.Equal(t, uint32(294916), resp.FragmentInfo.TotalChunkSizeBytes, "Unexpected total chunk size")
		assert.Equal(t, uint32(512*1024), resp.FragmentInfo.FragmentSizeBytes, "Unexpected fragment size")
	})

	expectedFragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: resp.FragmentInfo.TotalChunkSizeBytes,
		FragmentSizeBytes:   resp.FragmentInfo.FragmentSizeBytes,
	}

	// Verify chunk store data
	t.Run("Verify Chunk Store Data", func(t *testing.T) {
		// Check proofs
		assert.True(t, c.chunkStoreWriter.ProofExists(ctx, blobKey))
		proofs, err := c.chunkStoreReader.GetChunkProofs(ctx, blobKey)
		assert.NoError(t, err, "Failed to get chunk proofs")
		assert.Len(t, proofs, int(numChunks), "Unexpected number of proofs")

		// Check coefficients
		coefExist, fetchedFragmentInfo := c.chunkStoreWriter.CoefficientsExists(ctx, blobKey)
		assert.True(t, coefExist, "Coefficients should exist")
		assert.Equal(t, expectedFragmentInfo, fetchedFragmentInfo, "Unexpected fragment info")

		coefficients, err := c.chunkStoreReader.GetChunkCoefficients(ctx, blobKey, expectedFragmentInfo)
		assert.NoError(t, err, "Failed to get chunk coefficients")
		assert.Len(t, coefficients, int(numChunks), "Unexpected number of coefficients")
	})

	t.Run("Verify Re-encoding is prevented", func(t *testing.T) {
		assert.Equal(t, c.s3Client.Called["UploadObject"], expectedUploadCalls)
		assert.Equal(t, c.s3Client.Called["FragmentedUploadObject"], expectedFragmentedUploadObjectCalls)
		// Create and execute encoding request again
		resp, err := server.EncodeBlob(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, uint32(294916), resp.FragmentInfo.TotalChunkSizeBytes, "Unexpected total chunk size")
		assert.Equal(t, uint32(512*1024), resp.FragmentInfo.FragmentSizeBytes, "Unexpected fragment size")
		assert.Equal(t, c.s3Client.Called["UploadObject"], expectedUploadCalls)
		assert.Equal(t, c.s3Client.Called["FragmentedUploadObject"], expectedFragmentedUploadObjectCalls)
	})
}

// Helper function to create test blob header
func createTestBlobHeader(t *testing.T) *corev2.BlobHeader {
	t.Helper()
	return &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x1234",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
}

// Helper function to initialize encoder
func createTestComponents(t *testing.T) *testComponents {
	t.Helper()
	prover, err := makeTestProver(300000)
	require.NoError(t, err, "Failed to create prover")
	metrics := encoder.NewMetrics("9000", logger)
	s3Client := mock.NewS3Client()
	dynamoDBClient := &mock.MockDynamoDBClient{}
	blobStore := blobstore.NewBlobStore(s3BucketName, s3Client, logger)
	chunkStoreWriter := chunkstore.NewChunkWriter(logger, s3Client, s3BucketName, 512*1024)
	chunkStoreReader := chunkstore.NewChunkReader(logger, s3Client, s3BucketName)
	encoderServer := encoder.NewEncoderServerV2(encoder.ServerConfig{
		GrpcPort:              "8080",
		MaxConcurrentRequests: 10,
		RequestPoolSize:       5,
		PreventReencoding:     true,
	}, blobStore, chunkStoreWriter, logger, prover, metrics)

	return &testComponents{
		encoderServer:    encoderServer,
		blobStore:        blobStore,
		chunkStoreWriter: chunkStoreWriter,
		chunkStoreReader: chunkStoreReader,
		s3Client:         s3Client,
		dynamoDBClient:   dynamoDBClient,
	}
}
