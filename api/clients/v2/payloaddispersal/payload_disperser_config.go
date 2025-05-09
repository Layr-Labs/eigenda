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

	// BlobCompleteTimeout is the duration after which the PayloadDisperser will time out, while polling
	// the disperser for blob status, waiting for BlobStatus_COMPLETE
	BlobCompleteTimeout time.Duration

	// BlobStatusPollInterval is the tick rate for the PayloadDisperser to use, while polling the disperser with
	// GetBlobStatus.
	BlobStatusPollInterval time.Duration

	// The timeout duration for contract calls
	ContractCallTimeout time.Duration

	// // whether to EigenDACertVerifierRouter or historical immutable EigenDACertVerifierV2
	// UseRouter bool
}

// getDefaultPayloadDisperserConfig creates a PayloadDisperserConfig with default values
func getDefaultPayloadDisperserConfig() *PayloadDisperserConfig {
	return &PayloadDisperserConfig{
		PayloadClientConfig:    *clients.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCompleteTimeout:    2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. If a given field is 0, and 0 is not an acceptable value, then
// this method sets it to the default.
func (dc *PayloadDisperserConfig) checkAndSetDefaults() error {
	defaultConfig := getDefaultPayloadDisperserConfig()

	if dc.DisperseBlobTimeout == 0 {
		dc.DisperseBlobTimeout = defaultConfig.DisperseBlobTimeout
	}

	if dc.BlobCompleteTimeout == 0 {
		dc.BlobCompleteTimeout = defaultConfig.BlobCompleteTimeout
	}

	if dc.BlobStatusPollInterval == 0 {
		dc.BlobStatusPollInterval = defaultConfig.BlobStatusPollInterval
	}

	if dc.ContractCallTimeout == 0 {
		dc.ContractCallTimeout = defaultConfig.ContractCallTimeout
	}

	return nil
}
