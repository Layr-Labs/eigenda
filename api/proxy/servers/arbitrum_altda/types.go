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
type PreimagesMap map[PreimageType]map[common.Hash][]byte

/*
	These response types are copied verbatim (types, comments) from the upstream nitro reference implementation.
	Importing them into the EigenDA monorepo directly would overload dependency graph and create massive mgmt burden,
	requring delicate inter-play of different go-ethereum forks (especially since we already import from OP Stack).
*/

// PreimagesResult contains the collected preimages
type PreimagesResult struct {
	Preimages PreimagesMap
}

// PayloadResult contains the recovered payload data
type PayloadResult struct {
	Payload []byte
}

// SupportedHeaderBytesResult is the result struct that data availability providers should use to respond with
// their supported header bytes
type SupportedHeaderBytesResult struct {
	HeaderBytes hexutil.Bytes `json:"headerBytes,omitempty"`
}

// StoreResult is the result struct that data availability providers should use to respond with a commitment to a
// Store request for posting batch data to their DA service
type StoreResult struct {
	SerializedDACert hexutil.Bytes `json:"serialized-da-cert,omitempty"`
}

// GenerateReadPreimageProofResult is the result struct that data availability providers
// should use to respond with a proof for a specific preimage
type GenerateReadPreimageProofResult struct {
	Proof hexutil.Bytes `json:"proof,omitempty"`
}

// GenerateCertificateValidityProofResult is the result struct that data availability providers should use to
// respond with validity proof
type GenerateCertificateValidityProofResult struct {
	Proof hexutil.Bytes `json:"proof,omitempty"`
}
