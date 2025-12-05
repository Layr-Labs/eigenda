package signingrate

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// This function periodically flushes signing rate data from a SigningRateTracker to
// persistent storage using a SigningRateStorage. This function spins until the context is cancelled,
// and so it should be run  on a background goroutine.
func SigningRateStorageFlusher(
	ctx context.Context,
	logger logging.Logger,
	tracker SigningRateTracker,
	storage SigningRateStorage,
	flushPeriod time.Duration,
) {

	ticker := time.NewTicker(flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			buckets, err := tracker.GetUnflushedBuckets()
			if err != nil {
				logger.Errorf("Error getting unflushed buckets: %v", err)
				continue
			}

			if len(buckets) == 0 {
				// nothing to flush
				continue
			}

			err = storage.StoreBuckets(ctx, buckets)
			if err != nil {
				logger.Errorf("Error storing signing rate buckets: %v", err)
				continue
			}
		}
	}
}
