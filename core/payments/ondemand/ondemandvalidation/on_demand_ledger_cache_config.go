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

// Verify validates the OnDemandLedgerCacheConfig
func (c *OnDemandLedgerCacheConfig) Verify() error {
	if c.MaxLedgers <= 0 {
		return errors.New("max ledgers must be > 0")
	}

	if c.OnDemandTableName == "" {
		return errors.New("on-demand table name must not be empty")
	}

	if c.UpdateInterval <= 0 {
		return errors.New("update interval must be > 0")
	}

	return nil
}

// Creates a new config with validation
func NewOnDemandLedgerCacheConfig(
	maxLedgers int,
	onDemandTableName string,
	updateInterval time.Duration,
) (OnDemandLedgerCacheConfig, error) {
	config := OnDemandLedgerCacheConfig{
		MaxLedgers:        maxLedgers,
		OnDemandTableName: onDemandTableName,
		UpdateInterval:    updateInterval,
	}

	if err := config.Verify(); err != nil {
		return OnDemandLedgerCacheConfig{}, err
	}

	return config, nil
}
