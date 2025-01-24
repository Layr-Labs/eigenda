package clients

import "time"

// PayloadRetrieverConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadRetriever
type PayloadRetrieverConfig struct {
	PayloadClientConfig

	// The timeout duration for relay calls to retrieve blobs.
	RelayTimeout time.Duration
}
