package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	wmock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

func TestValidRequest(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)
}

func TestInvalidRequestWrongHash(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
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

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.Error(t, err)
}

func TestInvalidRequestWrongKey(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, _ := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0

	_, differentPrivateKey := rand.ECDSA()
	signature, err := SignStoreChunksRequest(differentPrivateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.Error(t, err)
}

func TestInvalidRequestInvalidDisperserID(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey0, privateKey0 := rand.ECDSA()
	disperserAddress0 := crypto.PubkeyToAddress(*publicKey0)

	// This disperser will be loaded on chain (simulated), but will fail the valid disperser ID filter.
	publicKey1, privateKey1 := rand.ECDSA()
	disperserAddress1 := crypto.PubkeyToAddress(*publicKey1)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress0, nil)
	chainReader.Mock.On("GetDisperserAddress", uint32(1)).Return(disperserAddress1, nil)
	chainReader.Mock.On("GetDisperserAddress", uint32(1234)).Return(
		nil, errors.New("disperser not found"))

	filterCallCount := atomic.Uint32{}

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		0, /* disable auth caching */
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
	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)
	require.Equal(t, uint32(2), filterCallCount.Load())

	request.DisperserID = 1
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.Error(t, err)
	require.Equal(t, uint32(3), filterCallCount.Load())

	request.DisperserID = 1234
	signature, err = SignStoreChunksRequest(privateKey1, request)
	require.NoError(t, err)
	request.Signature = signature
	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.Error(t, err)
	require.Equal(t, uint32(4), filterCallCount.Load())
}

func TestAuthCaching(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// The first request will actually be validated.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)

	// Make some more requests. Intentionally fiddle with the hash to make them invalid if checked.
	// With auth caching, those checks won't happen until the auth timeout has passed (configured to 1 minute).
	now := start
	for i := 0; i < 60; i++ {
		request = RandomStoreChunksRequest(rand)
		request.DisperserID = 0
		signature, err = SignStoreChunksRequest(privateKey, request)
		require.NoError(t, err)
		request.Signature = signature

		request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, now)
		now = now.Add(time.Second)
		require.NoError(t, err)

		// making the same request from a different origin should cause validation to happen and for it to fail
		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "otherhost", request, now)
		require.Error(t, err)
	}

	// The next request will be made after the auth timeout has passed, so it will be validated.
	// Since it is actually invalid, the authenticator should reject it.
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err = SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, now)
	require.Error(t, err)
}

func TestAuthCachingDisabled(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		0, // This disables auth caching
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// The first request will always be validated.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)

	// Make another request without moving time forward. It should be validated.
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err = SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.Error(t, err)
}

func TestKeyExpiry(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// Preloading the cache should have grabbed Disperser 0's key
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)

	// Since time hasn't advanced, the authenticator shouldn't have fetched the key again
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Move time forward to just before the key expires.
	now := start.Add(59 * time.Second)
	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, now)
	require.NoError(t, err)

	// The key should not yet have been fetched again.
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Move time forward to just after the key expires.
	now = now.Add(2 * time.Second)
	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, now)
	require.NoError(t, err)

	// The key should have been fetched again.
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", 2)
}

func TestAuthCacheSize(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	cacheSize := rand.Intn(10) + 2

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		cacheSize,
		time.Minute,
		time.Minute,
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// Make requests from cacheSize different origins.
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = 0
		signature, err := SignStoreChunksRequest(privateKey, request)
		require.NoError(t, err)
		request.Signature = signature

		origin := fmt.Sprintf("%d", i)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), origin, request, start)
		require.NoError(t, err)
	}

	// All origins should be authenticated in the auth cache. If we send invalid requests from the same origins,
	// they should still be authenticated (since the authenticator won't re-check).
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = 0
		signature, err := SignStoreChunksRequest(privateKey, request)
		require.NoError(t, err)
		request.Signature = signature

		request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

		origin := fmt.Sprintf("%d", i)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), origin, request, start)
		require.NoError(t, err)
	}

	// Make a request from a new origin. This should boot origin 0 from the cache.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "neworigin", request, start)
	require.NoError(t, err)

	for i := 0; i < cacheSize; i++ {
		request = RandomStoreChunksRequest(rand)
		request.DisperserID = 0
		signature, err = SignStoreChunksRequest(privateKey, request)
		require.NoError(t, err)
		request.Signature = signature

		request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)

		origin := fmt.Sprintf("%d", i)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), origin, request, start)

		if i == 0 {
			// Origin 0 should have been booted from the cache, so this request should be re-validated.
			require.Error(t, err)
		} else {
			// All other origins should still be in the cache.
			require.NoError(t, err)
		}
	}
}

func TestKeyCacheSize(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	cacheSize := rand.Intn(10) + 2

	chainReader := wmock.MockWriter{}
	keyMap := make(map[uint32]*ecdsa.PrivateKey, cacheSize+1)
	for i := 0; i < cacheSize+1; i++ {
		publicKey, privateKey := rand.ECDSA()
		disperserAddress := crypto.PubkeyToAddress(*publicKey)
		keyMap[uint32(i)] = privateKey

		chainReader.Mock.On("GetDisperserAddress", uint32(i)).Return(disperserAddress, nil)
	}

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		cacheSize,
		time.Minute,
		0, // disable auth caching
		func(uint32) bool { return true },
		start)
	require.NoError(t, err)

	// The authenticator will preload key 0 into the cache.
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", 1)

	// Make a request for each key (except for the last one, which won't fit in the cache).
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		origin := fmt.Sprintf("%d", i)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), origin, request, start)
		require.NoError(t, err)
	}

	// All keys should have required exactly one read except from 0, which was preloaded.
	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize)

	// Make another request for each key. None should require a read from the chain.
	for i := 0; i < cacheSize; i++ {
		request := RandomStoreChunksRequest(rand)
		request.DisperserID = uint32(i)
		signature, err := SignStoreChunksRequest(keyMap[uint32(i)], request)
		require.NoError(t, err)
		request.Signature = signature

		origin := fmt.Sprintf("%d", i)

		err = authenticator.AuthenticateStoreChunksRequest(context.Background(), origin, request, start)
		require.NoError(t, err)
	}

	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize)

	// Make a request for the last key. This should require a read from the chain and will boot key 0 from the cache.
	request := RandomStoreChunksRequest(rand)
	request.DisperserID = uint32(cacheSize)
	signature, err := SignStoreChunksRequest(keyMap[uint32(cacheSize)], request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(),
		fmt.Sprintf("%d", cacheSize), request, start)
	require.NoError(t, err)

	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize+1)

	// Make another request for key 0. This should require a read from the chain.
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err = SignStoreChunksRequest(keyMap[0], request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "0", request, start)
	require.NoError(t, err)

	chainReader.Mock.AssertNumberOfCalls(t, "GetDisperserAddress", cacheSize+2)
}
