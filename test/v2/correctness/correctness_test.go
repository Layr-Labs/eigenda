package correctness

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
)

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

	c := client.GetTestClient(t, client.PreprodEnv)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err := c.DisperseAndVerify(ctx, quorums, payload, rand.Uint32())
	if err != nil {
		return fmt.Errorf("failed to disperse and verify: %v", err)
	}

	return nil
}

// Disperse a 0 byte blob.
// Empty blobs are not allowed by the disperser
func TestEmptyBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	blobBytes := []byte{}
	quorums := []core.QuorumID{0, 1}

	c := client.GetTestClient(t, client.PreprodEnv)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// We have to use the disperser client directly, since it's not possible for the PayloadDisperser to
	// attempt dispersal of an empty blob
	// This should fail with "data is empty" error
	_, _, err := c.GetDisperserClient().DisperseBlob(ctx, blobBytes, 0, quorums, rand.Uint32())
	require.Error(t, err)
	require.ErrorContains(t, err, "blob size must be greater than 0")
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
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB).
func TestSmallBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := rand.VariableBytes(100*units.KiB, 200*units.KiB)
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 1MB and 2MB).
func TestLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)
	maxBlobSize := int(config.MaxBlobSize)

	payload := rand.VariableBytes(maxBlobSize/2, maxBlobSize*3/4)

	err = testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB) with a single quorum
func TestSmallBlobDispersalSingleQuorum(t *testing.T) {
	t.Skip("TODO: validation is borked for single quorum dispersal")

	rand := random.NewTestRandom(t)
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	err := testBasicDispersal(t, rand, payload, []core.QuorumID{0})
	require.NoError(t, err)
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func TestMaximumSizedBlobDispersal(t *testing.T) {
	t.Skip("it's really hard to figure out what the maximum payload size is, re-enable when that is resolved")

	quorums := []core.QuorumID{0, 1}

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)
	maxBlobSize := int(config.MaxBlobSize)
	dataLength := maxBlobSize

	rand := random.NewTestRandom(t)
	payload := rand.Bytes(dataLength)
	err = testBasicDispersal(t, rand, payload, quorums)
	require.NoError(t, err)
}

// Disperse a blob that is too large (>16MB after padding)
func TestTooLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	// TODO refactor this to use exactly 1 byte more than max size after padding and header data

	config, err := client.GetConfig(client.PreprodEnv)
	require.NoError(t, err)
	maxBlobSize := int(config.MaxBlobSize)

	dataLength := maxBlobSize + 1
	payload := rand.Bytes(dataLength)

	err = testBasicDispersal(t, rand, payload, []core.QuorumID{0, 1})
	require.Error(t, err)
	fmt.Println(err)
}

func TestDoubleDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t, client.PreprodEnv)

	quorums := []core.QuorumID{0, 1}
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	salt := rand.Uint32()
	err := c.DisperseAndVerify(ctx, quorums, payload, salt)
	require.NoError(t, err)

	// disperse again
	err = c.DisperseAndVerify(ctx, quorums, payload, salt)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob already exists"))
}

func TestUnauthorizedGetChunks(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t, client.PreprodEnv)

	quorums := []core.QuorumID{0, 1}
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	eigenDACert, err := c.DispersePayload(ctx, quorums, payload, rand.Uint32())
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

	rand := random.NewTestRandom(t)

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

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, _, err = disperserClient.DisperseBlob(ctx, paddedPayload, 0, quorums, rand.Uint32())
	require.Error(t, err)
	require.Contains(t, err.Error(), "error accounting blob")
}
