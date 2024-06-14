package commitments

import (
	"fmt"
	"log"
)

type DAServiceOPCommitmentType byte

const (
	EigenDAByte DAServiceOPCommitmentType = 0
)

// DAServiceOPCommitment represents a value of one of two possible types (Keccak256Commitment or DAServiceCommitment).
type DAServiceOPCommitment struct {
	eigendaCommitment *EigenDACommitment
}

var _ Commitment = (*DAServiceOPCommitment)(nil)

func OptimismEigenDACommitment(value EigenDACommitment) DAServiceOPCommitment {
	return DAServiceOPCommitment{eigendaCommitment: &value}
}

func (e DAServiceOPCommitment) IsEigenDA() bool {
	return e.eigendaCommitment != nil
}

func (e DAServiceOPCommitment) MustEigenDAValue() EigenDACommitment {
	if e.eigendaCommitment != nil {
		return *e.eigendaCommitment
	}
	log.Panic("CommitmentEither does not contain a Keccak256Commitment value")
	return EigenDACommitment{} // This will never be reached, but is required for compilation.
}

func (e DAServiceOPCommitment) Marshal() ([]byte, error) {
	if e.IsEigenDA() {
		eigenDABytes, err := e.MustEigenDAValue().Marshal()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(EigenDAByte)}, eigenDABytes...), nil
	} else {
		return nil, fmt.Errorf("DAServiceOPCommitment is neither a keccak256 commitment or a DA service commitment")
	}
}

func (e *DAServiceOPCommitment) Unmarshal(bz []byte) error {
	if len(bz) < 1 {
		return fmt.Errorf("OP commitment does not contain generic commitment type prefix byte")
	}
	head := DAServiceOPCommitmentType(bz[0])
	tail := bz[1:]
	switch head {
	case EigenDAByte:
		eigendaCommitment := EigenDACommitment{}
		err := eigendaCommitment.Unmarshal(tail)
		if err != nil {
			return err
		}
		e.eigendaCommitment = &eigendaCommitment
	default:
		return fmt.Errorf("unrecognized generic commitment type byte: %x", bz[0])
	}
	return nil
}
