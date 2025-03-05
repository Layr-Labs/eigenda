package commitments

type EigenDACommitmentType byte

const (
	// EigenDA V1
	CertV0 EigenDACommitmentType = iota
	// EigenDA V2
	CertV1
)

// CertCommitment is the binary representation of a commitment.
type CertCommitment interface {
	CommitmentType() EigenDACommitmentType
	Encode() []byte
	Verify(input []byte) error
}

type EigenDACommitment struct {
	prefix EigenDACommitmentType
	b      []byte
}

// NewEigenDACommitment creates a new commitment from the given input.
func NewEigenDACommitment(input []byte, commitmentType EigenDACommitmentType) EigenDACommitment {
	return EigenDACommitment{
		prefix: commitmentType,
		b:      input,
	}
}

// CommitmentType returns the commitment type of EigenDACommitment.
func (c EigenDACommitment) CommitmentType() EigenDACommitmentType {
	return c.prefix
}

// Encode adds a commitment type prefix self describing the commitment.
func (c EigenDACommitment) Encode() []byte {
	return append([]byte{byte(c.prefix)}, c.b...)
}
