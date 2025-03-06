package payloadretrieval

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
)

// RelayPayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a RelayPayloadRetriever
type RelayPayloadRetrieverConfig struct {
	clients.PayloadClientConfig

	// The timeout duration for relay calls to retrieve blobs.
	RelayTimeout time.Duration
}

// GetDefaultRelayPayloadRetrieverConfig creates a RelayPayloadRetrieverConfig with default values
func GetDefaultRelayPayloadRetrieverConfig() *RelayPayloadRetrieverConfig {
	return &RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clients.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. It performs one of the following actions for any contained 0 values:
//
// 1. If 0 is an acceptable value for the field, do nothing.
// 2. If 0 is NOT an acceptable value for the field, and a default value is defined, then set it to the default.
// 3. If 0 is NOT an acceptable value for the field, and a default value is NOT defined, return an error.
func (rc *RelayPayloadRetrieverConfig) checkAndSetDefaults() error {
	defaultConfig := GetDefaultRelayPayloadRetrieverConfig()
	if rc.RelayTimeout == 0 {
		rc.RelayTimeout = defaultConfig.RelayTimeout
	}

	return nil
}
