package v2

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
		SRSNumberToLoad:               2097152, // 2097152 is default in production, no need to load so much for tests
	}

	preprodLock   sync.Mutex
	preprodClient *TestClient
)

//mkdir srs
//mkdir srs/SRSTables
//wget https://srs-mainnet.s3.amazonaws.com/kzg/g1.point --output-document=./srs/g1.point
//wget https://srs-mainnet.s3.amazonaws.com/kzg/g2.point --output-document=./srs/g2.point
//wget https://srs-mainnet.s3.amazonaws.com/kzg/g2.point.powerOf2 --output-document=./srs/g2.point.powerOf2

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

// TODO: automatically download KZG points if they are not present

func getPreprodClient(t *testing.T) *TestClient {
	preprodLock.Lock()
	defer preprodLock.Unlock()

	setupFilesystem(t, preprodConfig)

	if preprodClient == nil {
		preprodClient = NewTestClient(t, preprodConfig)
	}

	return preprodClient
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}
}

// Tests the basic dispersal workflow:
// - disperse a blob
// - wait for it to be confirmed
// - read the blob from the relays
// - read the blob from the validators
func testBasicDispersal(t *testing.T, rand *random.TestRandom, payload []byte, requestedLength int, quorums []core.QuorumID) error {
	skipInCI(t)
	client := getPreprodClient(t)

	// Make sure the payload is the correct length
	fmt.Printf("requestedLength: %d, len(payload): %d\n", requestedLength, len(payload))
	require.Equal(t, requestedLength, len(payload))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.DisperseAndVerify(ctx, payload, quorums)
}

// Disperse a 0 byte payload.
// Empty blobs are not allowed by the disperser
func TestEmptyBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	payload := []byte{}
	// This should fail with "data is empty" error
	err := testBasicDispersal(t, rand, payload, 0, []core.QuorumID{0, 1})
	require.Error(t, err)
	require.ErrorContains(t, err, "data is empty")
}

// Disperse a 1 byte payload (no padding).
func TestMicroscopicBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	payload := []byte{1}
	err := testBasicDispersal(t, rand, payload, 1, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a 1 byte payload (with padding).
func TestMicroscopicBlobDispersalWithPadding(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	payload := []byte{1}
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, 2, len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, 2, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperser a payload without padding.
// This should fail with "encountered an error to convert a 32-bytes into a valid field element" error
func TestPaddingError(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	payload := rand.Bytes(33)
	err := testBasicDispersal(t, rand, payload, len(payload), []core.QuorumID{0, 1})
	require.Error(t, err, "encountered an error to convert a 32-bytes into a valid field element")
}

// Disperse a small payload (between 1KB and 2KB).
func TestSmallBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	dataLength := 1024 + rand.Intn(1024)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, calculateExpectedPaddedSize(dataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, len(paddedPayload), []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	dataLength := 1024 * (100 + rand.Intn(100))
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, calculateExpectedPaddedSize(dataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, len(paddedPayload), []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 1MB and 2MB).
func TestLargeBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	dataLength := int(1024 * 1024 * (1 + rand.Float64()))
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, calculateExpectedPaddedSize(dataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, len(paddedPayload), []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB) with a single quorum
func TestSmallBlobDispersalSingleQuorum(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	desiredDataLength := 1024 + rand.Intn(1024)
	payload := rand.Bytes(desiredDataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	require.Equal(t, calculateExpectedPaddedSize(desiredDataLength), len(paddedPayload))
	err := testBasicDispersal(t, rand, paddedPayload, len(paddedPayload), []core.QuorumID{0})
	require.NoError(t, err)
}

// TODO:(dmanc): This test is failing. "Timed out waiting for blob to be confirmed"
// Disperse a blob that is exactly at the maximum size after padding (16MB)
func TestMaximumSizedBlobDispersal(t *testing.T) {
	skipInCI(t)

	t.Skipf("2mb is the max size in preprod") // TODO

	rand := random.NewTestRandom(t)
	originalSize, err := calculateOriginalSize(16 * 1024 * 1024)
	require.NoError(t, err)
	payload := rand.Bytes(originalSize)
	padded := codec.ConvertByPaddingEmptyByte(payload)
	lengthPadded := len(padded)
	fmt.Printf("length: %d, originalSize: %d, lengthPadded: %d\n", len(payload), originalSize, lengthPadded)
	require.Equal(t, calculateExpectedPaddedSize(originalSize), lengthPadded)
	err = testBasicDispersal(t, rand, padded, lengthPadded, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a blob that is too large (>16MB after padding)
func TestTooLargeBlobDispersal(t *testing.T) {
	skipInCI(t)
	rand := random.NewTestRandom(t)
	originalSize, err := calculateOriginalSize(16*1024*1024 + 2) // 16MB + 2 bytes
	require.NoError(t, err)
	payload := rand.Bytes(originalSize)
	padded := codec.ConvertByPaddingEmptyByte(payload)
	lengthPadded := len(padded)
	fmt.Printf("length: %d, originalSize: %d, lengthPadded: %d\n", len(payload), originalSize, lengthPadded)
	require.Equal(t, calculateExpectedPaddedSize(originalSize), lengthPadded)
	err = testBasicDispersal(t, rand, padded, lengthPadded, []core.QuorumID{0, 1})
	require.Error(t, err)
	require.ErrorContains(t, err, "blob size cannot exceed 16777216 bytes")
}

// calculateExpectedPaddedSize calculates the expected size after padding
// For each complete chunk of 31 bytes, adds 1 padding byte (making it 32)
// For the remaining bytes (if any), adds 1 padding byte at the front
func calculateExpectedPaddedSize(inputSize int) int {
	if inputSize <= 0 {
		return 0
	}
	numFullChunks := inputSize / 31
	remainingBytes := inputSize % 31

	paddedSize := numFullChunks * 32
	if remainingBytes > 0 {
		paddedSize += remainingBytes + 1
	}
	return paddedSize
}

// calculateOriginalSize calculates the original size before padding, given a padded size.
// This is the inverse of calculateExpectedPaddedSize.
// Note: For invalid padded sizes (like n*32 + 1), this will return an error
func calculateOriginalSize(paddedSize int) (int, error) {
	if paddedSize <= 0 {
		return 0, fmt.Errorf("padded size must be greater than 0")
	}

	if !isValidPaddedSize(paddedSize) {
		return 0, fmt.Errorf("padded size is not valid")
	}

	remainder := paddedSize % 32
	numFullChunks := paddedSize / 32

	// Each full 32-byte chunk came from 31 original bytes
	originalFromFullChunks := numFullChunks * 31

	// For partial chunks, subtract 1 for the padding byte
	if remainder > 0 {
		return originalFromFullChunks + remainder - 1, nil
	}
	return originalFromFullChunks, nil
}

// isValidPaddedSize checks if a given size could be the result of our padding scheme.
// A valid padded size must be either:
// 1. A multiple of 32 (representing complete chunks), or
// 2. Have a remainder > 1 when divided by 32 (representing a partial chunk with at least 1 data byte)
func isValidPaddedSize(paddedSize int) bool {
	if paddedSize <= 0 {
		return false
	}

	remainder := paddedSize % 32
	return remainder == 0 || remainder > 1
}
