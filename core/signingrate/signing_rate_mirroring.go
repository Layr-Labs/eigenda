package signingrate

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// A function that fetches signing rate data from some source starting from the given time.
type SigningRateScraper func(ctx context.Context, startTime time.Time) ([]*validator.SigningRateBucket, error)

// Do an initial scrape of signing rate data from a remote source and ingest it into the given tracker.
// This makes it so that external callers never view an empty tracker at startup.
func DoInitialScrape(
	ctx context.Context,
	logger logging.Logger,
	// A function that can fetch signing rate data from some source.
	scraper SigningRateScraper,
	// The signing rate tracker that will mirror the remote data.
	tracker SigningRateTracker,
	// The amount of time to mirror data for. Data older than this period will not be mirrored.
	timePeriod time.Duration,
) error {

	logger.Info("Doing initial scrape of signing rate data", "time_period", timePeriod.String())

	startTime := time.Now().Add(-timePeriod)
	buckets, err := scraper(ctx, startTime)
	if err != nil {
		return fmt.Errorf("failed to do initial scrape of signing rate data: %w", err)
	}

	for _, bucket := range buckets {
		tracker.UpdateLastBucket(bucket)
	}

	logger.Info("Completed initial scrape of signing rate data", "num_buckets", len(buckets))

	return nil
}

// Call this function to mirror signing rate data from a remote source. This method does not return and should
// be run in its own goroutine.
func MirrorSigningRate(
	ctx context.Context,
	logger logging.Logger,
	// A function that can fetch signing rate data from some source.
	scraper SigningRateScraper,
	// The signing rate tracker that will mirror the remote data.
	tracker SigningRateTracker,
	// How often to poll the remote source for new data.
	interval time.Duration,
	// The amount of time to mirror data for. Data older than this period will not be mirrored.
	timePeriod time.Duration,
) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	previousScrapeTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping signing rate mirroring")
			return
		case <-ticker.C:
			currentTime := time.Now()

			buckets, err := scraper(ctx, previousScrapeTime)
			if err != nil {
				logger.Error("Failed to scrape signing rate data", "err", err)
				continue
			}

			for _, bucket := range buckets {
				tracker.UpdateLastBucket(bucket)
			}

			previousScrapeTime = currentTime
		}
	}
}
