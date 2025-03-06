package payloaddispersal

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
)

// PayloadDisperserConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadDisperser
type PayloadDisperserConfig struct {
	clients.PayloadClientConfig

	// DisperseBlobTimeout is the duration after which the PayloadDisperser will time out, when trying to disperse a
	// blob
	DisperseBlobTimeout time.Duration

	// BlobCertifiedTimeout is the duration after which the PayloadDisperser will time out, while polling
	// the disperser for blob status, waiting for BlobStatus_CERTIFIED
	BlobCertifiedTimeout time.Duration

	// BlobStatusPollInterval is the tick rate for the PayloadDisperser to use, while polling the disperser with
	// GetBlobStatus.
	BlobStatusPollInterval time.Duration

	// The timeout duration for contract calls
	ContractCallTimeout time.Duration
}

// GetDefaultPayloadDisperserConfig creates a PayloadDisperserConfig with default values
func GetDefaultPayloadDisperserConfig() *PayloadDisperserConfig {
	return &PayloadDisperserConfig{
		PayloadClientConfig:    *clients.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCertifiedTimeout:   2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (dc *PayloadDisperserConfig) checkAndSetDefaults() error {
	defaultConfig := GetDefaultPayloadDisperserConfig()

	if dc.DisperseBlobTimeout == 0 {
		dc.DisperseBlobTimeout = defaultConfig.DisperseBlobTimeout
	}

	if dc.BlobCertifiedTimeout == 0 {
		dc.BlobCertifiedTimeout = defaultConfig.BlobCertifiedTimeout
	}

	if dc.BlobStatusPollInterval == 0 {
		dc.BlobStatusPollInterval = defaultConfig.BlobStatusPollInterval
	}

	if dc.ContractCallTimeout == 0 {
		dc.ContractCallTimeout = defaultConfig.ContractCallTimeout
	}

	return nil
}
