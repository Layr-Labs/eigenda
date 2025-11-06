package common

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
)

// CompatibilityConfig ... CompatibilityConfig stores values useful to external services for checking compatibility
// with the proxy instance, such as version, chainID, and recency window size. These values are returned by the rest
// servers /config endpoint.
type CompatibilityConfig struct {
	// Current proxy version in the format {MAJOR}.{MINOR}.{PATCH}-{META} e.g: 2.4.0-43-g3b4f9f40. The version
	// is injected at build using `git describe --tags --always --dirty`. This allows a service to perform a
	// minimum version supported check.
	Version string `json:"version"`
	// The ChainID of the connected ethClient. This allows a service to check which chain the proxy is connected
	// to. If the proxy has memstore enabled, a ChainID of "" will be set.
	ChainID string `json:"chain_id"`
	// The EigenDA directory address. This allows a service to verify which contracts are being used by the proxy.
	DirectoryAddress string `json:"directory_address"`
	// The cert verifier router or immutable contract address. This allows a service to verify the cert verifier being
	// used by the proxy.
	CertVerifierAddress string `json:"cert_verifier_address"`
	// The max supported payload size in bytes supported by the proxy instance. Calculated from `MaxBlobSizeBytes`.
	MaxPayloadSizeBytes uint32 `json:"max_payload_size_bytes"`
	// The recency window size. This allows a service (e.g batch poster) to check alignment with the proxy instance.
	RecencyWindowSize uint32 `json:"recency_window_size"`
	// The APIs currently enabled on the rest server
	APIsEnabled []string `json:"apis_enabled,omitempty"`
	// Whether the proxy is in read-only mode (no signer payment key)
	ReadOnlyMode bool `json:"read_only_mode"`
}

func NewCompatibilityConfig(
	version string,
	chainID string,
	clientConfigV2 ClientConfigV2,
	readOnly bool,
	APIsEnabled []string,
) (CompatibilityConfig, error) {
	var maxPayloadSize uint32 = 0
	// If the proxy is in v1 mode (soon to be removed) a v2 MaxBlobSizeBytes is not set.
	if clientConfigV2.MaxBlobSizeBytes > 0 {
		var err error
		// BlobSymbolsToMaxPayloadSize returns an err if the given blob length symbols is 0
		maxPayloadSize, err = codec.BlobSymbolsToMaxPayloadSize(
			uint32(clientConfigV2.MaxBlobSizeBytes / encoding.BYTES_PER_SYMBOL))
		if err != nil {
			return CompatibilityConfig{}, fmt.Errorf("calculate max payload size: %w", err)
		}
	}

	// Remove 'v' prefix from version string if present for compatibility with eigenda/common/version helper funcs
	if len(version) > 0 {
		versionRunes := []rune(version)
		if versionRunes[0] == 'v' || versionRunes[0] == 'V' {
			version = string(versionRunes[1:])
		}
	}

	return CompatibilityConfig{
		Version:             version,
		ChainID:             chainID,
		DirectoryAddress:    clientConfigV2.EigenDADirectory,
		CertVerifierAddress: clientConfigV2.EigenDACertVerifierOrRouterAddress,
		MaxPayloadSizeBytes: maxPayloadSize,
		RecencyWindowSize:   clientConfigV2.RBNRecencyWindowSize,
		APIsEnabled:         APIsEnabled,
		ReadOnlyMode:        readOnly,
	}, nil
}
