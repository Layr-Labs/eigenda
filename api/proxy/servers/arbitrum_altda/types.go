package arbitrum_altda

import (
	proxy_common "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	/*
		MessageHeader is a 40 byte prefix encoding added to the SequencerMessage
		that is constructed during a batch poster tx to the SequencerInbox which
		appends a new SequencerMessage (e.g, DA Cert) to the safe/final rollup
		tx feed.

		MessageHeader is re-derived as part of the nitro derivation pipeline and is trustlessly
		enforced since keccak256(header + DA Cert) is committed by the SequencerInbox message accumulator
		which is used to referee one step proofs for READINBOXMESSAGE opcode disputes.

		the first 4 fields of the header are a "time boundary" that's
		computed based on the inbox tx block # and rollup provided
		"time variation" values where:

		  minTimeStamp = block.timestamp - delaySeconds
		  minBlockNumber = block.number - delayBlocks

		  maxTimeStamp = block.timestamp + futureSeconds
		  maxBlockNumber = block.number + futureBlocks


		1. MinTimestamp (bytes 0-7) - Minimum timestamp for the batch
		2. MaxTimestamp (bytes 8-15) - Maximum timestamp for the batch
		3. MinL1Block (bytes 16-23) - Minimum L1 block number
		4. MaxL1Block (bytes 24-31) - Maximum L1 block number
		5. AfterDelayedMessages (bytes 32-39) - Number of delayed messages processed
	*/

	// Offset used to determine the MessageHeader
	MessageHeaderOffset = 40

	// Number of DA Commitment encoding bytes prefixed to the DA Cert bytes
	// by the ArbitrumCommitment encoding
	DACommitPrefixBytes = 2

	// Offset used to determine where in the Sequencer Message that
	// the first DA Cert byte starts
	DACertOffset = MessageHeaderOffset + DACommitPrefixBytes
)

const (
	// trusted integration
	MethodGetSupportedHeaderBytes = "daprovider_getSupportedHeaderBytes"
	MethodStore                   = "daprovider_store"
	MethodRecoverPayload          = "daprovider_recoverPayload"
	MethodCollectPreimages        = "daprovider_collectPreimages"
	// trustless integration
	MethodGenerateReadPreimageProof = "daprovider_generateReadPreimageProof"
	MethodGenerateCertValidityProof = "daprovider_generateCertificateValidityProof"
	// compatibility check
	MethodCompatibilityConfig = "daprovider_compatibilityConfig"
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
	HeaderBytes []hexutil.Bytes `json:"headerBytes,omitempty"`
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

// CompatibilityConfigResult is the result struct used to check compatibility between the proxy instance and an
// external service
type CompatibilityConfigResult struct {
	proxy_common.CompatibilityConfig
}
