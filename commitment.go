package plasma

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

// ErrCommitmentLength is returned when the commitment length is invalid.
var ErrCommitmentLength = errors.New("invalid commitment length")

// ErrInvalidCommitment is returned when the commitment cannot be parsed into a known commitment type.
var ErrInvalidCommitment = errors.New("invalid commitment")

// ErrCommitmentMismatch is returned when the commitment does not match the given input.
var ErrCommitmentMismatch = errors.New("commitment mismatch")

// CommitmentType is the commitment type prefix.
type CommitmentType byte

const (
	// default commitment type for the DA storage.
	Keccak256CommitmentType CommitmentType = 0
	GenericCommitment       CommitmentType = 1
)

type ExtDAType byte

const (
	EigenDA ExtDAType = 0x00
)

type EigenDAVersion byte

const (
	EigenV0 EigenDAVersion = 0x00
)

// Keccak256Commitment is the default commitment type for op-plasma.
type Keccak256Commitment []byte

// Encode adds a commitment type prefix self describing the commitment.
func (c Keccak256Commitment) Encode() []byte {
	return append([]byte{byte(Keccak256CommitmentType)}, c...)
}

// TxData adds an extra version byte to signal it's a commitment.
func (c Keccak256Commitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

// Verify checks if the commitment matches the given input.
func (c Keccak256Commitment) Verify(input []byte) error {
	if !bytes.Equal(c, crypto.Keccak256(input)) {
		return ErrCommitmentMismatch
	}
	return nil
}

// Keccak256 creates a new commitment from the given input.
func Keccak256(input []byte) Keccak256Commitment {
	return Keccak256Commitment(crypto.Keccak256(input))
}

// DecodeKeccak256 validates and casts the commitment into a Keccak256Commitment.
func DecodeKeccak256(commitment []byte) (Keccak256Commitment, error) {
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	if commitment[0] != byte(Keccak256CommitmentType) {
		return nil, ErrInvalidCommitment
	}
	c := commitment[1:]
	if len(c) != 32 {
		return nil, ErrInvalidCommitment
	}
	return c, nil
}

type EigenDACommitment []byte

func (c EigenDACommitment) Encode() []byte {
	return append([]byte{byte(EigenDA), byte(EigenV0)}, c...)
}

func (c EigenDACommitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

func DecodeEigenDACommitment(commitment []byte) (EigenDACommitment, error) {
	if len(commitment) <= 3 {
		return nil, ErrCommitmentLength
	}
	if commitment[0] != byte(GenericCommitment) {
		return nil, ErrInvalidCommitment
	}

	if commitment[1] != byte(EigenDA) {
		return nil, ErrInvalidCommitment
	}

	// additional versions will need to be hardcoded here
	if commitment[2] != byte(EigenV0) {
		return nil, ErrInvalidCommitment
	}

	c := commitment[3:]

	// TODO - Add a length check
	return c, nil
}
