package payments

import (
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

// NOTE: Currently, it doesn't work to run these tests in sequence. Each test must be run as a separate command.
// The problem is that the cleanup logic sometimes randomly fails to free docker ports, so subsequent setups fail.
// Once we figure out why resources aren't being freed, then these tests will be runnable the "normal" way.

func TestReservationOnly_LegacyClient_LegacyController(t *testing.T) {
	t.Skip("Manual test for now")
	testReservationOnly(t, clientledger.ClientLedgerModeLegacy, false)
}

func TestReservationOnly_LegacyClient_NewController(t *testing.T) {
	t.Skip("Manual test for now")
	testReservationOnly(t, clientledger.ClientLedgerModeLegacy, true)
}

func TestReservationOnly_NewClient_LegacyController(t *testing.T) {
	t.Skip("Manual test for now")
	testReservationOnly(t, clientledger.ClientLedgerModeReservationOnly, false)
}

func TestReservationOnly_NewClient_NewController(t *testing.T) {
	t.Skip("Manual test for now")
	testReservationOnly(t, clientledger.ClientLedgerModeReservationOnly, true)
}

func testReservationOnly(t *testing.T, clientLedgerMode clientledger.ClientLedgerMode, controllerUseNewPayments bool) {
	// Save current working directory. The setup process in its current form changes working directory, which causes
	// subsequent executions to fail, since the process relies on relative paths. This is a workaround for now: we just
	// capture the original working directory, and switch back to it as a cleanup step.
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	})

	infraConfig := &integration.InfrastructureConfig{
		TemplateName:                    "testconfig-anvil.yaml",
		TestName:                        "",
		InMemoryBlobStore:               false,
		Logger:                          test.GetLogger(),
		RootPath:                        "../../../",
		UserReservationSymbolsPerSecond: 1024,
		ClientLedgerMode:                clientLedgerMode,
		ControllerUseNewPayments:        controllerUseNewPayments,
	}

	infra, err := integration.SetupInfrastructure(infraConfig)
	if infra != nil {
		t.Cleanup(func() {
			integration.TeardownInfrastructure(infra)
		})
	}
	require.NoError(t, err)

	testHarness, err := integration.NewTestHarnessWithSetup(infra)
	if testHarness != nil {
		t.Cleanup(func() {
			testHarness.Cleanup()
		})
	}
	require.NoError(t, err)

	integration.MineAnvilBlocks(t, testHarness.RPCClient, 6)

	payloadSize := 1000
	testDuration := 1 * time.Minute
	t.Logf("Test Duration: %s", testDuration)
	t.Logf("Payload size: %d bytes", payloadSize)

	t.Run("Within reservation limits", func(t *testing.T) {
		initialReservation, err := reservation.NewReservation(
			1024,
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		err = testHarness.UpdateReservationOnChain(t.Context(), t, initialReservation)
		require.NoError(t, err)

		// Wait for vault monitor to pick up changes
		time.Sleep(3 * time.Second)

		testRandom := random.NewTestRandom()
		// the reservation of 1024 symbols/second can support up .25 min size dispersals per second.
		// to account for non-determinism, disperse at half that rate, and assert no failures
		blobsPerSecond := float32(0.125)
		t.Logf("Blobs per second: %f", blobsPerSecond)

		resultChan := submitPayloads(
			t,
			testRandom,
			testHarness.PayloadDisperser,
			blobsPerSecond,
			payloadSize,
			testDuration)

		for err := range resultChan {
			require.NoError(t, err, "Payload submission failed")
		}

		// The next part of the test decreases the reservation size, and asserts the the same dispersal conditions now
		// yield errors. The legacy client ledger mode doesn't observe payment vault updates, so skip this next
		// part if that's the current configuration
		if clientLedgerMode == clientledger.ClientLedgerModeLegacy {
			return
		}

		newReservation, err := reservation.NewReservation(
			256, // this rate will not support the 0.125 dispersals/second rate
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		err = testHarness.UpdateReservationOnChain(t.Context(), t, newReservation)
		require.NoError(t, err)

		// the vault monitor checks every 1 second, so this should be plenty of time
		time.Sleep(3 * time.Second)

		t.Log("Dispersing with decreased reservation limits")
		resultChan = submitPayloads(
			t,
			testRandom,
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

		// With reduced reservation rate, expect roughly 50% success rate
		// To account for non-determinism, weaken assertion to just >25% of each
		require.GreaterOrEqual(t, successCount, quarter, "Expected >25%% of dispersals to succeed")
		require.GreaterOrEqual(t, failureCount, quarter, "Expected >25%% of dispersals to fail")
	})

	t.Run("Over reservation limits", func(t *testing.T) {
		initialReservation, err := reservation.NewReservation(
			1024,
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		err = testHarness.UpdateReservationOnChain(t.Context(), t, initialReservation)
		require.NoError(t, err)

		// Wait for vault monitor to pick up changes
		time.Sleep(3 * time.Second)

		testRandom := random.NewTestRandom()
		// 2x the rate of the what's permitted by the reservation
		blobsPerSecond := float32(0.5)
		t.Logf("Blobs per second: %f", blobsPerSecond)

		resultChan := submitPayloads(
			t,
			testRandom,
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

		// The next part of the test increases the reservation size, and asserts the the same dispersal conditions no
		// longer yield errors. The legacy client ledger mode doesn't observe payment vault updates, so skip this next
		// part if that's the current configuration
		if clientLedgerMode == clientledger.ClientLedgerModeLegacy {
			return
		}

		newReservation, err := reservation.NewReservation(
			4096, // this rate will easily support the 0.5 dispersals/second rate
			time.Now().Add(-1*time.Hour),
			time.Now().Add(24*time.Hour),
			[]core.QuorumID{0, 1},
		)
		require.NoError(t, err)
		err = testHarness.UpdateReservationOnChain(t.Context(), t, newReservation)
		require.NoError(t, err)

		// the vault monitor checks every 1 second, so this should be plenty of time
		time.Sleep(3 * time.Second)

		t.Log("Testing with increased reservation limits")
		resultChan = submitPayloads(
			t,
			testRandom,
			testHarness.PayloadDisperser,
			blobsPerSecond,
			payloadSize,
			testDuration)

		// with the updated reservation, we should expect no failures at all
		for err := range resultChan {
			require.NoError(t, err, "Payload submission failed")
		}
	})
}
