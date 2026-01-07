package auth

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/core"
	wmock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// setupMockChainReader sets up a mock chain reader with the given disperser addresses.
func setupMockChainReader(dispersers map[uint32]gethcommon.Address) *wmock.MockWriter {
	chainReader := &wmock.MockWriter{}

	for id, addr := range dispersers {
		chainReader.Mock.On("GetDisperserAddress", id).Return(addr, nil)
	}

	return chainReader
}

func TestValidRequest(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, privateKey, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperserAddress,
	})

	authenticator, err := NewRequestAuthenticator(
		ctx,
		chainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

func TestInvalidRequestWrongHash(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, privateKey, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperserAddress,
	})

	authenticator, err := NewRequestAuthenticator(
		ctx,
		chainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	// Modify the request so that the hash is different
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.Error(t, err)
}

func TestInvalidRequestWrongKey(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, _, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperserAddress,
	})

	authenticator, err := NewRequestAuthenticator(
		ctx,
		chainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0

	_, differentPrivateKey, err := rand.EthAccount()
	require.NoError(t, err)
	signature, err := SignStoreChunksRequest(differentPrivateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.Error(t, err)
}

func TestInvalidRequestInvalidDisperserID(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress0, privateKey0, err := rand.EthAccount()
	require.NoError(t, err)

	disperserAddress1, privateKey1, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperserAddress0,
		1: disperserAddress1,
	})
	// Add specific mock for disperser ID 1234 which should return an error
	chainReader.Mock.On("GetDisperserAddress", uint32(1234)).Return(
		gethcommon.Address{}, errors.New("disperser not found"))

	authenticator, err := NewRequestAuthenticator(
		ctx,
		chainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	// Test valid disperser ID 0
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey0, request)
	require.NoError(t, err)
	request.Signature = signature
	hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Test valid disperser ID 1 (should work now that we accept all disperser IDs)
	request.DisperserID = 1
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err) // Should succeed now

	// Test invalid disperser ID (not found on chain)
	request.DisperserID = 1234
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.Error(t, err) // Should still fail - disperser not found
}

func TestKeyExpiry(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, privateKey, err := rand.EthAccount()
	require.NoError(t, err)

	mockChainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperserAddress,
	})

	authenticator, err := NewRequestAuthenticator(
		ctx,
		mockChainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Move time forward to just before the key expires.
	now := start.Add(59 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Move time forward to just after the key expires.
	now = now.Add(2 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

func TestKeyCacheSize(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	cacheSize := rand.Intn(10) + 2

	mockChainReader := wmock.MockWriter{}
	keyMap := make(map[uint32]*ecdsa.PrivateKey, cacheSize+1)
	for i := 0; i < cacheSize+1; i++ {
		disperserAddress, privateKey, err := rand.EthAccount()
		require.NoError(t, err)
		keyMap[uint32(i)] = privateKey

		mockChainReader.Mock.On("GetDisperserAddress", uint32(i)).Return(disperserAddress, nil)
	}

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&mockChainReader,
		test.GetLogger(),
		cacheSize,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	// Make a request for each key (except for the last one, which won't fit in the cache).
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
		require.NoError(t, err)
		expectedHash, err := hashing.HashStoreChunksRequest(request)
		require.NoError(t, err)
		require.Equal(t, expectedHash, hash)
	}

	// Make another request for each key. None should require a read from the chain.
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
		require.NoError(t, err)
		expectedHash, err := hashing.HashStoreChunksRequest(request)
		require.NoError(t, err)
		require.Equal(t, expectedHash, hash)
	}

	// Make a request for the last key. This should require a read from the chain and will boot key 0 from the cache.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = uint32(cacheSize)
	signature, err := SignStoreChunksRequest(keyMap[uint32(cacheSize)], request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Make another request for key 0. This should require a read from the chain.
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err = SignStoreChunksRequest(keyMap[0], request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.NoError(t, err)
	expectedHash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

func TestOnDemandPaymentAuthorization(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperser0Address, _, err := rand.EthAccount()
	require.NoError(t, err)

	disperser1Address, _, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := setupMockChainReader(map[uint32]gethcommon.Address{
		0: disperser0Address,
		1: disperser1Address,
	})

	authenticator, err := NewRequestAuthenticator(
		ctx,
		chainReader,
		test.GetLogger(),
		10,
		time.Minute,
		[]uint32{0},
		start)
	require.NoError(t, err)

	onDemandBatch := &corev2.Batch{
		BlobCertificates: []*corev2.BlobCertificate{
			{BlobHeader: &corev2.BlobHeader{PaymentMetadata: core.PaymentMetadata{CumulativePayment: big.NewInt(10)}}},
			{BlobHeader: &corev2.BlobHeader{PaymentMetadata: core.PaymentMetadata{CumulativePayment: big.NewInt(0)}}},
		},
	}

	reservationBatch := &corev2.Batch{
		BlobCertificates: []*corev2.BlobCertificate{
			{BlobHeader: &corev2.BlobHeader{PaymentMetadata: core.PaymentMetadata{CumulativePayment: big.NewInt(0)}}},
		},
	}

	require.True(t, authenticator.IsDisperserAuthorized(0, onDemandBatch))
	require.True(t, authenticator.IsDisperserAuthorized(0, reservationBatch))

	require.False(t, authenticator.IsDisperserAuthorized(1, onDemandBatch))
	require.True(t, authenticator.IsDisperserAuthorized(1, reservationBatch))
}

func TestMultipleDisperserIDs(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	// Set up multiple disperser addresses
	disperser0Address, privateKey0, err := rand.EthAccount()
	require.NoError(t, err)
	disperser1Address, privateKey1, err := rand.EthAccount()
	require.NoError(t, err)
	disperser2Address, privateKey2, err := rand.EthAccount()
	require.NoError(t, err)

	mockChainReader := wmock.MockWriter{}
	mockChainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperser0Address, nil)
	mockChainReader.Mock.On("GetDisperserAddress", uint32(1)).Return(disperser1Address, nil)
	mockChainReader.Mock.On("GetDisperserAddress", uint32(2)).Return(disperser2Address, nil)

	// Create authenticator with cache size 3 to test preloading
	authenticator, err := NewRequestAuthenticator(
		ctx,
		&mockChainReader,
		test.GetLogger(),
		3,
		time.Minute,
		[]uint32{0}, // Only disperser 0 authorized for on-demand
		start)
	require.NoError(t, err)

	// Test authentication with different disperser IDs
	testCases := []struct {
		disperserID uint32
		privateKey  *ecdsa.PrivateKey
	}{
		{0, privateKey0},
		{1, privateKey1},
		{2, privateKey2},
	}

	for _, tc := range testCases {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = tc.disperserID
		signature, err := SignStoreChunksRequest(tc.privateKey, request)
		require.NoError(t, err)
		request.Signature = signature

		hash, err := authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
		require.NoError(t, err)
		require.NotNil(t, hash)
	}
}
