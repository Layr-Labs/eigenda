package commitments


type CertEncodingCommitment byte

const (
	CertV0 CertEncodingCommitment = 0
)

// OPCommitment is the binary representation of a commitment.
type CertCommitment interface {
	CommitmentType() CertEncodingCommitment
	Encode() []byte
	Verify(input []byte) error
}

type CertCommitmentV0 []byte


// NewV0CertCommitment creates a new commitment from the given input.
func NewV0CertCommitment(input []byte) CertCommitmentV0 {
	return CertCommitmentV0(input)
}

// DecodeCertCommitment validates and casts the commitment into a Keccak256Commitment.
func DecodeCertCommitment(commitment []byte) (CertCommitmentV0, error) {
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	return commitment, nil
}

// CommitmentType returns the commitment type of Keccak256.
func (c CertCommitmentV0) CommitmentType() CertEncodingCommitment {
	return CertV0
}

// Encode adds a commitment type prefix self describing the commitment.
func (c CertCommitmentV0) Encode() []byte {
	return append([]byte{byte(CertV0)}, c...)
}