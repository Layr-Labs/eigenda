package live

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

// TODO test zero sized blob
// TODO test small blob
// TODO test max sized blob
// TODO test blob that is too large

// Disperse a small payload (between 1KB and 2KB).
func smallBlobProxyDispersalTest(t *testing.T, environment string) {
	rand := random.NewTestRandom()
	payload := rand.VariableBytes(units.KiB, 2*units.KiB)

	config, err := client.GetConfig(environment)
	require.NoError(t, err)

	c := client.GetTestClient(t, environment)
	checkAndSetCertVerifierAddress(t, c, config.EigenDACertVerifierAddressQuorums0_1)

	err = c.DisperseAndVerifyWithProxy(t.Context(), payload)
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
