package server

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
)

type CommitmentMode string

const (
	OptimismCommitmentMode CommitmentMode = "optimism"
	SimpleCommitmentMode   CommitmentMode = "simple"
)

func StringToCommitmentMode(s string) (CommitmentMode, error) {
	switch s {
	case string(OptimismCommitmentMode):
		return OptimismCommitmentMode, nil
	case string(SimpleCommitmentMode):
		return SimpleCommitmentMode, nil
	default:
		return "", fmt.Errorf("unknown commitment mode: %s", s)
	}
}

func StringToCommitment(key string, c CommitmentMode) ([]byte, error) {
	if len(key) <= 2 {
		return nil, fmt.Errorf("commitment is empty")
	}

	if key[:2] != "0x" {
		return nil, fmt.Errorf("commitment parameter does not have 0x prefix")
	}

	b, err := hex.DecodeString(key[2:])
	if err != nil {
		return nil, err
	}

	switch c {
	case OptimismCommitmentMode:
		var comm commitments.OPCommitment
		err = comm.Unmarshal(b)
		if err != nil {
			return nil, err
		}
		if !comm.IsGenericCommitment() {
			return nil, fmt.Errorf("commitment is not a OP DA service commitment")
		}
		daComm := comm.MustGenericCommitmentValue()
		if !daComm.IsEigenDA() {
			return nil, fmt.Errorf("commitment is not an EigenDA OP DA service commitment")
		}
		eigendaComm := daComm.MustEigenDAValue()
		if !eigendaComm.IsCertV0() {
			return nil, fmt.Errorf("commitment is not a supported EigenDA cert encoding")
		}
		return eigendaComm.MustCertV0Value(), nil
	case SimpleCommitmentMode:
		var eigendaComm commitments.EigenDACommitment
		err = eigendaComm.Unmarshal(b)
		if err != nil {
			return nil, err
		}
		if !eigendaComm.IsCertV0() {
			return nil, fmt.Errorf("commitment is not a supported EigenDA cert encoding")
		}
		return eigendaComm.MustCertV0Value(), nil
	default:
		return nil, fmt.Errorf("unknown commitment type")
	}
}

func EncodeCommitment(s []byte, c CommitmentMode) ([]byte, error) {
	switch c {
	case OptimismCommitmentMode:
		comm := commitments.GenericCommitment(commitments.OptimismEigenDACommitment(commitments.EigenDACertV0(s)))
		return comm.Marshal()
	case SimpleCommitmentMode:
		comm := commitments.EigenDACertV0(s)
		return comm.Marshal()
	default:
		return nil, fmt.Errorf("unknown commitment type")
	}
}

func ReadCommitmentMode(r *http.Request) (CommitmentMode, error) {
	query := r.URL.Query()
	key := query.Get(CommitmentModeKey)
	if key == "" { // default
		return OptimismCommitmentMode, nil
	}
	return StringToCommitmentMode(key)
}
