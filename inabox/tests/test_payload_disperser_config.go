package integration

import (
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
)

// TestPayloadDisperserConfig configures how a PayloadDisperser client should be set up for testing.
//
// This struct is intentionally sparse, containing only fields that must be specifically set during testing. If any
// additional fields need modification in tests written in the future, they should be added here. Otherwise, all
// configuration fields for constructing a PayloadDisperser should simply be hardcoded in the test construction helpers.
type TestPayloadDisperserConfig struct {
	// Payment mode the client should use
	ClientLedgerMode clientledger.ClientLedgerMode

	// Private key to use for the disperser account (hex string with or without 0x prefix).
	// If empty string, a random private key will be generated.
	PrivateKey string
}

// Returns a PayloadDisperserConfig with default values for testing.
//
// The default private key is one that has a large reservation automatically allocated when setting up the payment
// vault.
func GetDefaultTestPayloadDisperserConfig() TestPayloadDisperserConfig {
	return TestPayloadDisperserConfig{
		ClientLedgerMode: clientledger.ClientLedgerModeLegacy,
		PrivateKey:       "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded",
	}
}
