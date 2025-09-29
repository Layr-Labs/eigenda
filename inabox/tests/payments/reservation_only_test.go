package payments

import (
	"testing"
	"time"

	integration_test "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestReservationOnly(t *testing.T) {
	infraConfig := &integration_test.InfrastructureConfig{
		TemplateName:                    "testconfig-anvil.yaml",
		TestName:                        "",
		InMemoryBlobStore:               false,
		Logger:                          test.GetLogger(),
		RootPath:                        "../../../",
		UserReservationSymbolsPerSecond: 1024,
	}

	infra, err := integration_test.SetupGlobalInfrastructure(infraConfig)
	require.NoError(t, err, "Failed to setup infrastructure")

	testHarness, err := integration_test.NewTestHarnessWithSetup(infra)
	require.NoError(t, err, "Failed to create test harness")

	t.Cleanup(func() {
		testHarness.Cleanup()
		integration_test.TeardownGlobalInfrastructure(infra)
	})

	t.Run("Within limits", func(t *testing.T) {
		integration_test.MineAnvilBlocks(t, testHarness.RPCClient, 6)

		testDuration := 1 * time.Minute
		blobsPerSecond := float32(0.25)
		payloadSize := 1000

		t.Logf("Starting payment exhaustion test")
		t.Logf("Test Duration: %s", testDuration)
		t.Logf("Blobs per second: %f", blobsPerSecond)
		t.Logf("Payload size: %d bytes", payloadSize)

		resultChan := SubmitPayloads(
			t,
			random.NewTestRandom(),
			testHarness.PayloadDisperser,
			blobsPerSecond,
			payloadSize,
			testDuration)

		for err := range resultChan {
			require.NoError(t, err, "Payload submission failed")
		}
	})
}
