package ejector

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/config/secret"
)

var _ config.DocumentedConfig = (*EjectorConfig)(nil)

// Configuration for the ejector.
type EjectorConfig struct {

	// The Ethereum RPC URL(s) to use for connecting to the blockchain.
	EthRpcUrls []*secret.Secret `docs:"required"`

	// The private key to use for signing ejection transactions, in hex.
	// Do not include the '0x' prefix. This is required if KMS is not configured.
	PrivateKey *secret.Secret `docs:"required"`

	// The address of the contract directory contract.
	ContractDirectoryAddress string `docs:"required"`

	// The URL of the Eigenda Data API to use for looking up signing rates.
	DataApiUrl string `docs:"required"`

	// The AWS KMS Key ID to use for signing transactions. Only required if the private key is not provided via
	// the Secret.PrivateKey field.
	KmsKeyId string `docs:"required"`

	// The AWS region where the KMS key is located. Only required if KmsKeyId is provided.
	KmsRegion string `docs:"required"`

	// The AWS KMS endpoint to use. Only required if using a custom endpoint (e.g., LocalStack).
	KmsEndpoint string

	// The number of times to retry a failed Ethereum RPC call.
	EthRpcRetryCount int

	// The number of block confirmations to wait for before considering an ejection transaction to be confirmed.
	EthBlockConfirmations int

	// The timeout to use when making requests to the Data API.
	DataApiTimeout time.Duration

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
	EjectionThrottle float64

	// The time period over which the ejection rate limit is calculated. The ejection manager will be allowed to eject
	// ejectionRateLimit fraction of stake every EjectionThrottleTimePeriod.
	EjectionThrottleTimePeriod time.Duration

	// If true, then the ejection manager will immediately be able to eject ejectionRateLimit fraction of stake when it
	// starts up. If false, then the ejection manager will need to wait before it has this capacity.
	StartEjectionThrottleFull bool

	// A list of validator addresses that we should never attempt to eject, even if they otherwise
	// meet the ejection criteria.
	DoNotEjectTheseValidators []string

	// The period at which to periodically attempt to finalize ejections that have been started.
	EjectionFinalizationPeriod time.Duration

	// The number of blocks to wait before using a reference block number. That is to say, do not always
	// read data from the latest block  we know about, but rather read from a block that is sufficiently old as to make
	// choosing the wrong fork unlikely.
	ReferenceBlockNumberOffset uint64

	// The interval at which to poll for a new reference block number.
	ReferenceBlockNumberPollInterval time.Duration

	// The size for the caches for on-chain data.
	ChainDataCacheSize uint64

	// The output type for logs, must be "json" or "text".
	LogOutputType string

	// Whether to enable color in log output (only applies to text output).
	LogColor bool

	// If non-zero, this value will be used as the gas limit for transactions, overriding the gas estimation.
	MaxGasOverride uint64

	// Flip this flag to true if you want to disable ejections. Useful for emergency situations where you want
	// to stop the ejector from ejecting validators, but without tearing down the kube infrastructure.
	DisableEjections bool

	// The period between verbose signing rate data dumps. If zero, then verbose signing rate logging is disabled.
	SigningRateLogPeriod time.Duration
}

// DefaultEjectorConfig returns a default configuration for the ejector.
func DefaultEjectorConfig() *EjectorConfig {
	return &EjectorConfig{
		EjectionPeriod:                       time.Minute,
		EjectionCriteriaTimeWindow:           10 * time.Minute,
		EjectionFinalizationDelay:            time.Hour,
		EjectionRetryDelay:                   24 * time.Hour,
		MaxConsecutiveFailedEjectionAttempts: 5,
		EjectionThrottle:                     0.05, // 5% of stake can be ejected every EjectionThrottleTimePeriod
		EjectionThrottleTimePeriod:           24 * time.Hour,
		StartEjectionThrottleFull:            false,
		EjectionFinalizationPeriod:           time.Minute,
		DataApiTimeout:                       60 * time.Second,
		EthRpcRetryCount:                     3,
		EthBlockConfirmations:                0,
		ReferenceBlockNumberOffset:           64,
		ReferenceBlockNumberPollInterval:     10 * time.Second,
		ChainDataCacheSize:                   1024,
		LogOutputType:                        string(common.JSONLogFormat),
		LogColor:                             false,
		MaxGasOverride:                       10_000_000,
		DisableEjections:                     false,
		SigningRateLogPeriod:                 time.Hour,
	}
}

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

	if c.EjectionThrottle <= 0 || c.EjectionThrottle > 1.0 {
		return fmt.Errorf("invalid ejection rate limit: %f", c.EjectionThrottle)
	}

	if c.EjectionThrottleTimePeriod <= 0 {
		return fmt.Errorf("invalid ejection throttle time period: %s", c.EjectionThrottleTimePeriod)
	}

	if c.DataApiUrl == "" {
		return fmt.Errorf("invalid data API URL: %s", c.DataApiUrl)
	}

	if c.DataApiTimeout <= 0 {
		return fmt.Errorf("invalid data API timeout: %s", c.DataApiTimeout)
	}

	if c.EjectionFinalizationPeriod <= 0 {
		return fmt.Errorf("invalid ejection finalization period: %s", c.EjectionFinalizationPeriod)
	}

	if c.ReferenceBlockNumberPollInterval <= 0 {
		return fmt.Errorf("invalid reference block number poll interval: %s", c.ReferenceBlockNumberPollInterval)
	}

	if c.ChainDataCacheSize <= 0 {
		return fmt.Errorf("invalid chain data cache size: %d", c.ChainDataCacheSize)
	}
	if c.SigningRateLogPeriod < 0 {
		return fmt.Errorf("invalid signing rate log period: %s", c.SigningRateLogPeriod)
	}
	if len(c.EthRpcUrls) == 0 {
		return fmt.Errorf("at least one Ethereum RPC URL must be provided")
	}
	for _, url := range c.EthRpcUrls {
		if url.Get() == "" {
			return fmt.Errorf("EthRpcUrls cannot be empty strings")
		}
	}

	// Either a private key must be provided or KMS must be configured.
	if c.PrivateKey.Get() == "" {
		if c.KmsKeyId == "" {
			return fmt.Errorf("either a private key or KMS Key ID must be provided")
		}
		if c.KmsRegion == "" {
			return fmt.Errorf("KMS region must be provided when KMS Key ID is set")
		}
	}

	return nil
}

func (c *EjectorConfig) GetEnvVarPrefix() string {
	return "EJECTOR"
}

func (c *EjectorConfig) GetName() string {
	return "Ejector"
}

func (c *EjectorConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/ejector",
		"github.com/Layr-Labs/eigenda/common/config/secret",
	}
}