package ejector

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

// The environment variable prefix to use for the ejector configuration.
const EjectorConfigEnvPrefix = "EJECTOR"

var _ config.VerifiableConfig = (*EjectorConfig)(nil)

// Configuration for the ejector.
type EjectorConfig struct {

	////////////////////////
	// Required arguments //
	////////////////////////

	// The address of the contract directory contract.
	ContractDirectoryAddress string

	// The Ethereum RPC URL to use for connecting to the blockchain.
	EthRPCURL string

	// The private key to use for signing ejection transactions.
	PrivateKey string

	// The URL of the Eigenda Data API to use for looking up signing rates.
	DataApiUrl string

	////////////////////////
	// Optional arguments //
	////////////////////////

	// The period with which to evaluate validators for ejection.
	EjectionPeriod time.Duration

	// The time window over which to evaluate signing metrics when deciding whether to eject a validator.
	EjectionCriteriaTimeWindow time.Duration

	// The time between starting an ejection and when the ejection can be finalized.
	EjectionFinalizationDelay time.Duration

	// The minimum time to wait before retrying a failed ejection.
	EjectionRetryDelay time.Duration

	// The maximum number of consecutive failed ejection attempts before giving up on ejecting a validator.
	MaxConsecutiveFailedEjectionAttempts uint32

	// The maximum fraction of stake (out of 1.0) that can be ejected during an ejection time period.
	EjectionRateLimit float64

	// The time period over which the ejection rate limit is calculated. The ejection manager will be allowed to eject
	// ejectionRateLimit fraction of stake every EjectionThrottleTimePeriod.
	EjectionThrottleTimePeriod time.Duration

	// If true, then the ejection manager will immediately be able to eject ejectionRateLimit fraction of stake when it
	// starts up. If false, then the ejection manager will need to wait before it has this capacity.
	StartEjectionThrottleFull bool
}

// DefaultEjectorConfig returns a default configuration for the ejector.
func DefaultEjectorConfig() *EjectorConfig {
	return &EjectorConfig{
		EjectionPeriod:                       time.Minute,
		EjectionCriteriaTimeWindow:           10 * time.Minute,
		EjectionFinalizationDelay:            time.Hour,
		EjectionRetryDelay:                   24 * time.Hour,
		MaxConsecutiveFailedEjectionAttempts: 5,
		EjectionRateLimit:                    0.05, // 5% of stake can be ejected every EjectionThrottleTimePeriod
		EjectionThrottleTimePeriod:           24 * time.Hour,
		StartEjectionThrottleFull:            false,
	}
}

// Verify checks that the configuration is valid.
func (c *EjectorConfig) Verify() error {
	if c.EjectionPeriod <= 0 {
		return fmt.Errorf("invalid ejection period: %s", c.EjectionPeriod)
	}

	if c.EjectionCriteriaTimeWindow <= 0 {
		return fmt.Errorf("invalid ejection criteria time window: %s", c.EjectionCriteriaTimeWindow)
	}

	if c.ContractDirectoryAddress == "" {
		return fmt.Errorf("invalid contract directory address: %s", c.ContractDirectoryAddress)
	}

	if c.EjectionFinalizationDelay <= 0 {
		return fmt.Errorf("invalid ejection finalization delay: %s", c.EjectionFinalizationDelay)
	}

	if c.EjectionRetryDelay <= 0 {
		return fmt.Errorf("invalid ejection retry delay: %s", c.EjectionRetryDelay)
	}

	if c.MaxConsecutiveFailedEjectionAttempts == 0 {
		return fmt.Errorf("invalid max consecutive failed ejection attempts: %d",
			c.MaxConsecutiveFailedEjectionAttempts)
	}

	if c.EjectionRateLimit <= 0 || c.EjectionRateLimit > 1.0 {
		return fmt.Errorf("invalid ejection rate limit: %f", c.EjectionRateLimit)
	}

	if c.EjectionThrottleTimePeriod <= 0 {
		return fmt.Errorf("invalid ejection throttle time period: %s", c.EjectionThrottleTimePeriod)
	}
	if c.EthRPCURL == "" {
		return fmt.Errorf("invalid Ethereum RPC URL: %s", c.EthRPCURL)
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("invalid private key")
	}
	if c.DataApiUrl == "" {
		return fmt.Errorf("invalid data API URL: %s", c.DataApiUrl)
	}

	return nil
}
