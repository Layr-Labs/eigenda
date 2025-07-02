package live

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// getEnvironmentName takes an environment string as listed in environments (aka a path to a config file describing
// the environment) and returns the name of the environment. Assumes the path is in the format of
// "path/to/ENVIRONMENT_NAME.json".
func getEnvironmentName(environment string) string {
	elements := strings.Split(environment, "/")
	fileName := elements[len(elements)-1]
	environmentName := strings.Split(fileName, ".")[0]
	return environmentName
}

// checkAndSetCertVerifierAddress checks whether the input address string is empty, and skips the test if it is
//
// If the input address is not empty, this method configures the test client to use the input address as the cert
// verifier address.
func checkAndSetCertVerifierAddress(t *testing.T, c *client.TestClient, certVerifierAddress string) {
	if certVerifierAddress == "" {
		t.Skip("Requested cert verifier address is not configured")
	}
	c.SetCertVerifierAddress(certVerifierAddress)
}

// Tests the basic dispersal workflow:
// - disperse a blob
// - wait for it to be confirmed
// - read the blob from the relays
// - read the blob from the validators
func testBasicDispersal(c *client.TestClient, payload []byte) error {
	err := c.DisperseAndVerify(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("failed to disperse and verify: %v", err)
	}

	return nil
}

// Disperse a 0 byte blob.
// Empty blobs are not allowed by the disperser
func emptyBlobDispersalTest(t *testing.T, environment string) {
	blobBytes := []byte{}
	quorums := []core.QuorumID{0, 1}

	c := client.GetTestClient(t, environment)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// We have to use the disperser client directly, since it's not possible for the PayloadDisperser to
	// attempt dispersal of an empty blob
	// This should fail with "data is empty" error
	_, _, err := c.GetDisperserClient().DisperseBlob(ctx, blobBytes, 0, quorums)
	require.Error(t, err)
	require.ErrorContains(t, err, "zero symbols requested")
}

func TestEmptyBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			emptyBlobDispersalTest(t, environment)
		})
	}
}

// Disperse an empty payload. Blob will not be empty, since payload encoding entails adding bytes
func emptyPayloadDispersalTest(t *testing.T, environment string) {
	payload := []byte{}

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestEmptyPayloadDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			emptyPayloadDispersalTest(t, environment)
		})
	}
}

// Disperse a payload that consists only of 0 bytes
func testZeroPayloadDispersalTest(t *testing.T, environment string) {
	payload := make([]byte, 1000)

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestZeroPayloadDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			testZeroPayloadDispersalTest(t, environment)
		})
	}
}

// Disperse a blob that consists only of 0 bytes. This should be permitted by eigenDA, even
// though it's not permitted by the default payload -> blob encoding scheme
func zeroBlobDispersalTest(t *testing.T, environment string) {
	blobBytes := make([]byte, 1000)
	quorums := []core.QuorumID{0, 1}

	c := client.GetTestClient(t, environment)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// We have to use the disperser client directly, since it's not possible for the PayloadDisperser to
	// attempt dispersal of a blob containing all 0s
	_, _, err := c.GetDisperserClient().DisperseBlob(ctx, blobBytes, 0, quorums)
	require.NoError(t, err)
}

func TestZeroBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			zeroBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a 1 byte payload (no padding).
func microscopicBlobDispersalTest(t *testing.T, environment string) {
	payload := []byte{1}

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestMicroscopicBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			microscopicBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a 1 byte payload (with padding).
func microscopicBlobDispersalWithPadding(t *testing.T, environment string) {
	payload := []byte{1}

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestMicroscopicBlobDispersalWithPadding(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			microscopicBlobDispersalWithPadding(t, environment)
		})
	}
}

// Disperse a small payload (between 1KB and 2KB).
func smallBlobDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestSmallBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			smallBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a medium payload (between 100KB and 200KB).
func mediumBlobDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(100*units.KiB, 200*units.KiB)

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestMediumBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			mediumBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a medium payload (between 1MB and 2MB).
func largeBlobDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()

	config, err := client.GetConfig(environment)
	require.NoError(t, err)
	maxBlobSize := int(config.MaxBlobSize)

	payload := rand.VariableBytes(maxBlobSize/2, maxBlobSize*3/4)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestLargeBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			largeBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a small payload (between 1KB and 2KB) with each of the defined quorum sets available
func smallBlobDispersalAllQuorumsSetsTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)

	t.Run("0 1", func(t *testing.T) {
		checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)
		err = testBasicDispersal(c, payload)
		require.NoError(t, err)
	})

	t.Run("0 1 2", func(t *testing.T) {
		checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1_2)
		err = testBasicDispersal(c, payload)
		require.NoError(t, err)
	})

	t.Run("2", func(t *testing.T) {
		checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums2)
		err = testBasicDispersal(c, payload)
		require.NoError(t, err)
	})
}

func TestSmallBlobDispersalAllQuorumsSets(t *testing.T) {
	t.Skip() // currently broken

	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			smallBlobDispersalAllQuorumsSetsTest(t, environment)
		})
	}
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func maximumSizedBlobDispersalTest(t *testing.T, environment string) {
	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	maxPermissibleDataLength, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(config.MaxBlobSize) / encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	rand := random.NewTestRandom()
	payload := rand.Bytes(int(maxPermissibleDataLength))

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.NoError(t, err)
}

func TestMaximumSizedBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			maximumSizedBlobDispersalTest(t, environment)
		})
	}
}

// Disperse a blob that is too large (>16MB after padding)
func tooLargeBlobDispersalTest(t *testing.T, environment string) {
	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	maxPermissibleDataLength, err := codec.BlobSymbolsToMaxPayloadSize(uint32(config.MaxBlobSize) / encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	rand := random.NewTestRandom()
	payload := rand.Bytes(int(maxPermissibleDataLength) + 1)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = testBasicDispersal(c, payload)
	require.Error(t, err)
}

func TestTooLargeBlobDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			tooLargeBlobDispersalTest(t, environment)
		})
	}
}

func doubleDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	c := client.GetTestClient(t, environment)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c.SetCertVerifierAddress(config.EigenDACertVerifierAddressQuorums0_1)

	err = c.DisperseAndVerify(ctx, payload)
	require.NoError(t, err)

	// disperse again
	err = c.DisperseAndVerify(ctx, payload)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob already exists"))
}

func TestDoubleDispersal(t *testing.T) {
	t.Skip("This test is not working ever since we removed the salt param from the top level client.")

	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			doubleDispersalTest(t, environment)
		})
	}
}

func unauthorizedGetChunksTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	c := client.GetTestClient(t, environment)
	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c.SetCertVerifierAddress(config.EigenDACertVerifierAddressQuorums0_1)

	eigenDACert, err := c.DispersePayload(ctx, payload)
	require.NoError(t, err)

	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	require.True(t, ok, "expected EigenDACertV3, got %T", eigenDACert)
	require.NotNil(t, eigenDAV3Cert)

	blobKey, err := eigenDAV3Cert.ComputeBlobKey()
	require.NoError(t, err)

	targetRelay := eigenDAV3Cert.RelayKeys()[0]

	chunkRequests := make([]*relay.ChunkRequestByRange, 1)
	chunkRequests[0] = &relay.ChunkRequestByRange{
		BlobKey: *blobKey,
		Start:   0,
		End:     1,
	}
	_, err = c.GetRelayClient().GetChunksByRange(ctx, targetRelay, chunkRequests)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get operator key: operator not found")
}

func TestUnauthorizedGetChunks(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			unauthorizedGetChunksTest(t, environment)
		})
	}
}

func dispersalWithInvalidSignatureTest(t *testing.T, environment string) {
	quorums := []core.QuorumID{0, 1}

	rand := random.NewTestRandom()

	c := client.GetTestClient(t, environment)

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	require.NoError(t, err)

	// Create a dispersal client with a random key
	signer, err := auth.NewLocalBlobRequestSigner(fmt.Sprintf("%x", rand.Bytes(32)))
	require.NoError(t, err)

	accountId, err := signer.GetAccountID()
	require.NoError(t, err)
	fmt.Printf("Account ID: %s\n", accountId.Hex())

	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          c.GetConfig().DisperserHostname,
		Port:              fmt.Sprintf("%d", c.GetConfig().DisperserPort),
		UseSecureGrpcFlag: true,
	}
	disperserClient, err := clients.NewDisperserClient(logger, disperserConfig, signer, nil, nil)
	require.NoError(t, err)

	payloadBytes := rand.VariableBytes(units.KiB, 2*units.KiB)

	payload := coretypes.NewPayload(payloadBytes)

	// TODO (litt3): make the blob form configurable. Using PolynomialFormCoeff means that the data isn't being
	//  FFTed/IFFTed, and it is important for both modes of operation to be tested.
	blob, err := payload.ToBlob(codecs.PolynomialFormCoeff)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, _, err = disperserClient.DisperseBlob(ctx, blob.Serialize(), 0, quorums)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error accounting blob")
}

func TestDispersalWithInvalidSignature(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			dispersalWithInvalidSignatureTest(t, environment)
		})
	}
}
