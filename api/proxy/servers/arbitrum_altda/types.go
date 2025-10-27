package arbitrum_altda

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type PreimageType uint8

// The ALT DA server only cares about type 3 Custom DA preimage types
const (
	CustomDAPreimageType PreimageType = 3
)

// TODO: Reduce this mapping logic to be less generalized to
//
//	multi PreimageType since EigenDA x CustomDA only
//	cares about the one key
//
// PreimagesMap maintains a nested mapping:
// //   preimage_type -> preimage_hash_key -> preimage bytes
//
// only the CustomDAPreimageType is used for EigenDAV2 batches
type PreimagesMap map[PreimageType]map[common.Hash][]byte

// PreimageRecorder is used to add (key,value) pair to the map accessed by key = ty of a bigger map, preimages.
// If ty doesn't exist as a key in the preimages map,
// then it is intialized to map[common.Hash][]byte and then (key,value) pair is added
type PreimageRecorder func(key common.Hash, value []byte, ty PreimageType)

// RecordPreimagesTo takes in preimages map and returns a function that can be used
// In recording (hash,preimage) key value pairs into preimages map,
// when fetching payload through RecoverPayloadFromBatch
func RecordPreimagesTo(preimages PreimagesMap) PreimageRecorder {
	if preimages == nil {
		return nil
	}
	return func(key common.Hash, value []byte, ty PreimageType) {
		if preimages[ty] == nil {
			preimages[ty] = make(map[common.Hash][]byte)
		}
		preimages[ty][key] = value
	}
}

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
