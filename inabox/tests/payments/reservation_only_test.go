package payments

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	integration_test "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestReservationOnlyLegacy(t *testing.T) {
	testReservationOnly(t, clientledger.ClientLedgerModeLegacy)
}

func TestReservationOnlyNewPayments(t *testing.T) {
	testReservationOnly(t, clientledger.ClientLedgerModeReservationOnly)
}

func testReservationOnly(t *testing.T, clientLedgerMode clientledger.ClientLedgerMode) {
	infraConfig := &integration_test.InfrastructureConfig{
		TemplateName:                    "testconfig-anvil.yaml",
		TestName:                        "",
		InMemoryBlobStore:               false,
		Logger:                          test.GetLogger(),
		RootPath:                        "../../../",
		UserReservationSymbolsPerSecond: 1024,
		UserOnDemandDeposit:             0,
		// choose a bin width value much lower than the default, so that we converge on the average faster
		ReservationPeriodInterval: 10,
		ClientLedgerMode:          clientLedgerMode,
	}

	infra, err := integration_test.SetupGlobalInfrastructure(infraConfig)
	require.NoError(t, err)

	testHarness, err := integration_test.NewTestHarnessWithSetup(infra)
	require.NoError(t, err)

	t.Cleanup(func() {
		testHarness.Cleanup()
		integration_test.TeardownGlobalInfrastructure(infra)
	})

	integration_test.MineAnvilBlocks(t, testHarness.RPCClient, 6)

	payloadSize := 1000
	testDuration := 1 * time.Minute
	t.Logf("Test Duration: %s", testDuration)
	t.Logf("Payload size: %d bytes", payloadSize)

	t.Run("Within reservation limits", func(t *testing.T) {
		// the reservation of 1024 symbols/second can support up .25 min size dispersals per second.
		// to account for non-determinism, disperse at half that rate, and assert no failures
		blobsPerSecond := float32(0.125)
		t.Logf("Blobs per second: %f", blobsPerSecond)

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

	t.Run("Over reservation limits", func(t *testing.T) {
		// 2x the rate of the what's permitted by the reservation
		blobsPerSecond := float32(0.5)
		t.Logf("Blobs per second: %f", blobsPerSecond)

		resultChan := SubmitPayloads(
			t,
			random.NewTestRandom(),
			testHarness.PayloadDisperser,
			blobsPerSecond,
			payloadSize,
			testDuration)

		successCount := 0
		failureCount := 0
		for err := range resultChan {
			if err != nil {
				failureCount++
			} else {
				successCount++
			}
		}

		quarter := (successCount + failureCount) / 4

		// With 2x the reservation rate, expect roughly 50% success rate
		// To account for non-determinism, weaken assertion to just >25% of each
		require.GreaterOrEqual(t, successCount, quarter, "Expected >25%% of dispersals to succeed")
		require.GreaterOrEqual(t, failureCount, quarter, "Expected >25%% of dispersals to fail")
	})
}
