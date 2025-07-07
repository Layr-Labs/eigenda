package memstore

import (
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
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
		LoadG2Points:    false,
	}

	kzgVerifier, err := kzgverifier.NewVerifier(&kzgConfig, nil)
	require.NoError(t, err)

	verifier, err := verify.NewVerifier(verifierConfig, kzgVerifier, nil)
	require.NoError(t, err)

	ms, err := New(
		t.Context(),
		verifier,
		testLogger,
		getDefaultMemStoreTestConfig(),
	)

	require.NoError(t, err)

	expected := []byte(testPreimage)
	key, err := ms.Put(t.Context(), expected)
	require.NoError(t, err)

	actual, err := ms.Get(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
