package memstore

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
)

const (
	testPreimage = "Four score and seven years ago"
)

func getDefaultMemStoreTestConfig() *memconfig.SafeConfig {
	return memconfig.NewSafeConfig(memconfig.Config{
		MaxBlobSizeBytes: 1024 * 1024,
		BlobExpiration:   0,
		PutLatency:       0,
		GetLatency:       0,
	})
}

func TestGetSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	verifierConfig := &verify.Config{
		VerifyCerts: false,
	}

	kzgConfig := kzg.KzgConfig{
		G1Path:          "../../../resources/g1.point",
		G2Path:          "../../../resources/g2.point",
		G2TrailingPath:  "../../../resources/g2.trailing.point",
		CacheDir:        "../../../resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	verifier, err := verify.NewVerifier(verifierConfig, kzgConfig, nil)
	require.NoError(t, err)

	ms, err := New(
		ctx,
		verifier,
		testLogger,
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
