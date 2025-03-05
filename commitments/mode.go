package commitments

import (
	"fmt"
)

type CommitmentMeta struct {
	Mode CommitmentMode
	// version is shared for all modes and denotes version of the EigenDA certificate
	Version EigenDACommitmentType
}

type CommitmentMode string

const (
	OptimismKeccak  CommitmentMode = "optimism_keccak256"
	OptimismGeneric CommitmentMode = "optimism_generic"
	Standard        CommitmentMode = "standard"
)

func StringToCommitmentMode(s string) (CommitmentMode, error) {
	switch s {
	case string(OptimismKeccak):
		return OptimismKeccak, nil
	case string(OptimismGeneric):
		return OptimismGeneric, nil
	case string(Standard):
		return Standard, nil
	default:
		return "", fmt.Errorf("unknown commitment mode: %s", s)
	}
}

func EncodeCommitment(
	bytes []byte,
	commitmentMode CommitmentMode,
	commitmentType EigenDACommitmentType,
) ([]byte, error) {
	switch commitmentMode {
	case OptimismKeccak:
		return Keccak256Commitment(bytes).Encode(), nil

	case OptimismGeneric:
		certCommit := NewEigenDACommitment(bytes, commitmentType).Encode()
		svcCommit := EigenDASvcCommitment(certCommit).Encode()
		altDACommit := NewGenericCommitment(svcCommit).Encode()
		return altDACommit, nil

	case Standard:
		return NewEigenDACommitment(bytes, commitmentType).Encode(), nil
	}

	return nil, fmt.Errorf("unknown commitment mode")
}
