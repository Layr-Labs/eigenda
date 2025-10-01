package ejector

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

var _ config.VerifiableConfig = (*EjectorConfig)(nil)

type EjectorConfig struct {
	// The period with which to evaluate validators for ejection.
	ejectionPeriod time.Duration

	// The time window over which to evaluate signing metrics when deciding whether to eject a validator.
	ejectionCriteriaTimeWindow time.Duration

	// The address of the contract directory contract.
	contractDirectoryAddress string

	// The time between starting an ejection and when the ejection can be finalized.
	ejectionFinalizationDelay time.Duration

	// The minimum time to wait before retrying a failed ejection.
	ejectionRetryDelay time.Duration

	// The maximum number of consecutive failed ejection attempts before giving up on ejecting a validator.
	maxConsecutiveFailedEjectionAttempts uint32

	// The maximum fraction of stake (out of 1.0) that can be ejected during an ejection time period.
	ejectionRateLimit float64

	// The time period over which the ejection rate limit is calculated. The ejection manager will be allowed to eject
	// ejectionRateLimit fraction of stake every ejectionThrottleTimePeriod.
	ejectionThrottleTimePeriod time.Duration

	// If true, then the ejection manager will immediately be able to eject ejectionRateLimit fraction of stake when it
	// starts up. If false, then the ejection manager will need to wait before it has this capacity.
	startEjectionThrottleFull bool

	// The Ethereum RPC URL to use for connecting to the blockchain.
	ethRPCURL string

	// The private key to use for signing ejection transactions.
	privateKey string

	// The URL of the Eigenda Data API to use for looking up signing rates.
	dataAPIURL string
}

// DefaultEjectorConfig returns a default configuration for the ejector.
func DefaultEjectorConfig() *EjectorConfig {
	return &EjectorConfig{
		ejectionPeriod:                       time.Minute,
		ejectionCriteriaTimeWindow:           10 * time.Minute,
		ejectionFinalizationDelay:            time.Hour,
		ejectionRetryDelay:                   24 * time.Hour,
		maxConsecutiveFailedEjectionAttempts: 5,
		ejectionRateLimit:                    0.05, // 5% of stake can be ejected every ejectionThrottleTimePeriod
		ejectionThrottleTimePeriod:           24 * time.Hour,
		startEjectionThrottleFull:            false,
	}
}

// Verify checks that the configuration is valid.
func (c *EjectorConfig) Verify() error {
	if c.ejectionPeriod <= 0 {
		return fmt.Errorf("invalid ejection period: %s", c.ejectionPeriod)
	}

	if c.ejectionCriteriaTimeWindow <= 0 {
		return fmt.Errorf("invalid ejection criteria time window: %s", c.ejectionCriteriaTimeWindow)
	}

	if c.contractDirectoryAddress == "" {
		return fmt.Errorf("invalid contract directory address: %s", c.contractDirectoryAddress)
	}

	if c.ejectionFinalizationDelay <= 0 {
		return fmt.Errorf("invalid ejection finalization delay: %s", c.ejectionFinalizationDelay)
	}

	if c.ejectionRetryDelay <= 0 {
		return fmt.Errorf("invalid ejection retry delay: %s", c.ejectionRetryDelay)
	}

	if c.maxConsecutiveFailedEjectionAttempts == 0 {
		return fmt.Errorf("invalid max consecutive failed ejection attempts: %d",
			c.maxConsecutiveFailedEjectionAttempts)
	}

	if c.ejectionRateLimit <= 0 || c.ejectionRateLimit > 1.0 {
		return fmt.Errorf("invalid ejection rate limit: %f", c.ejectionRateLimit)
	}

	if c.ejectionThrottleTimePeriod <= 0 {
		return fmt.Errorf("invalid ejection throttle time period: %s", c.ejectionThrottleTimePeriod)
	}
	if c.ethRPCURL == "" {
		return fmt.Errorf("invalid Ethereum RPC URL: %s", c.ethRPCURL)
	}
	if c.privateKey == "" {
		return fmt.Errorf("invalid private key")
	}
	if c.dataAPIURL == "" {
		return fmt.Errorf("invalid data API URL: %s", c.dataAPIURL)
	}

	return nil
}
