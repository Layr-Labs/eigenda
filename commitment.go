package plasma

import (
	"errors"
)

// ErrInvalidCommitment is returned when the commitment cannot be parsed into a known commitment type.
var ErrInvalidCommitment = errors.New("invalid commitment")

// ErrCommitmentMismatch is returned when the commitment does not match the given input.
var ErrCommitmentMismatch = errors.New("commitment mismatch")

// CommitmentType is the commitment type prefix.
type CommitmentType byte

// Max input size ensures the canonical chain cannot include input batches too large to
// challenge in the Data Availability Challenge contract. Value in number of bytes.
// This value can only be changed in a hard fork.
const MaxInputSize = 130672

// TxDataVersion1 is the version number for batcher transactions containing
// plasma commitments. It should not collide with DerivationVersion which is still
// used downstream when parsing the frames.
const TxDataVersion1 = 1

const (
	// default commitment type for the DA storage.
	Keccak256CommitmentType CommitmentType = 0x00
	DaService               CommitmentType = 0x01
)

type ExtDAType byte

const (
	EigenDA ExtDAType = 0x00
)

type EigenDAVersion byte

const (
	EigenV0 EigenDAVersion = 0x00
)

// NOTE - This logic will need to be migrated into layr-labs/op-stack directly
type EigenDACommitment []byte

func (c EigenDACommitment) Encode() []byte {
	return append([]byte{byte(DaService), byte(EigenDA), byte(EigenV0)}, c...)
}

func (c EigenDACommitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

func DecodeEigenDACommitment(commitment []byte) (EigenDACommitment, error) {
	if len(commitment) <= 3 {
		return nil, ErrInvalidCommitment
	}
	if commitment[0] != byte(DaService) {
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
