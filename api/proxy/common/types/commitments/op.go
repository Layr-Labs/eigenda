package commitments

import (
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/ethereum/go-ethereum/crypto"
)

// OPCommitmentByte is the commitment type prefix.
type OPCommitmentByte byte

// CommitmentType describes the binary format of the commitment.
// OPKeccak256CommitmentByte is the default commitment type for optimism's centralized DA storage.
// OPGenericCommitmentByte indicates an opaque bytestring that the op-node never opens.
const (
	OPKeccak256CommitmentByte OPCommitmentByte = 0
	OPGenericCommitmentByte   OPCommitmentByte = 1
)

// See https://specs.optimism.io/experimental/alt-da.html#example-commitments
const EigenDALayerByte = byte(0)

// OPKeccak256Commitment is an implementation of OPCommitment that uses Keccak256 as the commitment function.
type OPKeccak256Commitment []byte

// NewOPKeccak256Commitment creates a new commitment from the given input.
func NewOPKeccak256Commitment(input []byte) OPKeccak256Commitment {
	return OPKeccak256Commitment(crypto.Keccak256(input))
}

// Encode adds a 0x00 byte prefix in front of the keccak commitment.
// Encoding is thus [ 0x00 | keccak_commitment ]
// See https://specs.optimism.io/experimental/alt-da.html#example-commitments
func (c OPKeccak256Commitment) Encode() []byte {
	return append([]byte{byte(OPKeccak256CommitmentByte)}, c...)
}

// OPEigenDAGenericCommitment is an implementation of OPCommitment that treats the commitment as an opaque bytestring.
type OPEigenDAGenericCommitment struct {
	versionedCert certs.VersionedCert
}

// NewOPEigenDAGenericCommitment creates a new commitment from the given input.
func NewOPEigenDAGenericCommitment(versionedCert certs.VersionedCert) OPEigenDAGenericCommitment {
	return OPEigenDAGenericCommitment{versionedCert}
}

// Encode adds a 2 byte header in front of the serialized versioned cert,
// to turn it into an altda commitment. See https://specs.optimism.io/experimental/alt-da.html#example-commitments
// Encoding is thus [ commitment_type_byte | da_layer_byte | eigenda_commitment ]
// which for eigenda is [ 0x01 | 0x00 | serialized_versioned_cert ]
func (c OPEigenDAGenericCommitment) Encode() []byte {
	return append([]byte{byte(OPGenericCommitmentByte), EigenDALayerByte}, c.versionedCert.Encode()...)
}
