package payments

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/test/random"
)

// Submits payloads at a certain rate for a duration.
//
// Returns dispersal errors on a channel.
func submitPayloads(
	t *testing.T,
	testRandom *random.TestRandom,
	payloadDisperser *payloaddispersal.PayloadDisperser,
	blobsPerSecond float32,
	payloadSize int,
	testDuration time.Duration,
) <-chan error {
	resultsChan := make(chan error)

	go func() {
		ctx, cancel := context.WithTimeout(t.Context(), testDuration)
		defer cancel()
		startTime := time.Now()

		secondsPerBlob := time.Duration(1.0 / blobsPerSecond * float32(time.Second))
		ticker := time.NewTicker(secondsPerBlob)
		defer ticker.Stop()

		var wg sync.WaitGroup
		defer func() {
			wg.Wait()
			close(resultsChan)
		}()

		var successCount atomic.Uint32
		var failureCount atomic.Uint32
		var blobCount atomic.Uint32

		defer func() {
			successes := successCount.Load()
			failures := failureCount.Load()

			t.Logf("Test duration: %s", time.Since(startTime))
			t.Logf("Total attempts: %d", successes+failures)
			t.Logf("Successful dispersals: %d", successes)
			t.Logf("Failed dispersals: %d", failures)
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
					resultsChan <- err

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
	}()

	return resultsChan
}
