package signingrate

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

func LoadSigningRateDataFromStorage(
	ctx context.Context,
	logger logging.Logger,
	tracker SigningRateTracker,
	storage SigningRateStorage,
	signingRateRetentionPeriod time.Duration,
) error {

	startTimestamp := time.Now().Add(-signingRateRetentionPeriod)

	buckets, err := storage.LoadBuckets(ctx, startTimestamp)
	if err != nil {
		return fmt.Errorf("loading signing rate buckets from storage: %w", err)
	}

	logger.Debugf("Loaded signing rate data from storage starting at %v, found %d buckets",
		startTimestamp, len(buckets))

	for _, bucket := range buckets {
		tracker.UpdateLastBucket(bucket)
	}

	return nil
}
