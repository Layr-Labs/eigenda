package correctness

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// Tests the basic dispersal workflow:
// - disperse a blob
// - wait for it to be confirmed
// - read the blob from the relays
// - read the blob from the validators
func testBasicDispersal(
	t *testing.T,
	payload []byte,
	certVerifierAddress string,
) error {
	if certVerifierAddress == "" {
		t.Skip("Requested cert verifier address is not configured")
	}

	c := client.GetTestClient(t, client.PreprodEnv)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err := c.DisperseAndVerify(ctx, certVerifierAddress, payload)
	if err != nil {
		return fmt.Errorf("failed to disperse and verify: %v", err)
	}

	return nil
}

// Disperse a 0 byte blob.
// Empty blobs are not allowed by the disperser
func TestEmptyBlobDispersal(t *testing.T) {
	blobBytes := []byte{}
	quorums := []core.QuorumID{0, 1}

	c := client.GetTestClient(t, client.PreprodEnv)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// We have to use the disperser client directly, since it's not possible for the PayloadDisperser to
	// attempt dispersal of an empty blob
	// This should fail with "data is empty" error
	_, _, err := c.GetDisperserClient().DisperseBlob(ctx, blobBytes, 0, quorums)
	require.Error(t, err)
	require.ErrorContains(t, err, "blob size must be greater than 0")
}

// Disperse an empty payload. Blob will not be empty, since payload encoding entails adding bytes
func TestEmptyPayloadDispersal(t *testing.T) {
	payload := []byte{}

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a payload that consists only of 0 bytes
func TestZeroPayloadDispersal(t *testing.T) {
	payload := make([]byte, 1000)

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a blob that consists only of 0 bytes. This should be permitted by eigenDA, even
// though it's not permitted by the default payload -> blob encoding scheme
func TestZeroBlobDispersal(t *testing.T) {
	blobBytes := make([]byte, 1000)
	quorums := []core.QuorumID{0, 1}

	c := client.GetTestClient(t, client.PreprodEnv)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// We have to use the disperser client directly, since it's not possible for the PayloadDisperser to
	// attempt dispersal of a blob containing all 0s
	_, _, err := c.GetDisperserClient().DisperseBlob(ctx, blobBytes, 0, quorums)
	require.NoError(t, err)
}

// Disperse a 1 byte payload (no padding).
func TestMicroscopicBlobDispersal(t *testing.T) {
	payload := []byte{1}

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a 1 byte payload (with padding).
func TestMicroscopicBlobDispersalWithPadding(t *testing.T) {
	payload := []byte{1}

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB).
func TestSmallBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(100*units.KiB, 200*units.KiB)

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a medium payload (between 1MB and 2MB).
func TestLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom()

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)
	maxBlobSize := int(config.MaxBlobSize)

	payload := rand.VariableBytes(maxBlobSize/2, maxBlobSize*3/4)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB) with each of the defined quorum sets available
func TestSmallBlobDispersalAllQuorumsSets(t *testing.T) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1_2)
	require.NoError(t, err)
	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums2)
	require.NoError(t, err)
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func TestMaximumSizedBlobDispersal(t *testing.T) {
	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	maxPermissibleDataLength, err := codec.GetMaxPermissiblePayloadLength(uint32(config.MaxBlobSize) / encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	rand := random.NewTestRandom()
	payload := rand.Bytes(int(maxPermissibleDataLength))

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.NoError(t, err)
}

// Disperse a blob that is too large (>16MB after padding)
func TestTooLargeBlobDispersal(t *testing.T) {
	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	maxPermissibleDataLength, err := codec.GetMaxPermissiblePayloadLength(uint32(config.MaxBlobSize) / encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	rand := random.NewTestRandom()
	payload := rand.Bytes(int(maxPermissibleDataLength) + 1)

	err = testBasicDispersal(t, payload, config.EigenDACertVerifierAddressQuorums0_1)
	require.Error(t, err)
}

func TestDoubleDispersal(t *testing.T) {

	t.Skip("This test is not working ever since we removed the salt param from the top level client.")

	rand := random.NewTestRandom()
	c := client.GetTestClient(t, client.PreprodEnv)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	err = c.DisperseAndVerify(ctx, config.EigenDACertVerifierAddressQuorums0_1, payload)
	require.NoError(t, err)

	// disperse again
	err = c.DisperseAndVerify(ctx, config.EigenDACertVerifierAddressQuorums0_1, payload)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob already exists"))
}

func TestUnauthorizedGetChunks(t *testing.T) {
	rand := random.NewTestRandom()
	c := client.GetTestClient(t, client.PreprodEnv)
	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	eigenDACert, err := c.DispersePayload(ctx, config.EigenDACertVerifierAddressQuorums0_1, payload)
	require.NoError(t, err)

	blobKey, err := eigenDACert.ComputeBlobKey()
	require.NoError(t, err)

	targetRelay := eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys[0]

	chunkRequests := make([]*clients.ChunkRequestByRange, 1)
	chunkRequests[0] = &clients.ChunkRequestByRange{
		BlobKey: *blobKey,
		Start:   0,
		End:     1,
	}
	_, err = c.GetRelayClient().GetChunksByRange(ctx, targetRelay, chunkRequests)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get operator key: operator not found")
}

func TestDispersalWithInvalidSignature(t *testing.T) {
	quorums := []core.QuorumID{0, 1}

	rand := random.NewTestRandom()

	c := client.GetTestClient(t, client.PreprodEnv)

	// Create a dispersal client with a random key
	signer, err := auth.NewLocalBlobRequestSigner(fmt.Sprintf("%x", rand.Bytes(32)))
	require.NoError(t, err)

	signerAccountId, err := signer.GetAccountID()
	require.NoError(t, err)
	accountId := gethcommon.HexToAddress(signerAccountId)
	fmt.Printf("Account ID: %s\n", accountId.String())

	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          c.GetConfig().DisperserHostname,
		Port:              fmt.Sprintf("%d", c.GetConfig().DisperserPort),
		UseSecureGrpcFlag: true,
	}
	disperserClient, err := clients.NewDisperserClient(disperserConfig, signer, nil, nil)
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
