package reservation

import (
	"errors"
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

// Creates a new config with validation
func NewReservationLedgerCacheConfig(
	maxLedgers int,
	bucketCapacityPeriod time.Duration,
	overfillBehavior ratelimit.OverfillBehavior,
	updateInterval time.Duration,
) (ReservationLedgerCacheConfig, error) {
	if maxLedgers <= 0 {
		return ReservationLedgerCacheConfig{}, errors.New("max ledgers must be > 0")
	}

	if maxLedgers > maxReservationLRUCacheSize {
		return ReservationLedgerCacheConfig{}, errors.New("max ledgers exceeds maximum allowed cache size")
	}

	if bucketCapacityPeriod <= 0 {
		return ReservationLedgerCacheConfig{}, errors.New("bucket capacity period must be > 0")
	}

	if updateInterval <= 0 {
		return ReservationLedgerCacheConfig{}, errors.New("update interval must be > 0")
	}

	if overfillBehavior != ratelimit.OverfillNotPermitted && overfillBehavior != ratelimit.OverfillOncePermitted {
		return ReservationLedgerCacheConfig{}, errors.New("invalid overfill behavior")
	}

	return ReservationLedgerCacheConfig{
		MaxLedgers:           maxLedgers,
		BucketCapacityPeriod: bucketCapacityPeriod,
		OverfillBehavior:     overfillBehavior,
		UpdateInterval:       updateInterval,
	}, nil
}
