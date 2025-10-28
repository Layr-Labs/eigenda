package reservationvalidation

import (
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/ratelimit"
)

// Contains configuration for the reservation ledger cache
type ReservationLedgerCacheConfig struct {
	// The maximum number of ReservationLedger entries to be kept in the LRU cache. This may be automatically increased
	// at runtime if premature ledger evictions are detected by the underlying cache.
	MaxLedgers int
	// Duration used to calculate bucket capacity when creating new reservation ledgers
	BucketCapacityPeriod time.Duration
	// How to handle requests that would overfill the bucket
	OverfillBehavior ratelimit.OverfillBehavior
	// Interval for checking for payment updates
	UpdateInterval time.Duration
}

// Verify validates the ReservationLedgerCacheConfig
func (c *ReservationLedgerCacheConfig) Verify() error {
	if c.MaxLedgers <= 0 {
		return errors.New("max ledgers must be > 0")
	}

	if c.MaxLedgers > maxReservationLRUCacheSize {
		return errors.New("max ledgers exceeds maximum allowed cache size")
	}

	if c.BucketCapacityPeriod <= 0 {
		return errors.New("bucket capacity period must be > 0")
	}

	if c.UpdateInterval <= 0 {
		return errors.New("update interval must be > 0")
	}

	if c.OverfillBehavior != ratelimit.OverfillNotPermitted && c.OverfillBehavior != ratelimit.OverfillOncePermitted {
		return errors.New("invalid overfill behavior")
	}

	return nil
}

// Creates a new config with validation
func NewReservationLedgerCacheConfig(
	maxLedgers int,
	bucketCapacityPeriod time.Duration,
	overfillBehavior ratelimit.OverfillBehavior,
	updateInterval time.Duration,
) (ReservationLedgerCacheConfig, error) {
	config := ReservationLedgerCacheConfig{
		MaxLedgers:           maxLedgers,
		BucketCapacityPeriod: bucketCapacityPeriod,
		OverfillBehavior:     overfillBehavior,
		UpdateInterval:       updateInterval,
	}

	if err := config.Verify(); err != nil {
		return ReservationLedgerCacheConfig{}, fmt.Errorf("failed to verify reservation ledger cache config: %w", err)
	}

	return config, nil
}
