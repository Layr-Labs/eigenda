package common

import (
	"github.com/Layr-Labs/eigenda/api/clients"
)

// ClientConfigV1 wraps all the configuration values necessary to configure v1 eigenDA clients
//
// This struct wraps around an instance of clients.EigenDAClientConfig, and adds additional required values. Ideally,
// the extra values would just be part of clients.EigenDAClientConfig. Since these additions would require core changes,
// though, and v1 is slated for deprecation, this wrapper is just a stopgap to better organize configs in the proxy
// repo in the short term.
type ClientConfigV1 struct {
	EdaClientCfg     clients.EigenDAClientConfig
	MaxBlobSizeBytes uint64
	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	PutTries int
}
