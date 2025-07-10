package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	wmock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestValidRequest(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	publicKey, privateKey, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

func TestInvalidRequestWrongHash(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	publicKey, privateKey, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	// Modify the request so that the hash is different
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

	_, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.Error(t, err)
}

func TestInvalidRequestWrongKey(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	publicKey, _, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0

	_, differentPrivateKey, err := rand.ECDSA()
	require.NoError(t, err)
	signature, err := SignStoreChunksRequest(differentPrivateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	_, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.Error(t, err)
}

func TestInvalidRequestInvalidDisperserID(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	publicKey0, privateKey0, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress0 := crypto.PubkeyToAddress(*publicKey0)

	// Non EigenLabs disperser is registered on-chain
	publicKey1, privateKey1, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress1 := crypto.PubkeyToAddress(*publicKey1)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress0, nil)
	chainReader.Mock.On("GetDisperserAddress", uint32(1)).Return(disperserAddress1, nil)
	chainReader.Mock.On("GetDisperserAddress", uint32(1234)).Return(
		nil, errors.New("disperser not found"))

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey0, request)
	require.NoError(t, err)
	request.Signature = signature
	hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	request.DisperserID = 1
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	hash, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	request.DisperserID = 1234
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	_, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.Error(t, err)
}

func TestKeyExpiry(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	publicKey, privateKey, err := rand.ECDSA()
	require.NoError(t, err)
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	mockChainReader := wmock.MockWriter{}
	mockChainReader.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&mockChainReader,
		10,
		time.Minute,
		start)
	require.NoError(t, err)

	// Preloading the cache should have grabbed Disperser 0's key
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Since time hasn't advanced, the authenticator shouldn't have fetched the key again
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Move time forward to just before the key expires.
	now := start.Add(59 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// The key should not yet have been fetched again.
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Move time forward to just after the key expires.
	now = now.Add(2 * time.Second)
	hash, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, now)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// The key should have been fetched again.
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", 2)
}

func TestKeyCacheSize(t *testing.T) {
	rand := random.NewTestRandom()

	start := rand.Time()

	cacheSize := rand.Intn(10) + 2

	mockChainReader := wmock.MockWriter{}
	keyMap := make(map[uint32]*ecdsa.PrivateKey, cacheSize+1)
	for i := 0; i < cacheSize+1; i++ {
		publicKey, privateKey, err := rand.ECDSA()
		require.NoError(t, err)
		disperserAddress := crypto.PubkeyToAddress(*publicKey)
		keyMap[uint32(i)] = privateKey

		mockChainReader.Mock.On("GetDisperserAddress", uint32(i)).Return(disperserAddress, nil)
	}

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&mockChainReader,
		cacheSize,
		time.Minute,
		start)
	require.NoError(t, err)

	// The authenticator will preload key 0 into the cache.
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Make a request for each key (except for the last one, which won't fit in the cache).
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
		require.NoError(t, err)
		expectedHash, err := hashing.HashStoreChunksRequest(request)
		require.NoError(t, err)
		require.Equal(t, expectedHash, hash)
	}

	// All keys should have required exactly one read except from 0, which was preloaded.
	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize)

	// Make another request for each key. None should require a read from the chain.
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
		require.NoError(t, err)
		expectedHash, err := hashing.HashStoreChunksRequest(request)
		require.NoError(t, err)
		require.Equal(t, expectedHash, hash)
	}

	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize)

	// Make a request for the last key. This should require a read from the chain and will boot key 0 from the cache.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = uint32(cacheSize)
	signature, err := SignStoreChunksRequest(keyMap[uint32(cacheSize)], request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	mockChainReader.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize+1)

	// Make another request for key 0. This should require a read from the chain.
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err = SignStoreChunksRequest(keyMap[0], request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err = authenticator.AuthenticateStoreChunksRequest(context.Background(), request, start)
	require.NoError(t, err)
	expectedHash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	mockChainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize+2)
}
