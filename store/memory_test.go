package store

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

const (
	testPreimage = "Four score and seven years ago"
)

func getDefaultMemStoreTestConfig() MemStoreConfig {
	return MemStoreConfig{
		MaxBlobSizeBytes: 1024 * 1024,
		BlobExpiration:   0,
		PutLatency:       0,
		GetLatency:       0,
	}
}

func getDefaultVerifierTestConfig() *verify.Config {
	return &verify.Config{
		Verify: false,
		KzgConfig: &kzg.KzgConfig{
			G1Path:          "../resources/g1.point",
			G2PowerOf2Path:  "../resources/g2.point.powerOf2",
			CacheDir:        "../resources/SRSTables",
			SRSOrder:        3000,
			SRSNumberToLoad: 3000,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		},
	}
}

func TestGetSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	verifier, err := verify.NewVerifier(getDefaultVerifierTestConfig(), nil)
	require.NoError(t, err)

	ms, err := NewMemStore(
		ctx,
		verifier,
		log.New(),
		getDefaultMemStoreTestConfig(),
	)

	require.NoError(t, err)

	expected := []byte(testPreimage)
	key, err := ms.Put(ctx, expected)
	require.NoError(t, err)

	actual, err := ms.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestExpiration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	verifier, err := verify.NewVerifier(getDefaultVerifierTestConfig(), nil)
	require.NoError(t, err)

	memstoreConfig := getDefaultMemStoreTestConfig()
	memstoreConfig.BlobExpiration = 10 * time.Millisecond
	ms, err := NewMemStore(
		ctx,
		verifier,
		log.New(),
		memstoreConfig,
	)

	require.NoError(t, err)

	preimage := []byte(testPreimage)
	key, err := ms.Put(ctx, preimage)
	require.NoError(t, err)

	// sleep 1 second and verify that older blob entries are removed
	time.Sleep(time.Second * 1)

	_, err = ms.Get(ctx, key)
	require.Error(t, err)

}

func TestLatency(t *testing.T) {
	t.Parallel()

	putLatency := 1 * time.Second
	getLatency := 1 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	verifier, err := verify.NewVerifier(getDefaultVerifierTestConfig(), nil)
	require.NoError(t, err)

	config := getDefaultMemStoreTestConfig()
	config.PutLatency = putLatency
	config.GetLatency = getLatency
	ms, err := NewMemStore(ctx, verifier, log.New(), config)

	require.NoError(t, err)

	preimage := []byte(testPreimage)
	timeBeforePut := time.Now()
	key, err := ms.Put(ctx, preimage)
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(timeBeforePut), putLatency)

	timeBeforeGet := time.Now()
	_, err = ms.Get(ctx, key)
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(timeBeforeGet), getLatency)

}
