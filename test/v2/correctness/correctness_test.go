package correctness

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"strings"
	"testing"
	"time"

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

	c := client.GetTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err := c.DisperseAndVerify(ctx, payload, quorums, rand.Uint32())
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
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := rand.VariableBytes(100*units.KiB, 200*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a medium payload (between 1MB and 2MB).
func TestLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)
	dataLength := int(rand.Uint64n(c.Config.MaxBlobSize/2) + c.Config.MaxBlobSize/4)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a small payload (between 1KB and 2KB) with a single quorum
func TestSmallBlobDispersalSingleQuorum(t *testing.T) {
	rand := random.NewTestRandom(t)
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)
	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0})
	require.NoError(t, err)
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func TestMaximumSizedBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)
	dataLength := int(c.Config.MaxBlobSize)
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)[:dataLength]

	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.NoError(t, err)
}

// Disperse a blob that is too large (>16MB after padding)
func TestTooLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)
	dataLength := int(c.Config.MaxBlobSize) + 1
	payload := rand.Bytes(dataLength)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)[:dataLength+1]

	err := testBasicDispersal(t, rand, paddedPayload, []core.QuorumID{0, 1})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "blob size cannot exceed"))
}

func TestDoubleDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
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

func TestUnauthorizedGetChunks(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	key, err := c.DispersePayload(ctx, paddedPayload, []core.QuorumID{0, 1}, rand.Uint32())
	require.NoError(t, err)

	// Wait for blob to become certified
	cert, err := c.WaitForCertification(ctx, *key, []core.QuorumID{0, 1})
	require.NoError(t, err)

	targetRelay := cert.RelayKeys[0]

	chunkRequests := make([]*clients.ChunkRequestByRange, 1)
	chunkRequests[0] = &clients.ChunkRequestByRange{
		BlobKey: *key,
		Start:   0,
		End:     1,
	}
	_, err = c.RelayClient.GetChunksByRange(ctx, targetRelay, chunkRequests)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get operator key: operator not found")
}

func TestDispersalWithInvalidSignature(t *testing.T) {
	rand := random.NewTestRandom(t)

	c := client.GetTestClient(t)

	// Create a dispersal client with a random key
	signer, err := auth.NewLocalBlobRequestSigner(fmt.Sprintf("%x", rand.Bytes(32)))
	require.NoError(t, err)

	signerAccountId, err := signer.GetAccountID()
	require.NoError(t, err)
	accountId := gethcommon.HexToAddress(signerAccountId)
	fmt.Printf("Account ID: %s\n", accountId.String())

	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          c.Config.DisperserHostname,
		Port:              fmt.Sprintf("%d", c.Config.DisperserPort),
		UseSecureGrpcFlag: true,
	}
	disperserClient, err := clients.NewDisperserClient(disperserConfig, signer, nil, nil)
	require.NoError(t, err)

	payload := rand.VariableBytes(units.KiB, 2*units.KiB)
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, _, err = disperserClient.DisperseBlob(ctx, paddedPayload, 0, []core.QuorumID{0, 1}, rand.Uint32())
	require.Error(t, err)
	require.Contains(t, err.Error(), "error accounting blob")
}
