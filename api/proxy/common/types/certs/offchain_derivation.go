package certs

import "github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"

// OffchainDerivationParameters holds parameters for offchain derivation for a given derivation version.
// Version 0 is currently the only offchain derivation version, which only contains the RBN recency window size
// parameter. However this struct is designed to be extensible for future offchain derivation versions.
type OffchainDerivationParameters struct {
	// Allowed distance (in L1 blocks) between the eigenDA cert's reference block number (RBN)
	// and the L1 block number at which the cert was included in the rollup's batch inbox.
	// If cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize, an
	// [RBNRecencyCheckFailedError] is returned.
	// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation
	RBNRecencyWindowSize uint64
}

// OffchainDerivationMap maps offchain derivation versions to their parameters.
type OffchainDerivationMap = map[coretypes.OffchainDerivationVersion]OffchainDerivationParameters
