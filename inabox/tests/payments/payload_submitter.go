package payments

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

// Submits payloads at a certain rate for a duration. Asserts the actual success rate is within tolerance of expected.
func mustSubmitPayloads(
	t *testing.T,
	testRandom *random.TestRandom,
	payloadDisperser *dispersal.PayloadDisperser,
	blobsPerSecond float32,
	payloadSize int,
	testDuration time.Duration,
	expectedSuccessRate float32,
	tolerance float32,
) {
	ctx, cancel := context.WithTimeout(t.Context(), testDuration)
	defer cancel()
	startTime := time.Now()

	secondsPerBlob := time.Duration(1.0 / blobsPerSecond * float32(time.Second))
	ticker := time.NewTicker(secondsPerBlob)
	defer ticker.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	var successCount atomic.Uint32
	var failureCount atomic.Uint32
	var blobCount atomic.Uint32

	defer func() {
		successes := successCount.Load()
		failures := failureCount.Load()
		total := successes + failures

		t.Logf("Test duration: %s", time.Since(startTime))
		t.Logf("Total attempts: %d", total)
		t.Logf("Successful dispersals: %d", successes)
		t.Logf("Failed dispersals: %d", failures)

		require.Greater(t, total, uint32(0), "no dispersals attempted")

		actualSuccessRate := float32(successes) / float32(total)

		t.Logf("Actual success rate: %.2f%%", actualSuccessRate*100)
		t.Logf("Expected success rate: %.2f%% ± %.2f%%", expectedSuccessRate*100, tolerance*100)

		minAcceptableRate := expectedSuccessRate - tolerance
		maxAcceptableRate := expectedSuccessRate + tolerance

		require.GreaterOrEqual(t, actualSuccessRate, minAcceptableRate,
			"Success rate %.2f%% below minimum %.2f%%", actualSuccessRate*100, minAcceptableRate*100)
		require.LessOrEqual(t, actualSuccessRate, maxAcceptableRate,
			"Success rate %.2f%% above maximum %.2f%%", actualSuccessRate*100, maxAcceptableRate*100)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				currentBlob := blobCount.Add(1)
				payload := coretypes.Payload(testRandom.Bytes(payloadSize))
				timestamp := time.Since(startTime)

				t.Logf("[%s] Dispersing blob #%d...", timestamp, currentBlob)

				_, err := payloadDisperser.SendPayload(t.Context(), payload)

				if err != nil {
					failureCount.Add(1)
					t.Logf("[%s] ❌ Blob #%d failed: %v", timestamp, currentBlob, err)
				} else {
					successCount.Add(1)
					t.Logf("[%s] ✅ Blob #%d succeeded", timestamp, currentBlob)
				}
			}()
		}
	}
}
