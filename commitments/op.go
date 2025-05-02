package commitments

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// ErrInvalidCommitment is returned when the commitment cannot be parsed into a known commitment type.
var ErrInvalidCommitment = errors.New("invalid commitment")

// ErrCommitmentMismatch is returned when the commitment does not match the given input.
var ErrCommitmentMismatch = errors.New("commitment mismatch")

// OPCommitmentType is the commitment type prefix.
type OPCommitmentType byte

// CommitmentType describes the binary format of the commitment.
// KeccakCommitmentStringType is the default commitment type for optimism's centralized DA storage.
// GenericCommitmentType indicates an opaque bytestring that the op-node never opens.
const (
	Keccak256CommitmentType OPCommitmentType = 0
	GenericCommitmentType   OPCommitmentType = 1
	KeccakCommitmentString  string           = "KeccakCommitment"
	GenericCommitmentString string           = "GenericCommitment"
)

func CommitmentTypeFromString(s string) (OPCommitmentType, error) {
	switch s {
	case KeccakCommitmentString:
		return Keccak256CommitmentType, nil
	case GenericCommitmentString:
		return GenericCommitmentType, nil
	default:
		return 0, fmt.Errorf("invalid commitment type: %s", s)
	}
}

// OPCommitment is the binary representation of a commitment.
type OPCommitment interface {
	CommitmentType() OPCommitmentType
	Encode() []byte
	Verify(input []byte) error
}

// Keccak256Commitment is an implementation of OPCommitment that uses Keccak256 as the commitment function.
type Keccak256Commitment []byte

// GenericCommitment is an implementation of OPCommitment that treats the commitment as an opaque bytestring.
type GenericCommitment []byte

// NewOPCommitment creates a new commitment from the given input and desired type.
func NewOPCommitment(t OPCommitmentType, input []byte) OPCommitment {
	switch t {
	case Keccak256CommitmentType:
		return NewKeccak256Commitment(input)
	case GenericCommitmentType:
		return NewGenericCommitment(input)
	default:
		return nil
	}
}

// DecodeOPCommitment parses the commitment into a known commitment type.
// The input type is determined by the first byte of the raw data.
// The input type is discarded and the commitment is passed to the appropriate constructor.
func DecodeOPCommitment(input []byte) (OPCommitment, error) {
	if len(input) == 0 {
		return nil, ErrInvalidCommitment
	}
	t := OPCommitmentType(input[0])
	data := input[1:]
	switch t {
	case Keccak256CommitmentType:
		return DecodeKeccak256(data)
	case GenericCommitmentType:
		return DecodeGenericCommitment(data)
	default:
		return nil, ErrInvalidCommitment
	}
}

// NewKeccak256Commitment creates a new commitment from the given input.
func NewKeccak256Commitment(input []byte) Keccak256Commitment {
	return Keccak256Commitment(crypto.Keccak256(input))
}

// DecodeKeccak256 validates and casts the commitment into a Keccak256Commitment.
func DecodeKeccak256(commitment []byte) (Keccak256Commitment, error) {
	// guard against empty commitments
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	// keccak commitments are always 32 bytes
	if len(commitment) != 32 {
		return nil, ErrInvalidCommitment
	}
	return commitment, nil
}

// CommitmentType returns the commitment type of Keccak256.
func (c Keccak256Commitment) CommitmentType() OPCommitmentType {
	return Keccak256CommitmentType
}

// Encode adds a commitment type prefix self describing the commitment.
func (c Keccak256Commitment) Encode() []byte {
	return append([]byte{byte(Keccak256CommitmentType)}, c...)
}

// Verify checks if the commitment matches the given input.
func (c Keccak256Commitment) Verify(input []byte) error {
	if !bytes.Equal(c, crypto.Keccak256(input)) {
		return ErrCommitmentMismatch
	}
	return nil
}

// NewGenericCommitment creates a new commitment from the given input.
func NewGenericCommitment(input []byte) GenericCommitment {
	return GenericCommitment(input)
}

// DecodeGenericCommitment validates and casts the commitment into a GenericCommitment.
func DecodeGenericCommitment(commitment []byte) (GenericCommitment, error) {
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	return commitment, nil
}

// CommitmentType returns the commitment type of Generic Commitment.
func (c GenericCommitment) CommitmentType() OPCommitmentType {
	return GenericCommitmentType
}

// Encode adds a commitment type prefix self describing the commitment.
func (c GenericCommitment) Encode() []byte {
	return append([]byte{byte(GenericCommitmentType)}, c...)
}

// Verify always returns true for GenericCommitment because the DA Server
// must validate the data before returning it to the op-node.
func (c GenericCommitment) Verify(_ []byte) error {
	return nil
}
