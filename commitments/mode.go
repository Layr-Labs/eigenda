package commitments

import (
	"encoding/hex"
	"fmt"
)

type CommitmentMeta struct {
	Mode CommitmentMode
	// CertVersion is shared for all modes and denotes version of the EigenDA certificate
	CertVersion byte
}

type CommitmentMode string

const (
	OptimismKeccak       CommitmentMode = "optimism_keccak256"
	OptimismGeneric      CommitmentMode = "optimism_generic"
	SimpleCommitmentMode CommitmentMode = "simple"
)

func StringToCommitmentMode(s string) (CommitmentMode, error) {
	switch s {
	case string(OptimismKeccak):
		return OptimismKeccak, nil
	case string(OptimismGeneric):
		return OptimismGeneric, nil
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
	case OptimismKeccak: // [op_type, ...]
		return b[1:], nil

	case OptimismGeneric: // [op_type, da_provider, cert_version, ...]
		return b[3:], nil

	case SimpleCommitmentMode: // [cert_version, ...]
		return b[1:], nil

	default:
		return nil, fmt.Errorf("unknown commitment type")
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

	case SimpleCommitmentMode:
		return NewV0CertCommitment(b).Encode(), nil
	}

	return nil, fmt.Errorf("unknown commitment mode")
}
