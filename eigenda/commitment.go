package eigenda

import (
	"errors"

	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ErrCommitmentLength is returned when the commitment length is invalid.
var ErrCommitmentLength = errors.New("invalid commitment length")

// ErrInvalidCommitment is returned when the commitment cannot be parsed into a known commitment type.
var ErrInvalidCommitment = errors.New("invalid commitment")

// ErrCommitmentMismatch is returned when the commitment does not match the given input.
var ErrCommitmentMismatch = errors.New("commitment mismatch")

// ExtDAType is the DA provider type.
type ExtDAType byte

const (
	EigenDA ExtDAType = 0x00
)

// EigenDAVersion is the version being used for EigenDA.
type EigenDAVersion byte

const (
	EigenV0 EigenDAVersion = 0x00
)

type Commitment []byte

func (c Commitment) Encode() []byte {
	return append([]byte{byte(op_plasma.GenericCommitmentType), byte(EigenDA), byte(EigenV0)}, c...)
}

func StringToCommit(key string) (Commitment, error) {
	comm, err := hexutil.Decode(key)
	if err != nil {
		return nil, err
	}
	return DecodeCommitment(comm)
}

// DecodeCommitment verifies and decodes an EigenDACommit from raw encoded bytes.
func DecodeCommitment(commitment []byte) (Commitment, error) {
	if len(commitment) <= 3 {
		return nil, ErrCommitmentLength
	}
	if commitment[0] != byte(op_plasma.GenericCommitmentType) {
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
