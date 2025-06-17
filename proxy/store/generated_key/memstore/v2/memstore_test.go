package memstore

import (
	"context"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
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
	g1Srs, err := kzg.ReadG1Points("../../../../resources/g1.point", 3000, 2)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, err)

	msV2, err := New(
		ctx,
		testLogger,
		getDefaultMemStoreTestConfig(),
		g1Srs,
	)

	require.NoError(t, err)

	expected := []byte(testPreimage)
	key, err := msV2.Put(ctx, expected)
	require.NoError(t, err)

	cert := certs.NewVersionedCert(key, coretypes.VersionThreeCert)

	actual, err := msV2.Get(ctx, cert)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
