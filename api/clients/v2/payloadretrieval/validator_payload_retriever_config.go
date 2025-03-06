package payloadretrieval

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
)

// ValidatorPayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values
// needed by a ValidatorPayloadRetriever
type ValidatorPayloadRetrieverConfig struct {
	clients.PayloadClientConfig

	// The timeout duration for retrieving chunks from a given quorum, and reassembling the chunks into a blob.
	// Once this timeout triggers, the retriever will give up on the quorum, and retry with the next quorum (if one exists)
	RetrievalTimeout time.Duration
}

// GetDefaultValidatorPayloadRetrieverConfig creates a ValidatorPayloadRetrieverConfig with default values
func GetDefaultValidatorPayloadRetrieverConfig() *ValidatorPayloadRetrieverConfig {
	return &ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clients.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    30 * time.Second,
	}
}

// checkAndSetDefaults checks an existing config struct. If a given field is 0, and 0 is not an acceptable value, then
// this method sets it to the default.
func (rc *ValidatorPayloadRetrieverConfig) checkAndSetDefaults() error {
	defaultConfig := GetDefaultValidatorPayloadRetrieverConfig()
	if rc.RetrievalTimeout == 0 {
		rc.RetrievalTimeout = defaultConfig.RetrievalTimeout
	}

	return nil
}
