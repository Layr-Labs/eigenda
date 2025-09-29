package payments

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	integration_test "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestPaymentExhaustion(t *testing.T) {
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

	t.Run("PaymentExhaustion", func(t *testing.T) {
		integration_test.MineAnvilBlocks(t, testHarness.RPCClient, 6)

		testDuration := 1 * time.Minute
		blobsPerSecond := float32(1)
		payloadSize := 1000

		t.Logf("Starting payment exhaustion test")
		t.Logf("Test Duration: %s", testDuration)
		t.Logf("Blobs per second: %f", blobsPerSecond)
		t.Logf("Payload size: %d bytes", payloadSize)

		controlLoop(t, testHarness, testDuration, blobsPerSecond, payloadSize)
	})
}

func controlLoop(
	t *testing.T,
	testHarness *integration_test.TestHarness,
	testDuration time.Duration,
	blobsPerSecond float32,
	payloadSize int,
) {
	var successCount atomic.Uint32
	var failureCount atomic.Uint32
	var blobCount atomic.Uint32

	testRandom := random.NewTestRandom()
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(t.Context(), testDuration)
	defer cancel()
	ticker := time.NewTicker(time.Duration(1000/blobsPerSecond) * time.Millisecond)
	defer ticker.Stop()
	defer func() {
		totalTime := time.Since(startTime)
		finalSuccessCount := successCount.Load()
		finalFailureCount := failureCount.Load()

		t.Logf("Test duration: %s seconds", totalTime)
		t.Logf("Total attempts: %d", blobCount.Load())
		t.Logf("Successful dispersals: %d", finalSuccessCount)
		t.Logf("Failed dispersals: %d", finalFailureCount)

		effectiveRate := float64(finalSuccessCount) / totalTime.Seconds()
		t.Logf("Effective success rate: %.2f blobs/second", effectiveRate)
		t.Logf("Total payload data dispersed: %d bytes", finalSuccessCount*uint32(payloadSize))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go func() {
				currentBlob := blobCount.Add(1)
				payload := coretypes.Payload(testRandom.Bytes(payloadSize))

				t.Logf("[%s] Dispersing blob #%d...", time.Since(startTime), currentBlob)

				_, err := testHarness.PayloadDisperser.SendPayload(t.Context(), payload)

				if err != nil {
					failureCount.Add(1)
					t.Logf("[%s] ❌ Blob #%d failed: %v", time.Since(startTime), currentBlob, err)
				} else {
					successCount.Add(1)
					t.Logf("[%s] ✅ Blob #%d succeeded", time.Since(startTime), currentBlob)
				}
			}()
		}
	}
}
