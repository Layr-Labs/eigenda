package chunkstore

import (
	"context"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/mock"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_cpu "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
	bucket         = "eigen-test"
)

type clientBuilder struct {
	// This method is called at the beginning of the test.
	start func() error
	// This method is called to build a new client.
	build func() (s3.Client, error)
	// This method is called at the end of the test when all operations are done.
	finish func() error
}

var clientBuilders = []*clientBuilder{
	{
		start: func() error {
			return nil
		},
		build: func() (s3.Client, error) {
			return mock.NewS3Client(), nil
		},
		finish: func() error {
			return nil
		},
	},
	{
		start: func() error {
			return setupLocalstack()
		},
		build: func() (s3.Client, error) {

			logger, err := common.NewLogger(common.DefaultLoggerConfig())
			if err != nil {
				return nil, err
			}

			config := aws.DefaultClientConfig()
			config.EndpointURL = localstackHost
			config.Region = "us-east-1"

			err = os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
			if err != nil {
				return nil, err
			}
			err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
			if err != nil {
				return nil, err
			}

			client, err := s3.NewClient(context.Background(), *config, logger)
			if err != nil {
				return nil, err
			}

			err = client.CreateBucket(context.Background(), bucket)
			if err != nil {
				return nil, err
			}

			return client, nil
		},
		finish: func() error {
			teardownLocalstack()
			return nil
		},
	},
}

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

func getProofs(t *testing.T, count int) []*encoding.Proof {
	proofs := make([]*encoding.Proof, count)

	// Note from Cody: I'd rather use randomized proofs here, but I'm not sure how to generate them.
	// Using random data breaks since the deserialization logic rejects invalid proofs.
	var x, y fp.Element
	_, err := x.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	require.NoError(t, err)
	_, err = y.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	require.NoError(t, err)

	for i := 0; i < count; i++ {
		proof := encoding.Proof{
			X: x,
			Y: y,
		}
		proofs[i] = &proof

	}

	return proofs
}

func RandomProofsTest(t *testing.T, client s3.Client) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	fragmentSize := rand.Intn(1024) + 100 // ignored since we aren't writing coefficients

	writer := NewChunkWriter(logger, client, bucket, fragmentSize)
	reader := NewChunkReader(logger, nil, client, bucket, make([]uint32, 0))

	expectedValues := make(map[corev2.BlobKey][]*encoding.Proof)

	// Write data
	for i := 0; i < 100; i++ {
		key := corev2.BlobKey(tu.RandomBytes(32))

		proofs := getProofs(t, rand.Intn(100)+100)
		expectedValues[key] = proofs

		err := writer.PutChunkProofs(context.Background(), key, proofs)
		require.NoError(t, err)
	}

	// Read data
	for key, expectedProofs := range expectedValues {
		proofs, err := reader.GetChunkProofs(context.Background(), key)
		require.NoError(t, err)
		require.Equal(t, expectedProofs, proofs)
	}
}

func TestRandomProofs(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		require.NoError(t, err)

		client, err := builder.build()
		require.NoError(t, err)
		RandomProofsTest(t, client)

		err = builder.finish()
		require.NoError(t, err)
	}
}

func generateRandomFrames(t *testing.T, encoder *rs.Encoder, size int) []*rs.Frame {
	frames, _, err := encoder.EncodeBytes(codec.ConvertByPaddingEmptyByte(tu.RandomBytes(size)))
	result := make([]*rs.Frame, len(frames))
	require.NoError(t, err)

	for i := 0; i < len(frames); i++ {
		result[i] = &frames[i]
	}

	return result
}

func RandomCoefficientsTest(t *testing.T, client s3.Client) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	chunkSize := uint64(rand.Intn(1024) + 100)
	fragmentSize := int(chunkSize / 2)

	params := encoding.ParamsFromSysPar(3, 1, chunkSize)
	encoder, _ := rs.NewEncoder(params, true)

	n := uint8(math.Log2(float64(encoder.NumEvaluations())))
	if encoder.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * encoder.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
		Fs:             fs,
		EncodingParams: params,
	}

	encoder.Computer = RsComputeDevice
	require.NotNil(t, encoder)

	writer := NewChunkWriter(logger, client, bucket, fragmentSize)
	reader := NewChunkReader(logger, nil, client, bucket, make([]uint32, 0))

	expectedValues := make(map[corev2.BlobKey][]*rs.Frame)
	metadataMap := make(map[corev2.BlobKey]*encoding.FragmentInfo)

	// Write data
	for i := 0; i < 100; i++ {
		key := corev2.BlobKey(tu.RandomBytes(32))

		coefficients := generateRandomFrames(t, encoder, int(chunkSize))
		expectedValues[key] = coefficients

		metadata, err := writer.PutChunkCoefficients(context.Background(), key, coefficients)
		require.NoError(t, err)
		metadataMap[key] = metadata
	}

	// Read data
	for key, expectedCoefficients := range expectedValues {
		coefficients, err := reader.GetChunkCoefficients(context.Background(), key, metadataMap[key])
		require.NoError(t, err)
		require.Equal(t, len(expectedCoefficients), len(coefficients))
		for i := 0; i < len(expectedCoefficients); i++ {
			require.Equal(t, *expectedCoefficients[i], *coefficients[i])
		}
	}
}

func TestRandomCoefficients(t *testing.T) {
	tu.InitializeRandom()
	for _, builder := range clientBuilders {
		err := builder.start()
		require.NoError(t, err)

		client, err := builder.build()
		require.NoError(t, err)
		RandomCoefficientsTest(t, client)

		err = builder.finish()
		require.NoError(t, err)
	}
}
