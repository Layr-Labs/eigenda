package ondemandvalidation

import (
	"errors"
	"time"
)

// Contains configuration for the on-demand ledger cache
type OnDemandLedgerCacheConfig struct {
	// The maximum number of OnDemandLedger entries to be kept in the LRU cache
	MaxLedgers int
	// The name of the dynamo table where on-demand payment information is stored
	OnDemandTableName string
	// Interval for checking for payment updates
	UpdateInterval time.Duration
}

// Creates a new config with validation
func NewOnDemandLedgerCacheConfig(
	maxLedgers int,
	onDemandTableName string,
	updateInterval time.Duration,
) (OnDemandLedgerCacheConfig, error) {
	if maxLedgers <= 0 {
		return OnDemandLedgerCacheConfig{}, errors.New("max ledgers must be > 0")
	}

	if onDemandTableName == "" {
		return OnDemandLedgerCacheConfig{}, errors.New("on-demand table name must not be empty")
	}

	if updateInterval <= 0 {
		return OnDemandLedgerCacheConfig{}, errors.New("update interval must be > 0")
	}

	return OnDemandLedgerCacheConfig{
		MaxLedgers:        maxLedgers,
		OnDemandTableName: onDemandTableName,
		UpdateInterval:    updateInterval,
	}, nil
}
