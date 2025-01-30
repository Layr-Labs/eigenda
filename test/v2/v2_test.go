package v2

import (
	"context"
	"fmt"
	"github.com/docker/go-units"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
)

var (
	preprodConfig = &TestClientConfig{
		TestDataPath:                  "~/.test-v2",
		DisperserHostname:             "disperser-preprod-holesky.eigenda.xyz",
		DisperserPort:                 443,
		EthRPCURLs:                    []string{"https://ethereum-holesky-rpc.publicnode.com"},
		BLSOperatorStateRetrieverAddr: "0x93545e3b9013CcaBc31E80898fef7569a4024C0C",
		EigenDAServiceManagerAddr:     "0x54A03db2784E3D0aCC08344D05385d0b62d4F432",
		SubgraphURL:                   "https://subgraph.satsuma-prod.com/51caed8fa9cb/eigenlabs/eigenda-operator-state-preprod-holesky/version/v0.7.0/api",
		SRSOrder:                      268435456,
		SRSNumberToLoad:               2097152,
		MaxBlobSize:                   16 * units.MiB,
	}

	lock   sync.Mutex
	client *TestClient

	targetConfig = preprodConfig
)

// TODO test dispersing the same blob twice in a row
// TODO test salt 0

func setupFilesystem(t *testing.T, config *TestClientConfig) {
	// Create the test data directory if it does not exist
	err := os.MkdirAll(config.TestDataPath, 0755)
	require.NoError(t, err)

	// Create the SRS directories if they do not exist
	err = os.MkdirAll(config.path(t, SRSPath), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(config.path(t, SRSPathSRSTables), 0755)
	require.NoError(t, err)

	// If any of the srs files do not exist, download them.
	filePath := config.path(t, SRSPathG1)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g1.point"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	filePath = config.path(t, SRSPathG2)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	filePath = config.path(t, SRSPathG2PowerOf2)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point.powerOf2"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	// Check to see if the private key file exists. If not, stop the test.
	filePath = config.path(t, KeyPath)
	_, err = os.Stat(filePath)
	require.NoError(t, err,
		"private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex.",
		filePath)
}

// getClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes ~1 minute
// to read the SRS points, so it's the lesser of two evils to keep it around.
func getClient(t *testing.T) *TestClient {
	lock.Lock()
	defer lock.Unlock()

	skipInCI(t)
	setupFilesystem(t, targetConfig)

	if client == nil {
		client = NewTestClient(t, targetConfig)
	}

	return client
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

// Tests the basic dispersal workflow:
// - disperse a blob
// - wait for it to be confirmed
// - read the blob from the relays
// - read the blob from the validators
func testBasicDispersal(
	t *testing.T,
	rand *random.TestRandom,
	payload []byte,
	quorums []core.QuorumID) error {

	client := getClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err := client.DisperseAndVerify(ctx, payload, quorums, rand.Uint32())
	if err != nil {
		return fmt.Errorf("failed to disperse and verify: %v", err)
	}

	return nil
}

// Disperse a 0 byte payload.
// Empty blobs are not allowed by the disperser
func TestEmptyBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := []byte{}
	// This should fail with "data is empty" error
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.Error(t, err)
	require.ErrorContains(t, err, "data is empty")
}

// Disperse a 1 byte payload (no padding).
func TestMicroscopicBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := []byte{1}
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a 1 byte payload (with padding).
func TestMicroscopicBlobDispersalWithPadding(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := []byte{1}
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, 2, len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperser a payload without padding.
// This should fail with "encountered an error to convert a 32-bytes into a valid field element" error
func TestPaddingError(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := rand.Bytes(33)
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.Error(t, err, "encountered an error to convert a 32-bytes into a valid field element")
}

// Disperse a small payload (between 1KB and 2KB).
func TestSmallBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := 1024 + rand.Intn(1024)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := 1024 * (100 + rand.Intn(100))
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 1MB and 2MB).
func TestLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := int(rand.Uint64n(targetConfig.MaxBlobSize/2) + targetConfig.MaxBlobSize/4)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB) with a single quorum
func TestSmallBlobDispersalSingleQuorum(t *testing.T) {
	rand := random.NewTestRandom(t)
	desiredDataLength := 1024 + rand.Intn(1024)
	payload := rand.Bytes(desiredDataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0})
	require.NoError(t, err)
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func TestMaximumSizedBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := int(targetConfig.MaxBlobSize)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)[:dataLength]

	//require.Equal(t, calculateExpectedPaddedSize(dataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a blob that is too large (>16MB after padding)
func TestTooLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := int(targetConfig.MaxBlobSize) + 1
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)[:dataLength+1]

	//require.Equal(t, calculateExpectedPaddedSize(dataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob size cannot exceed"))
}

func TestDoubleDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := getClient(t)

	dataLength := 1024 + rand.Intn(1024)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	salt := rand.Uint32()
	err := c.DisperseAndVerify(ctx, paddedPayload, []core.QuorumID{0, 1}, salt)
	require.NoError(t, err)

	// disperse again
	err = c.DisperseAndVerify(ctx, paddedPayload, []core.QuorumID{0, 1}, salt)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob already exists"))
}
