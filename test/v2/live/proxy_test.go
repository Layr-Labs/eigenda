package live

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

// Disperse an empty payload. Blob will not be empty, since payload encoding entails adding bytes
func emptyPayloadProxyDispersalTest(t *testing.T, environment string) {
	var payload []byte

	c := client.GetTestClient(t, common.TestLogger(t), environment)

	err := c.DisperseAndVerifyWithProxy(t.Context(), payload)
	require.NoError(t, err)
}

func TestEmptyPayloadProxyDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			emptyPayloadProxyDispersalTest(t, environment)
		})
	}
}

// Disperse a 1 byte payload (no padding).
func microscopicBlobProxyDispersalTest(t *testing.T, environment string) {
	payload := []byte{1}

	c := client.GetTestClient(t, common.TestLogger(t), environment)

	err := c.DisperseAndVerifyWithProxy(t.Context(), payload)
	require.NoError(t, err)
}

func TestMicroscopicBlobProxyDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			microscopicBlobProxyDispersalTest(t, environment)
		})
	}
}

// Disperse a small payload (between 1KB and 2KB).
func smallBlobProxyDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	c := client.GetTestClient(t, common.TestLogger(t), environment)

	err := c.DisperseAndVerifyWithProxy(t.Context(), payload)
	require.NoError(t, err)
}

func TestSmallBlobProxyDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			smallBlobProxyDispersalTest(t, environment)
		})
	}
}

// Disperse a blob that is exactly at the maximum size after padding (16MB)
func maximumSizedBlobProxyDispersalTest(t *testing.T, environment string) {
	config, err := client.GetConfig(common.TestLogger(t), "LIVE_TEST", environment)
	require.NoError(t, err)

	maxPermissibleDataLength, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(config.MaxBlobSize) / encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	rand := random.NewTestRandom()
	payload := rand.Bytes(int(maxPermissibleDataLength))

	c := client.GetTestClient(t, common.TestLogger(t), environment)

	err = c.DisperseAndVerifyWithProxy(t.Context(), payload)
	require.NoError(t, err)
}

func TestMaximumSizedBlobProxyDispersal(t *testing.T) {
	environments, err := client.GetEnvironmentConfigPaths()
	require.NoError(t, err)

	for _, environment := range environments {
		t.Run(getEnvironmentName(environment), func(t *testing.T) {
			maximumSizedBlobProxyDispersalTest(t, environment)
		})
	}
}
