package commitments

import (
	"encoding/hex"
	"fmt"
)

type CommitmentMode string

const (
	OptimismGeneric      CommitmentMode = "optimism_keccak256"
	OptimismAltDA        CommitmentMode = "optimism_generic"
	SimpleCommitmentMode CommitmentMode = "simple"
)

func StringToCommitmentMode(s string) (CommitmentMode, error) {
	switch s {
	case string(OptimismGeneric):
		return OptimismGeneric, nil
	case string(OptimismAltDA):
		return OptimismAltDA, nil
	case string(SimpleCommitmentMode):
		return SimpleCommitmentMode, nil
	default:
		return "", fmt.Errorf("unknown commitment mode: %s", s)
	}
}

func StringToDecodedCommitment(key string, c CommitmentMode) ([]byte, error) {
	offset := 0
	if key[:2] == "0x" {
		offset = 2
	}

	b, err := hex.DecodeString(key[offset:])
	if err != nil {
		return nil, err
	}

	if len(b) < 3 {
		return nil, fmt.Errorf("commitment is too short")
	}

	switch c {
	case OptimismGeneric: // [op_type, ...]
		return b[1:], nil

	case OptimismAltDA: // [op_type, da_provider, cert_version, ...]
		return b[3:], nil

	case SimpleCommitmentMode: // [cert_version, ...]
		return b[1:], nil

	default:
		return nil, fmt.Errorf("unknown commitment type")
	}
}

func EncodeCommitment(b []byte, c CommitmentMode) ([]byte, error) {
	switch c {
	case OptimismGeneric:
		return Keccak256Commitment(b).Encode(), nil

	case OptimismAltDA:
		certCommit := NewV0CertCommitment(b).Encode()
		svcCommit := EigenDASvcCommitment(certCommit).Encode()
		altDACommit := NewGenericCommitment(svcCommit).Encode()
		return altDACommit, nil

	case SimpleCommitmentMode:
		return NewV0CertCommitment(b).Encode(), nil
	}

	return nil, fmt.Errorf("unknown commitment mode")
}
