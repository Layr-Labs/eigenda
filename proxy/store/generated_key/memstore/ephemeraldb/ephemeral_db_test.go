package ephemeraldb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
)

const (
	testPreimage = "Four score and seven years ago"
)

func testConfig() *memconfig.SafeConfig {
	return memconfig.NewSafeConfig(
		memconfig.Config{
			MaxBlobSizeBytes: 1024 * 1024,
			BlobExpiration:   0,
			PutLatency:       0,
			GetLatency:       0,
		})
}

func TestGetSet(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := New(ctx, testConfig(), testLogger)

	testKey := []byte("bland")
	expected := []byte(testPreimage)
	err := db.InsertEntry(testKey, expected)
	require.NoError(t, err)

	actual, err := db.FetchEntry(testKey)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestExpiration(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := testConfig()
	cfg.SetBlobExpiration(10 * time.Millisecond)
	db := New(ctx, cfg, testLogger)

	preimage := []byte(testPreimage)
	testKey := []byte("bland")

	err := db.InsertEntry(testKey, preimage)
	require.NoError(t, err)

	// sleep 1 second and verify that older blob entries are removed
	time.Sleep(time.Second * 1)

	_, err = db.FetchEntry(testKey)
	require.Error(t, err)
}

func TestLatency(t *testing.T) {
	t.Parallel()

	putLatency := 1 * time.Second
	getLatency := 1 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := testConfig()
	config.SetLatencyPUTRoute(putLatency)
	config.SetLatencyGETRoute(getLatency)
	db := New(ctx, config, testLogger)

	preimage := []byte(testPreimage)
	testKey := []byte("bland")

	timeBeforePut := time.Now()
	err := db.InsertEntry(testKey, preimage)
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(timeBeforePut), putLatency)

	timeBeforeGet := time.Now()
	_, err = db.FetchEntry(testKey)
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(timeBeforeGet), getLatency)

}

func TestPutReturnsFailoverErrorConfig(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := testConfig()
	db := New(ctx, config, testLogger)
	testKey := []byte("som-key")

	err := db.InsertEntry(testKey, []byte("some-value"))
	require.NoError(t, err)

	config.SetPUTReturnsFailoverError(true)

	// failover mode should only affect Put route
	_, err = db.FetchEntry(testKey)
	require.NoError(t, err)

	err = db.InsertEntry(testKey, []byte("some-value"))
	require.ErrorIs(t, err, &api.ErrorFailover{})
}
