package commitments

import (
	"fmt"
)

type CommitmentMeta struct {
	Mode CommitmentMode
	// CertVersion is shared for all modes and denotes version of the EigenDA certificate
	CertVersion uint8
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

func EncodeCommitment(b []byte, c CommitmentMode) ([]byte, error) {
	switch c {
	case OptimismKeccak:
		return Keccak256Commitment(b).Encode(), nil

	case OptimismGeneric:
		certCommit := NewV0CertCommitment(b).Encode()
		svcCommit := EigenDASvcCommitment(certCommit).Encode()
		altDACommit := NewGenericCommitment(svcCommit).Encode()
		return altDACommit, nil

	case Standard:
		return NewV0CertCommitment(b).Encode(), nil
	}

	return nil, fmt.Errorf("unknown commitment mode")
}
