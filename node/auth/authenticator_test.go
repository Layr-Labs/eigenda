package auth

import (
	"crypto/ecdsa"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/hashing"
	wmock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/test/random"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestValidRequest(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, privateKey, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(0), uint32(3)).Return([]gethcommon.Address{
		disperserAddress,
	}, nil)

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&chainReader,
		10,
		time.Minute,
		3,
		func(uint32) bool { return true },
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

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(0), uint32(3)).Return([]gethcommon.Address{
		disperserAddress,
	}, nil)

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&chainReader,
		10,
		time.Minute,
		3,
		func(uint32) bool { return true },
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

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(0), uint32(3)).Return([]gethcommon.Address{
		disperserAddress,
	}, nil)

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&chainReader,
		10,
		time.Minute,
		3,
		func(uint32) bool { return true },
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

	// This disperser will be loaded on chain (simulated), but will fail the valid disperser ID filter.
	disperserAddress1, privateKey1, err := rand.EthAccount()
	require.NoError(t, err)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(0), uint32(3)).Return([]gethcommon.Address{
		disperserAddress0,
	}, nil)
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(1), uint32(3)).Return([]gethcommon.Address{
		disperserAddress1,
	}, nil)
	chainReader.Mock.On("GetAllDisperserAddresses", uint32(1234), uint32(3)).Return(
		nil, errors.New("disperser not found"))

	filterCallCount := atomic.Uint32{}

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&chainReader,
		10,
		time.Minute,
		3,
		func(id uint32) bool {
			filterCallCount.Add(1)
			return id != uint32(1)
		},
		start)
	require.NoError(t, err)
	require.Equal(t, uint32(1), filterCallCount.Load())

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
	require.Equal(t, uint32(2), filterCallCount.Load())

	request.DisperserID = 1
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.Error(t, err)

	request.DisperserID = 1234
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	_, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, start)
	require.Error(t, err)
}

func TestKeyExpiry(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	start := rand.Time()

	disperserAddress, privateKey, err := rand.EthAccount()
	require.NoError(t, err)

	mockChainReader := wmock.MockWriter{}
	mockChainReader.On("GetAllDisperserAddresses", uint32(0), uint32(3)).Return([]gethcommon.Address{
		disperserAddress,
	}, nil)

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&mockChainReader,
		10,
		time.Minute,
		3,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// Preloading the cache should have grabbed Disperser 0's key
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", 1)

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

	// Since time hasn't advanced, the authenticator shouldn't have fetched the key again
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", 1)

	// Move time forward to just before the key expires.
	now := start.Add(59 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// The key should not yet have been fetched again.
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", 1)

	// Move time forward to just after the key expires.
	now = now.Add(2 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(ctx, request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// The key should have been fetched again.
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", 2)
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

		mockChainReader.Mock.On("GetAllDisperserAddresses", uint32(i), uint32(3)).Return([]gethcommon.Address{
		disperserAddress,
	}, nil)
	}

	authenticator, err := NewRequestAuthenticator(
		ctx,
		&mockChainReader,
		cacheSize,
		time.Minute,
		3,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// The authenticator will preload key 0 into the cache.
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", 1)

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

	// All keys should have required exactly one read except from 0, which was preloaded.
	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", cacheSize)

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

	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", cacheSize)

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

	mockChainReader.AssertNumberOfCalls(t, "GetAllDisperserAddresses", cacheSize+1)

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

	mockChainReader.Mock.AssertNumberOfCalls(t, "GetAllDisperserAddresses", cacheSize+2)
}
