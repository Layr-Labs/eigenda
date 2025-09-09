package arbitrum_altda

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EigenDAV2MessageHeaderByte is the unique EigenDAV2MessageHeaderByte.
// The value chosen is purely arbitrary.
//
// TODO: See if there can be social consensus on the value we assume,
// otherwise could eventually conflict with competitors or even OCL.
// maybe we could reuse the existing OP social contract for DA layers
// and assume 0x01?
const EigenDAV2MessageHeaderByte byte = 0x42

type PreimageType uint8

// The ALT DA server only cares about type 3 Custom DA preimage types
const (
	CustomDAPreimageType PreimageType = 3
)

// PreimagesMap maintains a nested mapping:
// //   preimage_type -> preimage_hash_key -> preimage bytes
//
// only the CustomDAPreimageType is used for EigenDAV2 batches
//
// TODO: Figure out what's getting passed over that wire. We're
// being passed a mapping and populating it with EigenDAV2 Arbitrum batch
// context. It'd be good to know when an EigenDA batch is populated
// during the preimage mapping population given we use 16mib-buffer_to_ensure_padding batch
// sizing which arbitrum doesn't
// natively support. If other preimage entries exist for the node prestate trie then data size limits
// could unknowningly be hit/exceeded.
type PreimagesMap map[PreimageType]map[common.Hash][]byte

/*
	These response types are copied verbatim (types, comments) from the upstream nitro reference implementation.
	Importing them into the EigenDA monorepo directly would overload dependency graph and create massive mgmt burden,
	requring delicate inter-play of different go-ethereum forks (especially since we already import from OP Stack).
*/

// IsValidHeaderByteResult is the result struct that data availability providers should use to
// respond if the given headerByte corresponds to their DA service
type IsValidHeaderByteResult struct {
	IsValid bool `json:"is-valid,omitempty"`
}

// RecoverPayloadFromBatchResult is the result struct that data availability providers should use
// to respond with underlying payload and updated preimages map to a RecoverPayloadFromBatch
// fetch request
type RecoverPayloadFromBatchResult struct {
	Payload hexutil.Bytes `json:"payload,omitempty"`
	// TODO: Understand the "preimage population lifecycle" to assess for potential max size risks
	Preimages PreimagesMap `json:"preimages,omitempty"`
}

// StoreResult is the result struct that data availability providers should use to respond with
// a commitment to a Store request for posting batch data to their DA service
type StoreResult struct {
	// TODO: Encoding schema
	SerializedDACert hexutil.Bytes `json:"serialized-da-cert,omitempty"`
}

// GenerateProofResult is the result struct that data availability providers should use to
// respond with a proof for a specific preimage
type GenerateProofResult struct {
	// TODO: encoding schema
	Proof hexutil.Bytes `json:"proof,omitempty"`
}

// GenerateCertificateValidityProofResult is the result struct that data availability
// providers should use to respond with validity proof
type GenerateCertificateValidityProofResult struct {
	// TODO: Figure out how to best NOOP this
	Proof hexutil.Bytes `json:"proof,omitempty"`
}
