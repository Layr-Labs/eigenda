package commitments

import "github.com/ethereum/go-ethereum/crypto"

type DAServiceOPCommitmentType byte

const (
	EigenDACommitmentType DAServiceOPCommitmentType = 0
)

// OPCommitment is the binary representation of a commitment.
type DaSvcCommitment interface {
	CommitmentType() DAServiceOPCommitmentType
	Encode() []byte
	Verify(input []byte) error
}

type EigenDASvcCommitment []byte


// NewEigenDASvcCommitment creates a new commitment from the given input.
func NewEigenDASvcCommitment(input []byte) EigenDASvcCommitment {
	return EigenDASvcCommitment(crypto.Keccak256(input))
}

// DecodeEigenDASvcCommitment validates and casts the commitment into a Keccak256Commitment.
func DecodeEigenDASvcCommitment(commitment []byte) (EigenDASvcCommitment, error) {
	// guard against empty commitments
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	return commitment, nil
}

// CommitmentType returns the commitment type of Keccak256.
func (c EigenDASvcCommitment) CommitmentType() DAServiceOPCommitmentType {
	return EigenDACommitmentType
}

// Encode adds a commitment type prefix self describing the commitment.
func (c EigenDASvcCommitment) Encode() []byte {
	return append([]byte{byte(EigenDACommitmentType)}, c...)
}