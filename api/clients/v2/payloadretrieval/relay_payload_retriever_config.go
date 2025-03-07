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

// getDefaultRelayPayloadRetrieverConfig creates a RelayPayloadRetrieverConfig with default values
func getDefaultRelayPayloadRetrieverConfig() *RelayPayloadRetrieverConfig {
	return &RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clients.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. If a given field is 0, and 0 is not an acceptable value, then
// this method sets it to the default.
func (rc *RelayPayloadRetrieverConfig) checkAndSetDefaults() error {
	defaultConfig := getDefaultRelayPayloadRetrieverConfig()
	if rc.RelayTimeout == 0 {
		rc.RelayTimeout = defaultConfig.RelayTimeout
	}

	return nil
}
